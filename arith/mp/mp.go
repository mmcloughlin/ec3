// Package mp provides multi-precision operations for arithmetic programs.
package mp

import (
	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/ints"
)

// ConditionalMove moves src to dst if f≡eq.
func ConditionalMove(ctx *build.Context, dst, src ir.Int, f ir.Operand, eq ir.Flag) {
	k := ints.Min(dst.Len(), src.Len())
	for i := 0; i < k; i++ {
		ctx.CMOV(src.Limb(i), dst.Limb(i), f, eq)
	}
}

// AddInto sets z = x+y using carry register c.
func AddInto(ctx *build.Context, z, x, y ir.Int, c ir.Register) {
	var cin ir.Operand = ir.Flag(0)
	for i := 0; i < z.Len(); i++ {
		ctx.ADD(ir.Limb(x, i), ir.Limb(y, i), cin, z.Limb(i), c)
		cin = c
	}
}

// Add x and y, returning an integer with size the larger of x and y.
func Add(ctx *build.Context, x, y ir.Int) ir.Int {
	return add(ctx, x, y, ints.Max(x.Len(), y.Len()))
}

// AddFull adds x and y, returning an integer with one more limb than the larger
// of x and y.
func AddFull(ctx *build.Context, x, y ir.Int) ir.Int {
	return add(ctx, x, y, ints.Max(x.Len(), y.Len())+1)
}

// add x+y returning a k-limb result.
func add(ctx *build.Context, x, y ir.Int, k int) ir.Int {
	z := ctx.Int("sum", k)
	c := ctx.Register("c")
	AddInto(ctx, z, x, y, c)
	return z
}

// SubInto sets z = x-y using borrow register b.
func SubInto(ctx *build.Context, z, x, y ir.Int, b ir.Register) {
	var bin ir.Operand = ir.Flag(0)
	for i := 0; i < z.Len(); i++ {
		ctx.SUB(ir.Limb(x, i), ir.Limb(y, i), bin, z.Limb(i), b)
		bin = b
	}
}

// MulInto sets z = x*y.
func MulInto(ctx *build.Context, z, x, y ir.Int) {
	acc := NewAccumulator(ctx, z)
	for i, a := range x.Limbs() {
		for j, b := range y.Limbs() {
			lo, hi := ctx.Register("lo"), ctx.Register("hi")
			ctx.MUL(a, b, hi, lo)
			acc.AddProduct(hi, lo, i+j)
		}
		acc.Flush()
	}
}
