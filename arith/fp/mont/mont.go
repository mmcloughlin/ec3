package mont

import (
	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/arith/mp"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

type Field struct {
	p prime.Prime
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

// Modulus returns the prime modulus p as a multi-precision integer.
func (f *Field) Modulus() ir.Constants {
	return ir.NewConstantsFromInt(f.p.Int(), 64)
}

func (f *Field) Add(ctx *build.Context, z, x, y ir.Int) {
	// TODO(mbm): consider case when prime size is not a multiple of 64, so carry would not overflow.

	// Add as multi-precision integers, allowing a carry into a high word.
	carry := ctx.Register("carry")
	sum := ir.Extend(z, carry)
	mp.AddInto(ctx, sum, x, y, ctx.Register("c"))

	// Possibly subtract modulus.
	f.ConditionalSubtractModulus(ctx, sum)
}

// ConditionalSubtractModulus subtracts p from x if x ⩾ p in constant time.
func (f *Field) ConditionalSubtractModulus(ctx *build.Context, x ir.Int) {
	k := f.Limbs()
	mod := f.Modulus()

	// Subtract modulus.
	subp := ctx.Int("subp", k+1)
	b := ctx.Register("b")
	mp.SubInto(ctx, subp, x, mod, b)

	// Keep the subtracted value if borrow is 0.
	mp.ConditionalMove(ctx, x, subp, b, 0)
}
