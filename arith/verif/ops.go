package verif

import (
	"math/big"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/z3"
)

// Add returns a specification for addition on multi-precision k*s-bit integers.
func Add(ctx *z3.Context, s, k uint) *Spec {
	return NewBinarySpec(ctx, s, k, (*z3.BV).Add)
}

// Sub returns a specification for subtraction on multi-precision k*s-bit integers.
func Sub(ctx *z3.Context, s, k uint) *Spec {
	return NewBinarySpec(ctx, s, k, (*z3.BV).Sub)
}

// AddMod returns a specification for addition modulo mod on k*s-bit integers.
func AddMod(ctx *z3.Context, s, k uint, mod *big.Int) (*Spec, error) {
	f, err := NewField(ctx.BVSort(s*k), mod)
	if err != nil {
		return nil, err
	}
	return NewBinarySpec(ctx, s, k, f.Add), nil
}

// NewBinarySpec returns a specification for a binary operator on k*s-bit integers.
func NewBinarySpec(ctx *z3.Context, s, k uint, op func(x, y *z3.BV) *z3.BV) *Spec {
	t := ir.Integer{K: k}
	sig := &ir.Signature{
		Params:  ir.NewVars(t, "x", "y"),
		Results: ir.NewVars(t, "z"),
	}
	spec := NewSpec(ctx, sig, s)
	x := must(spec.Param("x"))
	y := must(spec.Param("y"))
	z := op(x, y)
	spec.SetResult("z", z)
	return spec
}

func must(x *z3.BV, err error) *z3.BV {
	if err != nil {
		panic(err)
	}
	return x
}
