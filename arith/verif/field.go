package verif

import (
	"math/big"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/z3"
)

type Field struct {
	*z3.BVSort

	m *z3.BV
}

func NewField(sort *z3.BVSort, mod *big.Int) (*Field, error) {
	if mod.BitLen() > int(sort.Bits()) {
		return nil, xerrors.Errorf("modulus larger than %d bits", sort.Bits())
	}
	return &Field{
		BVSort: sort,
		m:      sort.Int(mod),
	}, nil
}

// Modulus returns the field modulus.
func (f *Field) Modulus() *z3.BV { return f.m }

// Add returns an expression for the addition of x and y in the field.
func (f *Field) Add(x, y *z3.BV) *z3.BV {
	xext := x.ZeroExt(1)
	yext := y.ZeroExt(1)
	sum := xext.Add(yext)
	mext := f.m.ZeroExt(1)
	reduced := sum.Urem(mext)
	return reduced.Extract(f.Bits()-1, 0)
}
