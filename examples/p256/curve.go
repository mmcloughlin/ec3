package p256

import (
	"crypto/elliptic"
	"math/big"
)

// P256 returns a Curve which implements P-256.
func P256() elliptic.Curve { return p256 }

type curve struct{ *elliptic.CurveParams }

var p256 curve

func init() {
	p256.CurveParams = &elliptic.CurveParams{Name: "P-256"}
	p256.P, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	p256.N, _ = new(big.Int).SetString("115792089210356248762697446949407573529996955224135760342422259061068512044369", 10)
	p256.B, _ = new(big.Int).SetString("5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b", 16)
	p256.Gx, _ = new(big.Int).SetString("6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296", 16)
	p256.Gy, _ = new(big.Int).SetString("4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5", 16)
	p256.BitSize = 256
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
