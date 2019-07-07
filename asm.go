// +build ignore

package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/mmcloughlin/avo/attr"
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
func addmod(ctx *build.Context, p prime.Crandall) {
	ctx.Function("Add" + Slug(p))
	ctx.Attributes(attr.NOSPLIT)
	ctx.SignatureExpr("func(x, y *[32]byte)")

	// TODO(mbm): helper for loading integer from memory
	xb := operand.Mem{Base: ctx.Load(ctx.Param("x"), ctx.GP64())}
	x := mp.NewIntLimb64(ctx, 4)
	for i := 0; i < 4; i++ {
		ctx.MOVQ(xb.Offset(8*i), x[i])
	}

	yb := operand.Mem{Base: ctx.Load(ctx.Param("y"), ctx.GP64())}
	y := mp.NewIntLimb64(ctx, 4)
	for i := 0; i < 4; i++ {
		ctx.MOVQ(yb.Offset(8*i), y[i])
	}

	fp.AddModP(ctx, x, y, p)

	for i := 0; i < 4; i++ {
		ctx.MOVQ(x[i], xb.Offset(8*i))
	}

	ctx.RET()
}

// mul builds a multiplication function.
func mul(ctx *build.Context) {
	ctx.Function("Mul")
	ctx.Attributes(attr.NOSPLIT)
	ctx.SignatureExpr("func(z *[64]byte, x, y *[32]byte)")

	zb := operand.Mem{Base: ctx.Load(ctx.Param("z"), ctx.GP64())}
	z := mp.NewIntFromMem(zb, 8)

	xb := operand.Mem{Base: ctx.Load(ctx.Param("x"), ctx.GP64())}
	x := mp.NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: ctx.Load(ctx.Param("y"), ctx.GP64())}
	y := mp.NewIntFromMem(yb, 4)

	mp.Mul(ctx, z, x, y)

	ctx.RET()
}

// mulmod builds a modular multiplication function.
func mulmod(ctx *build.Context, p prime.Crandall) {
	ctx.Function("Mul" + Slug(p))
	ctx.Attributes(attr.NOSPLIT)
	ctx.SignatureExpr("func(z, x, y *[32]byte)")

	// Load arguments.
	zb := operand.Mem{Base: ctx.Load(ctx.Param("z"), ctx.GP64())}
	z := mp.NewIntFromMem(zb, 4)

	xb := operand.Mem{Base: ctx.Load(ctx.Param("x"), ctx.GP64())}
	x := mp.NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: ctx.Load(ctx.Param("y"), ctx.GP64())}
	y := mp.NewIntFromMem(yb, 4)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	mb := ctx.AllocLocal(8 * 8)
	m := mp.NewIntFromMem(mb, 8)

	mp.Mul(ctx, m, x, y)

	// Reduce.
	ctx.Comment("Reduction.")
	fp.ReduceDouble(ctx, z, m, p)

	ctx.RET()
}

var (
	commandline = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags       = build.NewFlags(commandline)
)

func main() {
	ctx := build.NewContext()

	// Multi-precision.
	mul(ctx)

	// Fp25519
	p := prime.Crandall{N: 255, C: 19}
	addmod(ctx, p)
	mulmod(ctx, p)

	// Process the command-line.
	flag.Parse()
	cfg := flags.Config()
	status := build.Main(cfg, ctx)
	os.Exit(status)
}
