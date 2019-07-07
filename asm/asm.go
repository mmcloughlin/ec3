package asm

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/reg"
)

// Zero64 returns a 64-bit register initialized to zero.
func Zero64() reg.Register {
	zero := build.GP64()
	build.XORQ(zero, zero)
	return zero
}
