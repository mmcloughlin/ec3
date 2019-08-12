package crandall

import (
	"math/big"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"

	"github.com/mmcloughlin/ec3/asm"
	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

func New(p prime.Crandall) fp.Field {
	return Field{p: p}
}

type Field struct {
	p prime.Crandall
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

// ReductionMultiplier returns the value an element must be multiplied by to reduce it modulo p upon overflow of the element size.
// Note that for the Crandall prime we have
//
//	2ⁿ - c ≡ 0 (mod p)
//	2ⁿ ≡ c (mod p)
//
// However n may not be on a limb boundary, so we actually need the identity
//
//	2ˡ ≡ 2ˡ⁻ⁿ * c (mod p)
//
// We'll call this the reduction multiplier.
func (f Field) ReductionMultiplier() uint32 {
	n := f.p.Bits()
	l := f.ElementBits()
	// TODO(mbm): check for overflow
	return uint32((1 << uint(l-n)) * f.p.C)
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
}

func (b builder) Add(x, y mp.Int) {
	k := b.Limbs()

	// Prepare a zero register.
	zero := asm.Zero64(b.Context)

	// Load reduction multiplier.
	d := b.ReductionMultiplier()
	dreg := b.GP64()
	b.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

	// Add y into x.
	b.ADDQ(y[0], x[0]) // TODO(mbm): can we replace this with `ADCX`? need to ensure the carry flag is 0
	for i := 1; i < k; i++ {
		b.ADCXQ(y[i], x[i])
	}

	// Both inputs are < 2ˡ so the result is < 2ˡ⁺¹.
	// If the last addition caused a carry into the l'th bit we need to perform a reduction.
	// Prepare the value we will add in to perform the reduction. The addend may be
	// zero or d depending on the carry.
	addend := b.GP64()
	b.MOVQ(zero, addend)
	b.CMOVQCS(dreg, addend)

	// Now add the addend into x.
	b.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
	for i := 1; i < k; i++ {
		b.ADCXQ(zero, x[i])
	}

	// We have added d into the low l bits. Therefore the result is less than 2ˡ + d.
	// But note that it could still be 2ˡ or higher, so we need to perform a
	// second reduction.

	// As before, the addend is either 0 or d depending on the carry from the last add.
	b.MOVQ(zero, addend)
	b.CMOVQCS(dreg, addend)

	// This time we only need to perform one add. The result must be less than 2ˡ + 2*d,
	// therefore provided 2*d does not exceed the size of a limb we can be sure there
	// will be no carry.
	// TODO(mbm): assert d is within an acceptable range
	b.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
}

func (b *builder) Sub(x, y mp.Int) {
	// TODO(mbm): implement subtraction in Crandall fields.
	panic(errutil.ErrNotImplemented)
}

// ReduceDouble computes z congruent to x modulo p. Let the element size be 2ˡ.
// This function assumes x < 2²ˡ and produces z < 2ˡ. Note that z is not
// guaranteed to be less than p.
func (b builder) ReduceDouble(z, x mp.Int) {
	k := b.Limbs()

	// Prepare a zero register.
	zero := asm.Zero64(b.Context)

	// Compute the reduction multiplier.
	d := b.ReductionMultiplier()
	dreg := b.GP64()
	b.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

	// Stage 1: upper bound 2²ˡ → 2ˡ + d*2ˡ.
	//
	// At this point x is a 2*l bit quantity which we'll view as two l-bit halfs.
	//
	//	x = H ∥ L = 2ˡ * H + L ≡ d*H + L (mod p)
	//
	// Therefore the first reduction stage multiplies the top limbs by d and
	// accumulates the result into the low limbs. Note the result will have an
	// additional limb.

	// Multiply r = d*H.
	r := mp.NewIntLimb64(b.Context, k+1)
	b.MOVQ(dreg, reg.RDX)
	b.XORQ(r[0], r[0]) // also clears flags
	for i := 0; i < k; i++ {
		lo := b.GP64()
		b.MULXQ(x[i+k], lo, r[i+1])
		b.ADCXQ(lo, r[i])
	}

	// Add r += x.
	for i := 0; i < k; i++ {
		b.ADOXQ(x[i], r[i])
	}
	b.ADOXQ(zero, r[k])

	// Stage 2: (d+1)*2ˡ → 2ˡ + (d+1)*d
	//
	// Currently r has one too many limbs, so we need to reduce again. The value in
	// the top limb is ⩽ d. When we reduce we have to multiply by d again, so the
	// result cannot exceed d^2. Provided d is small, the result will not exceed a single limb.

	// TODO(mbm): assert d is within an acceptable range

	top := r[k]
	b.IMULQ(dreg, top) // clears flags
	b.ADCXQ(top, r[0])
	for i := 1; i < k; i++ {
		b.ADCXQ(zero, r[i])
	}

	// Stage 3: finish
	//
	// It is still possible that the final add carried, in which case we need one final
	// add to complete the reduction.
	addend := b.GP64()
	b.MOVQ(zero, addend)
	b.CMOVQCS(dreg, addend)
	b.ADDQ(addend, r[0])

	// Write out the result.
	for i := 0; i < k; i++ {
		b.MOVQ(r[i], z[i])
	}
}
