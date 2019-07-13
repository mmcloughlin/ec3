package mont

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"

	"github.com/mmcloughlin/ec3/asm"
	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

func New(p prime.Prime) fp.Field {
	return field{p: p}
}

type field struct {
	p prime.Prime
}

func (f field) ElementBits() int {
	n := f.p.Bits()
	return ints.NextMultiple(n, 64)
}

func (f field) ElementSize() int {
	return f.ElementBits() / 8
}

func (f field) Limbs() int {
	return f.ElementBits() / 64
}

func (f field) Build(ctx *build.Context) fp.Builder {
	return &builder{
		field:   f,
		Context: ctx,
	}
}

type builder struct {
	field
	*build.Context

	modulus mp.Int
}

func (b builder) Add(x, y mp.Int) {
	k := b.Limbs()

	// Add as multi-precision integers, allowing a carry into a high word.
	// TODO(mbm): consider case when prime size is not a multiple of 64, so carry would not overflow.
	carry := asm.Zero64(b.Context)
	sum := mp.Int{}
	sum = append(sum, x...)
	sum = append(sum, carry)

	// TODO(mbm): Add() function in mp package
	b.ADDQ(y[0], sum[0])
	for i := 1; i < k; i++ {
		b.ADCQ(y[i], sum[i])
	}
	b.ADCQ(operand.U32(0), sum[k])

	b.ConditionalSubtractModulus(sum)
}

// ConditionalSubtractModulus subtracts p from x if x ⩾ p in constant time.
func (b builder) ConditionalSubtractModulus(x mp.Int) {
	subp := mp.CopyIntoRegisters(b.Context, x)
	p := b.P()

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

// P returns the prime modulus p as a multi-precision integer.
func (b *builder) P() mp.Int {
	if b.modulus != nil {
		return b.modulus
	}
	limbs := bigint.Uint64s(b.p.Int())
	b.modulus = mp.StaticGlobal(b.Context, "p", limbs)
	return b.modulus
}

func (b builder) ReduceDouble(z, x mp.Int) {
}
