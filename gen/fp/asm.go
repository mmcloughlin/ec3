package fp

import (
	"github.com/mmcloughlin/avo/attr"
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/gotypes"

	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
)

type Asm struct {
	cfg   Config
	field fp.Builder
	ctx   *build.Context
}

func NewAsm(cfg Config) *Asm {
	ctx := build.NewContext()
	return &Asm{
		cfg:   cfg,
		field: cfg.Field.Build(ctx),
		ctx:   ctx,
	}
}

func (a Asm) Context() *build.Context {
	return a.ctx
}

func (a Asm) Function(name string, params ...string) {
	a.ctx.Function(name)
	a.ctx.Attributes(attr.NOSPLIT)

	sig := a.cfg.Signature(params...)
	a.ctx.Signature(gotypes.NewSignature(nil, sig))
}

func (a Asm) Add() {
	a.Function("Add", "x", "y")

	// Load parameters.
	xp := mp.Param(a.ctx, "x", a.field.Limbs())
	yp := mp.Param(a.ctx, "y", a.field.Limbs())

	// Bring into registers.
	x := mp.Registers(a.ctx, xp)
	y := mp.Registers(a.ctx, yp)

	// Add.
	a.field.Add(x, y)

	// Write back to registers.
	mp.Copy(a.ctx, xp, x)

	a.ctx.RET()
}

func (a Asm) Mul() {
	a.Function("Mul", "z", "x", "y")
	k := a.field.Limbs()

	// Load parameters.
	z := mp.Param(a.ctx, "z", k)
	x := mp.Param(a.ctx, "x", k)
	y := mp.Param(a.ctx, "y", k)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	stack := a.ctx.AllocLocal(8 * 2 * k)
	m := mp.NewIntFromMem(stack, 2*k)
	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")
	a.field.ReduceDouble(z, m)

	a.ctx.RET()
}
