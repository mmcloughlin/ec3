package verif

import (
	"math/big"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/z3"
)

// Add returns a specification for addition on multi-precision k*s-bit integers.
func Add(ctx *z3.Context, s, k uint) *Spec {
	return NewSpecBinary(ctx, s, k, (*z3.BV).Add)
}

// Sub returns a specification for subtraction on multi-precision k*s-bit integers.
func Sub(ctx *z3.Context, s, k uint) *Spec {
	return NewSpecBinary(ctx, s, k, (*z3.BV).Sub)
}

// MulFull returns a specification for a full multiply on multi-precision k*s-bit integers.
func MulFull(ctx *z3.Context, s, k uint) *Spec {
	bits := s * k
	return NewSpecBinarySizes(ctx, s, k, 2*k, func(x, y *z3.BV) *z3.BV {
		xext := x.ZeroExt(bits)
		yext := y.ZeroExt(bits)
		return xext.Mul(yext)
	})
}

// AddMod returns a specification for addition modulo mod on k*s-bit integers.
func AddMod(ctx *z3.Context, s, k uint, mod *big.Int) (*Spec, error) {
	f, err := NewField(ctx.BVSort(s*k), mod)
	if err != nil {
		return nil, err
	}

	spec := NewSpecBinary(ctx, s, k, f.Add)

	m := f.Modulus()
	for _, name := range []string{"x", "y"} {
		x := must(spec.Param(name))
		spec.AddPrecondition(x.ULT(m))
	}

	return spec, nil
}

// SubMod returns a specification for subtraction modulo mod on k*s-bit integers.
func SubMod(ctx *z3.Context, s, k uint, mod *big.Int) (*Spec, error) {
	f, err := NewField(ctx.BVSort(s*k), mod)
	if err != nil {
		return nil, err
	}

	spec := NewSpecBinary(ctx, s, k, f.Sub)

	m := f.Modulus()
	for _, name := range []string{"x", "y"} {
		x := must(spec.Param(name))
		spec.AddPrecondition(x.ULT(m))
	}

	return spec, nil
}

// NewSpecBinary returns a specification for a binary operator on k*s-bit integers.
func NewSpecBinary(ctx *z3.Context, s, k uint, op func(x, y *z3.BV) *z3.BV) *Spec {
	return NewSpecBinarySizes(ctx, s, k, k, op)
}

// NewSpecBinarySizes returns a specification for a binary operator on kin*s-bit
// integers, returning kout*s-bit integers.
func NewSpecBinarySizes(ctx *z3.Context, s, kin, kout uint, op func(x, y *z3.BV) *z3.BV) *Spec {
	tin := ir.Integer{K: kin}
	tout := ir.Integer{K: kout}
	sig := &ir.Signature{
		Params:  ir.NewVars(tin, "x", "y"),
		Results: ir.NewVars(tout, "z"),
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
