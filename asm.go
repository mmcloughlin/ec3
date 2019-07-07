// +build ignore

package main

import (
	"strconv"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"

	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/prime"
)

// Slug returns a concise name for p.
func Slug(p prime.Crandall) string {
	return strconv.Itoa(p.N) + strconv.Itoa(p.C)
}

// BitsToQuadWords returns the number of 64-bit quad-words required to hold bits.
func BitsToQuadWords(bits int) int {
	return (bits + 63) / 64
}

// addmod builds a modular addition function.
func addmod(p prime.Crandall) {
	build.TEXT("Add"+Slug(p), build.NOSPLIT, "func(x, y *[32]byte)")

	// TODO(mbm): helper for loading integer from memory
	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := mp.NewIntLimb64(4)
	for i := 0; i < 4; i++ {
		build.MOVQ(xb.Offset(8*i), x[i])
	}

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := mp.NewIntLimb64(4)
	for i := 0; i < 4; i++ {
		build.MOVQ(yb.Offset(8*i), y[i])
	}

	fp.AddModP(x, y, p)

	for i := 0; i < 4; i++ {
		build.MOVQ(x[i], xb.Offset(8*i))
	}

	build.RET()
}

// mul builds a multiplication function.
func mul() {
	build.TEXT("Mul", build.NOSPLIT, "func(z *[64]byte, x, y *[32]byte)")

	zb := operand.Mem{Base: build.Load(build.Param("z"), build.GP64())}
	z := mp.NewIntFromMem(zb, 8)

	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := mp.NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := mp.NewIntFromMem(yb, 4)

	mp.Mul(z, x, y)

	build.RET()
}

// mulmod builds a modular multiplication function.
func mulmod(p prime.Crandall) {
	build.TEXT("Mul"+Slug(p), build.NOSPLIT, "func(z *[32]byte, x, y *[32]byte)")

	// Load arguments.
	zb := operand.Mem{Base: build.Load(build.Param("z"), build.GP64())}
	z := mp.NewIntFromMem(zb, 4)

	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := mp.NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := mp.NewIntFromMem(yb, 4)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	mb := build.AllocLocal(8 * 8)
	m := mp.NewIntFromMem(mb, 8)

	mp.Mul(m, x, y)

	// Reduce.
	build.Comment("Reduction.")
	fp.ReduceDouble(z, m, p)

	build.RET()
}

func main() {
	// Multi-precision.
	mul()

	// Fp25519
	p := prime.Crandall{N: 255, C: 19}
	addmod(p)
	mulmod(p)

	build.Generate()
}
