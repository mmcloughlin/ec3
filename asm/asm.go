package asm

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/ir"
	"github.com/mmcloughlin/avo/pass"
	"github.com/mmcloughlin/avo/reg"
)

// Compile is a convenience for compiling generated assembly code.
func Compile(ctx *build.Context) (*ir.File, error) {
	f, err := ctx.Result()
	if err != nil {
		return nil, err
	}

	if err := pass.Compile.Execute(f); err != nil {
		return nil, err
	}

	return f, nil
}

// Zero64 returns a 64-bit register initialized to zero.
func Zero64(ctx *build.Context) reg.Register {
	zero := ctx.GP64()
	ctx.XORQ(zero, zero)
	return zero
}

// IsRegisterName reports whether name is a register name.
func IsRegisterName(name string) bool {
	for _, family := range reg.Families {
		for _, r := range family.Registers() {
			if r.Asm() == name {
				return true
			}
		}
	}
	return false
}
