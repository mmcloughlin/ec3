package shortw

import "math/big"

const (
	ConstCanonicalName = "Curve-Name"
	ConstPDecimal      = "39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319"
	ConstNDecimal      = "39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643"
	ConstBHex          = "b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef"
	ConstGxHex         = "aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7"
	ConstGyHex         = "3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f"
	ConstBitSize       = 384
)

type Affine struct {
	X, Y *big.Int
}

func NewAffine(x, y *big.Int) Affine {
	return Affine{X: x, Y: y}
}

func (p Affine) Coordinates() (X, Y *big.Int) {
	return p.X, p.Y
}

type Jacobian Affine

func NewFromAffine(a Affine) Jacobian {
	return Jacobian(a)
}

func (p Jacobian) Affine() Affine {
	return Affine(p)
}

func (p *Jacobian) Add(q, r Jacobian) {
	p.X, p.Y = curvename.Params().Add(q.X, q.Y, r.X, r.Y)
}

func (p *Jacobian) Double(q Jacobian) {
	p.X, p.Y = curvename.Params().Double(q.X, q.Y)
}
