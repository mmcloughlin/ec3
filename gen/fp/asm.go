package fp

import (
	"github.com/mmcloughlin/avo/attr"
	"github.com/mmcloughlin/avo/build"

	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
)

type Asm struct {
	cfg   Config
	field fp.Crandall
	ctx   *build.Context
}

func NewAsm(cfg Config) *Asm {
	return &Asm{
		cfg: cfg,
		field: fp.Crandall{
			P: cfg.Prime,
		},
		ctx: build.NewContext(),
	}
}

func (a Asm) Context() *build.Context {
	return a.ctx
}

func (a Asm) Add() {
	a.ctx.Function("Add")
	a.ctx.Attributes(attr.NOSPLIT)
	a.ctx.SignatureExpr("func(x, y *[32]byte)")

	// Load parameters.
	xp := mp.Param(a.ctx, "x", 4)
	yp := mp.Param(a.ctx, "y", 4)

	// Bring into registers.
	x := mp.Registers(a.ctx, xp)
	y := mp.Registers(a.ctx, yp)

	// Add.
	a.field.Add(a.ctx, x, y)

	// Write back to registers.
	mp.Copy(a.ctx, xp, x)

	a.ctx.RET()
}

func (a Asm) Mul() {
	a.ctx.Function("Mul")
	a.ctx.Attributes(attr.NOSPLIT)
	a.ctx.SignatureExpr("func(z, x, y *[32]byte)")

	// Load parameters.
	z := mp.Param(a.ctx, "z", 4)
	x := mp.Param(a.ctx, "x", 4)
	y := mp.Param(a.ctx, "y", 4)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	stack := a.ctx.AllocLocal(8 * 8)
	m := mp.NewIntFromMem(stack, 8)
	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")
	a.field.ReduceDouble(a.ctx, z, m)

	a.ctx.RET()
}
