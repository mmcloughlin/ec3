package mp

import (
	"strconv"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/ir"
)

// Int represents a multi-precision integer.
type Int []ir.Register

func NewInt(prefix string, k uint) Int {
	x := make([]ir.Register, k)
	for i := 0; i < int(k); i++ {
		x[i] = ir.Register(prefix + strconv.Itoa(i))
	}
	return x
}

// Add sets z = x+y using carry register c.
func Add(ctx *build.Context, z, x, y Int, c ir.Register) {
	ctx.ADD(x[0], y[0], ir.Flag(0), z[0], c)
	for i := 1; i < len(x); i++ {
		ctx.ADD(x[i], y[i], c, z[i], c)
	}
}

// Sub sets z = x-y using borrow register b.
func Sub(ctx *build.Context, z, x, y Int, b ir.Register) {
	ctx.SUB(x[0], y[0], ir.Flag(0), z[0], b)
	for i := 1; i < len(x); i++ {
		ctx.SUB(x[i], y[i], b, z[i], b)
	}
}
