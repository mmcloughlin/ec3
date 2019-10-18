package verif

import (
	"math/big"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/internal/z3"
)

type Field struct {
	n   uint
	mod *big.Int

	ctx *z3.Context
	bv  *z3.BVSort
	m   *z3.BV
}

func NewField(ctx *z3.Context, n uint, mod *big.Int) (*Field, error) {
	if mod.BitLen() > int(n) {
		return nil, xerrors.Errorf("modulus larger than %d bits", n)
	}
	bv := ctx.BVSort(n)
	return &Field{
		n:   n,
		mod: mod,

		ctx: ctx,
		bv:  bv,
		m:   bv.Int(mod),
	}, nil
}

// Var builds a named variable.
func (f *Field) Var(name string) *z3.BV {
	return f.bv.Const(name)
}

// Limbs breaks x into k-bit limbs.
func (f *Field) Limbs(x *z3.BV, k uint) []*z3.BV {
	limbs := []*z3.BV{}
	for l := uint(0); l < f.n; l += k {
		limb := x.Extract(l+k-1, l)
		limbs = append(limbs, limb)
	}
	return limbs
}

// FromLimbs builds a full integer from component limbs.
func (f *Field) FromLimbs(limbs []*z3.BV) *z3.BV {
	x := limbs[0]
	for i := 1; i < len(limbs); i++ {
		x = limbs[i].Concat(x)
	}
	return x
}

// Add returns an expression for the addition of x and y in the field.
func (f *Field) Add(x, y *z3.BV) *z3.BV {
	xext := x.ZeroExt(1)
	yext := y.ZeroExt(1)
	sum := xext.Add(yext)
	mext := f.m.ZeroExt(1)
	reduced := sum.Urem(mext)
	return reduced.Extract(f.n-1, 0)
}
