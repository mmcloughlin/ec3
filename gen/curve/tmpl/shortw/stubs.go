package shortw

import "math/big"

// Curve parameters.
const (
	ConstCanonicalName = "Curve-Name"
	ConstPDecimal      = "39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319"
	ConstNDecimal      = "39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643"
	ConstBHex          = "b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef"
	ConstGxHex         = "aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7"
	ConstGyHex         = "3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f"
	ConstBitSize       = 384
)

// Implementation parameters.
const (
	ConstW = 6
)

// Affine is a stub affine point type.
type Affine struct {
	X, Y *big.Int
}

func NewAffine(x, y *big.Int) Affine {
	return Affine{X: x, Y: y}
}

func (p Affine) Coordinates() (X, Y *big.Int) {
	return p.X, p.Y
}

// Jacobian is a stub jacobian point type.
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

// scalarsize is the size of a scalar field element in bytes.
const scalarsize = ConstBitSize / 8

// scalar is a stub scalar field element type.
type scalar [scalarsize]byte

// TODO(mbm): use ec3 itself to codegen the scalar type stub
// This will reduce duplication and help retain compatibility.

func (k *scalar) SetInt(x *big.Int) {
	if x.Sign() < 0 || x.Cmp(curvename.N) >= 0 {
		x = new(big.Int).Mod(x, curvename.N)
	}

	for i := range k {
		k[i] = 0
	}

	bs := x.Bytes()
	for i, b := range bs {
		k[len(bs)-1-i] = b
	}
}

// Int converts to a big integer.
func (k *scalar) Int() *big.Int {
	// Endianness swap.
	var be scalar
	for i := 0; i < scalarsize; i++ {
		be[scalarsize-1-i] = k[i]
	}
	// Build big.Int.
	return new(big.Int).SetBytes(be[:])
}
