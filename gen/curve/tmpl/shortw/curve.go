// CodeGenerationWarning

package shortw

import (
	"crypto/elliptic"
	"math/big"
)

// CURVENAME returns a Curve which implements CanonicalName.
func CURVENAME() elliptic.Curve { return curvename }

type curve struct{ *elliptic.CurveParams }

var curvename = curve{
	CurveParams: &elliptic.CurveParams{Name: ConstCanonicalName},
}

func init() {
	curvename.P, _ = new(big.Int).SetString(ConstPDecimal, 10)
	curvename.N, _ = new(big.Int).SetString(ConstNDecimal, 10)
	curvename.B, _ = new(big.Int).SetString(ConstBHex, 16)
	curvename.Gx, _ = new(big.Int).SetString(ConstGxHex, 16)
	curvename.Gy, _ = new(big.Int).SetString(ConstGyHex, 16)
	curvename.BitSize = ConstBitSize
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (c curve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	a2 := NewAffine(x2, y2)
	j1 := NewFromAffine(a1)
	j2 := NewFromAffine(a2)
	s := new(Jacobian)
	s.Add(j1, j2)
	return s.Affine().Coordinates()
}

// Double returns 2*(x,y)
func (c curve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	j1 := NewFromAffine(a1)
	d := new(Jacobian)
	d.Double(j1)
	return d.Affine().Coordinates()
}
