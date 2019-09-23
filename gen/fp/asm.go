package fp

import (
	"go/token"
	"go/types"

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
	a.ctx.Function(a.cfg.Name(name))
	a.ctx.Pragma("noescape")
	a.ctx.Attributes(attr.NOSPLIT)

	sig := a.cfg.Signature(params...)
	a.ctx.Signature(gotypes.NewSignature(nil, sig))
}

func (a Asm) CMov() {
	// Declare the function. We can't use the standard helper here since this
	// doesn't only take field element arguments.
	a.ctx.Function(a.cfg.Name("CMov"))
	a.ctx.Pragma("noescape")
	a.ctx.Attributes(attr.NOSPLIT)
	params := types.NewTuple(
		a.cfg.Param("y"),
		a.cfg.Param("x"),
		types.NewParam(token.NoPos, nil, "c", types.Typ[types.Uint]),
	)
	sig := types.NewSignature(nil, params, nil, false)
	a.ctx.Signature(gotypes.NewSignature(nil, sig))

	// Load parameters.
	yp := mp.Param(a.ctx, "y", a.field.Limbs())
	xp := mp.Param(a.ctx, "x", a.field.Limbs())
	c := a.ctx.Load(a.ctx.Param("c"), a.ctx.GP64())

	// Bring into registers.
	y := mp.CopyIntoRegisters(a.ctx, yp)
	x := mp.CopyIntoRegisters(a.ctx, xp)

	// Do the conditional move.
	mp.ConditionalMove(a.ctx, y, x, c)

	// Write back to memory.
	mp.Copy(a.ctx, yp, y)

	a.ctx.RET()
}

func (a Asm) Add() {
	a.Function("Add", "z", "x", "y")

	// Load parameters.
	zp := mp.Param(a.ctx, "z", a.field.Limbs())
	xp := mp.Param(a.ctx, "x", a.field.Limbs())
	yp := mp.Param(a.ctx, "y", a.field.Limbs())

	// Bring into registers.
	x := mp.CopyIntoRegisters(a.ctx, xp)
	y := mp.CopyIntoRegisters(a.ctx, yp)

	// Add.
	a.field.Add(x, y)

	// Write back to memory.
	mp.Copy(a.ctx, zp, x)

	a.ctx.RET()
}

func (a Asm) Sub() {
	a.Function("Sub", "z", "x", "y")

	// Load parameters.
	zp := mp.Param(a.ctx, "z", a.field.Limbs())
	xp := mp.Param(a.ctx, "x", a.field.Limbs())
	yp := mp.Param(a.ctx, "y", a.field.Limbs())

	// Bring into registers.
	x := mp.CopyIntoRegisters(a.ctx, xp)
	y := mp.CopyIntoRegisters(a.ctx, yp)

	// Subtract.
	a.field.Sub(x, y)

	// Write to z.
	mp.Copy(a.ctx, zp, x)

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
	m := mp.AllocLocal(a.ctx, 2*k)
	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")
	a.field.ReduceDouble(z, m)

	a.ctx.RET()
}
