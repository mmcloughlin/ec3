package fp

import (
	"github.com/mmcloughlin/avo/attr"
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"

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

	// TODO(mbm): helper for loading integer from memory
	xb := operand.Mem{Base: a.ctx.Load(a.ctx.Param("x"), a.ctx.GP64())}
	x := mp.NewIntLimb64(a.ctx, 4)
	for i := 0; i < 4; i++ {
		a.ctx.MOVQ(xb.Offset(8*i), x[i])
	}

	yb := operand.Mem{Base: a.ctx.Load(a.ctx.Param("y"), a.ctx.GP64())}
	y := mp.NewIntLimb64(a.ctx, 4)
	for i := 0; i < 4; i++ {
		a.ctx.MOVQ(yb.Offset(8*i), y[i])
	}

	a.field.Add(a.ctx, x, y)

	for i := 0; i < 4; i++ {
		a.ctx.MOVQ(x[i], xb.Offset(8*i))
	}

	a.ctx.RET()
}

func (a Asm) Mul() {
	a.ctx.Function("Mul")
	a.ctx.Attributes(attr.NOSPLIT)
	a.ctx.SignatureExpr("func(z, x, y *[32]byte)")

	// Load arguments.
	zb := operand.Mem{Base: a.ctx.Load(a.ctx.Param("z"), a.ctx.GP64())}
	z := mp.NewIntFromMem(zb, 4)

	xb := operand.Mem{Base: a.ctx.Load(a.ctx.Param("x"), a.ctx.GP64())}
	x := mp.NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: a.ctx.Load(a.ctx.Param("y"), a.ctx.GP64())}
	y := mp.NewIntFromMem(yb, 4)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	mb := a.ctx.AllocLocal(8 * 8)
	m := mp.NewIntFromMem(mb, 8)

	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")
	a.field.ReduceDouble(a.ctx, z, m)

	a.ctx.RET()
}
