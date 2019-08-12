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
	stack := a.ctx.AllocLocal(8 * 2 * k)
	m := mp.NewIntFromMem(stack, 2*k)
	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")
	a.field.ReduceDouble(z, m)

	a.ctx.RET()
}

func (a Asm) Adhoc() {
	a.ctx.Function("Adhoc")
	a.ctx.SignatureExpr("func(z *[72]byte, x, y *[32]byte)")
	k := a.field.Limbs()

	// Load parameters.
	z := mp.Param(a.ctx, "z", 2*k+1)
	x := mp.Param(a.ctx, "x", k)
	y := mp.Param(a.ctx, "y", k)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	stack := a.ctx.AllocLocal(8 * 2 * k)
	m := mp.NewIntFromMem(stack, 2*k)
	mp.Mul(a.ctx, m, x, y)

	// Reduce.
	a.ctx.Comment("Reduction.")

	fld := a.field.(interface {
		MontReduce(z, x mp.Int)
	})
	fld.MontReduce(z, m)

	a.ctx.RET()
}
