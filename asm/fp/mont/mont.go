package mont

import (
	"math/big"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"

	"github.com/mmcloughlin/ec3/asm"
	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

// References:
//
//	[fastprimefieldecc]  Shay Gueron and Vlad Krasnov. Fast Prime Field Elliptic Curve Cryptography with
//	                     256 Bit Primes. Cryptology ePrint Archive, Report 2013/816. 2013.
//	                     https://eprint.iacr.org/2013/816
//	[hac:impl]           Alfred J. Menezes, Paul C. van Oorschot and Scott A. Vanstone. Efficient
//	                     Implementation. Handbook of Applied Cryptography, chapter 14. 1996.
//	                     http://cacr.uwaterloo.ca/hac/about/chap14.pdf

func New(p prime.Prime) fp.Field {
	return Field{p: p}
}

type Field struct {
	p prime.Prime
}

func (f Field) Prime() *big.Int {
	return f.p.Int()
}

func (f Field) ElementBits() int {
	n := f.p.Bits()
	return ints.NextMultiple(n, 64)
}

func (f Field) ElementSize() int {
	return f.ElementBits() / 8
}

func (f Field) Limbs() int {
	return f.ElementBits() / 64
}

func (f Field) Build(ctx *build.Context) fp.Builder {
	return &builder{
		Field:   f,
		Context: ctx,
	}
}

type builder struct {
	Field
	*build.Context

	modulus mp.Int
	mprime  operand.Op
}

func (b *builder) Add(x, y mp.Int) {
	k := b.Limbs()

	// Add as multi-precision integers, allowing a carry into a high word.
	// TODO(mbm): consider case when prime size is not a multiple of 64, so carry would not overflow.
	carry := asm.Zero64(b.Context)
	sum := x.Extend(carry)

	// TODO(mbm): Add() function in mp package
	b.ADDQ(y[0], sum[0])
	for i := 1; i < k; i++ {
		b.ADCQ(y[i], sum[i])
	}
	b.ADCQ(operand.U32(0), sum[k])

	b.ConditionalSubtractModulus(sum)
}

func (b *builder) Sub(x, y mp.Int) {
	k := b.Limbs()

	// Subtract multi-precision integers, allowing a borrow into a high word.
	borrow := asm.Zero64(b.Context)

	// TODO(mbm): Sub() function in mp package
	b.SUBQ(y[0], x[0])
	for i := 1; i < k; i++ {
		b.SBBQ(y[i], x[i])
	}
	b.SBBQ(operand.U32(0), borrow)

	// Compute x + p.
	addp := mp.CopyIntoRegisters(b.Context, x)

	p := b.Modulus()
	b.ADDQ(p[0], addp[0])
	for i := 1; i < k; i++ {
		b.ADCQ(p[i], addp[i])
	}

	// If the borrow is non-zero, that means we need to take x+p.
	b.ANDQ(operand.U32(1), borrow)
	for i := 0; i < k; i++ {
		b.CMOVQNE(addp[i], x[i])
	}
}

// ConditionalSubtractModulus subtracts p from x if x ⩾ p in constant time.
func (b *builder) ConditionalSubtractModulus(x mp.Int) {
	subp := mp.CopyIntoRegisters(b.Context, x)
	p := b.Modulus()

	// Subtract p.
	// TODO(mbm): Sub() function in mp package
	b.SUBQ(p[0], subp[0])
	for i := 1; i < len(p); i++ {
		b.SBBQ(p[i], subp[i])
	}
	for i := len(p); i < len(subp); i++ {
		b.SBBQ(operand.U32(0), subp[i])
	}

	// Conditionally move.
	for i := 0; i < b.Limbs(); i++ {
		b.CMOVQCC(subp[i], x[i])
	}
}

// Modulus returns the prime modulus p as a multi-precision integer.
func (b *builder) Modulus() mp.Int {
	if b.modulus != nil {
		return b.modulus
	}
	limbs := bigint.Uint64s(b.p.Int())
	b.modulus = mp.StaticGlobal(b.Context, "p", limbs)
	return b.modulus
}

// IsFriendly reports whether the modulus is "Montgomery friendly". This concept
// is introduced in [fastprimefieldecc], Section 3. The modulus is considered
// friendly if the m' value is one, which allows us to save a multiply in each
// Montgomery reduction step.
func (b *builder) IsFriendly() bool {
	// Note that m' ≡ -1/m (mod b), therefore m' is 1 when m ≡ -1 (mod b).
	m := b.p.Int()
	mplus1mod64 := new(big.Int).Add(m, bigint.One())
	mplus1mod64.Mod(mplus1mod64, bigint.Pow2(64))
	return bigint.IsZero(mplus1mod64)
}

// ModulusPrime returns the m' value required for Montgomery reduction, specifically:
//	m' ≡ -1/m (mod b)
// where b is the base 2⁶⁴. Note for certain friendly primes this value will be
// one, in which case it is not necessary to generate a constant for it. See
// IsFriendly to check for this case.
func (b *builder) ModulusPrime() operand.Op {
	if b.mprime != nil {
		return b.mprime
	}
	base := bigint.Pow2(64)
	m := b.p.Int()
	mprime := new(big.Int).ModInverse(m, base)
	mprime.Sub(base, mprime)
	mprimeint := mp.StaticGlobal(b.Context, "mprime", []uint64{mprime.Uint64()})
	b.mprime = mprimeint[0]
	return b.mprime
}

func (b builder) ReduceDouble(z, x mp.Int) {
	// Reduction is performed with multi-word Montgomery reduction. See [hac:impl]
	// Algorithm 14.32.

	k := b.Limbs()

	// We'll need a zero register.
	zero := asm.Zero64(b.Context)

	// Set up accumulator registers.
	acc := mp.NewIntLimb64(b.Context, 2*k+1)
	mp.Copy(b.Context, acc, x[:k])

	// Step 2: iterate over limbs
	pending := mp.NewIntLimb64(b.Context, 2*k+1)
	pending[k] = zero
	for i := 0; i < k; i++ {
		// Step 2.1: u_i = x_i * m' (mod b)
		var u operand.Op
		if b.IsFriendly() {
			u = acc[i]
		} else {
			mprime := b.ModulusPrime()
			b.MOVQ(mprime, reg.RDX)
			lo, hi := b.GP64(), b.GP64()
			b.MULXQ(acc[i], lo, hi)
			u = lo
		}

		// Step 2.2: x += u_i * m * b^i
		b.MOVQ(x.Limb(i+k), acc[i+k])
		b.MOVQ(u, reg.RDX)
		m := b.Modulus()
		b.XORQ(pending[i+k+1], pending[i+k+1]) // also clears flags
		for j := 0; j < k; j++ {
			lo, hi := b.GP64(), b.GP64()
			b.MULXQ(m[j], lo, hi)
			b.ADCXQ(lo, acc[i+j])
			b.ADOXQ(hi, acc[i+j+1])
		}
		b.ADCXQ(pending[i+k], acc[i+k])
		b.ADCXQ(zero, pending[i+k+1])
		b.ADOXQ(zero, pending[i+k+1])
	}
	acc[2*k] = pending[2*k]

	// Step 4: if x ⩾ m subtract m
	result := acc[k:]
	b.ConditionalSubtractModulus(result)

	// Write result.
	for i := 0; i < k; i++ {
		b.MOVQ(result[i], z[i])
	}
}
