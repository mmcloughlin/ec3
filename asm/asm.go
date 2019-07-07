package asm

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/reg"
)

// Zero64 returns a 64-bit register initialized to zero.
func Zero64(ctx *build.Context) reg.Register {
	zero := ctx.GP64()
	ctx.XORQ(zero, zero)
	return zero
}
