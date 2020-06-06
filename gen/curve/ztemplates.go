// Code generated by assets. DO NOT EDIT.

package curve

import "fmt"

func loadtemplate(name string) ([]byte, error) {
	switch name {
	case "tmpl/shortw/curve.go":
		return []byte(`// CodeGenerationWarning

package shortw

import (
	"crypto/elliptic"
	"math/big"
)

// References:
//
//	[msrecclibpaper]  Joppe W. Bos, Craig Costello, Patrick Longa and Michael Naehrig. Selecting
//	                  Elliptic Curves for Cryptography: An Efficiency and Security Analysis.
//	                  Cryptology ePrint Archive, Report 2014/130. 2014.
//	                  https://eprint.iacr.org/2014/130

// Curve extends the standard elliptic.Curve interface.
type Curve interface {
	elliptic.Curve

	// Inverse computes the inverse of k modulo the order N. Satisfies the
	// crypto/ecdsa.invertable interface.
	Inverse(k *big.Int) *big.Int
}

// CURVENAME returns a Curve which implements CanonicalName.
func CURVENAME() Curve { return curvename }

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
	p1 := a1.Projective()
	p2 := a2.Projective()
	s := new(Projective)
	s.CompleteAdd(p1, p2)
	return s.Affine().Coordinates()
}

// Double returns 2*(x1,y1)
func (c curve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	j1 := a1.Jacobian()
	d := new(Jacobian)
	d.Double(j1)
	return d.Affine().Coordinates()
}

// ScalarMult returns k*(x1,y1) where k is a number in big-endian form.
func (c curve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	// Implementation follows [msrecclibpaper] Algorithm 1.

	// Scalar recoding window size.
	const w = ConstW

	// Convert point from affine.
	a := NewAffine(x1, y1)
	p := a.Jacobian()

	// Step 1: scalar validation.

	// TODO(mbm): exit if scalar is 0.

	var K scalar
	K.SetBytesRaw(k)

	// Step 5: odd = k mod 2
	// Step 6: if odd = 0 then k = r − k
	even := K.ConvertToOdd()

	// Step 7: Recode k to (k_t, ..., k₀) using Algorithm 6.
	digits := K.FixedWindowRecode()

	// Step 4: Compute P[i] = (2i + 1)P for 0 ⩽ i < 2^{w−2}.
	var tbl table
	tbl.Precompute(p)

	// Step 8: Q = s_t * P[(|k_t| − 1)/2]
	var q, r Jacobian

	t := len(digits) - 1
	tbl.Lookup(&q, digits[t])

	// Step 9: for i = (t − 1) to 1
	for i := t - 1; i >= 1; i-- {
		// Step 14: Q = 2^{w−1}Q
		for j := 0; j < w-1; j++ {
			q.Double(&q)
		}

		// Step 15: Q = Q + s_i * P[(|k_i| − 1)/2]
		tbl.Lookup(&r, digits[i])
		q.Add(&q, &r)
	}

	// Step 18: Q = 2^{w−1}Q
	for j := 0; j < w-1; j++ {
		q.Double(&q)
	}

	// Step 19: Q = Q ⊕ s₀ * P[(|k₀| − 1)/2]
	tbl.Lookup(&r, digits[0])
	rp := r.Projective()
	qp := q.Projective()
	qp.CompleteAdd(qp, rp)

	// Step 20: if odd = 0 then Q = −Q
	qp.CNeg(even)

	// Step 21: Convert Q to affine coordinates (x, y).
	return qp.Affine().Coordinates()
}

// tablesize is the size of the lookup table used by ScalarMult.
const tablesize = 1 << (ConstW - 1)

// table is a lookup table used by ScalarMult.
type table [tablesize]Jacobian

// Precompute odd multiples of p.
func (t *table) Precompute(p *Jacobian) {
	t[0].Set(p)

	var _2p Jacobian
	_2p.Double(p)

	for i := 1; i < tablesize; i++ {
		t[i].Add(&t[i-1], &_2p)
	}
}

// Lookup an odd multiple from the table and store it in p. The provided digit
// may be negative, in which case p will be negated.
func (t *table) Lookup(p *Jacobian, digit int32) {
	idx := abs(digit) / 2
	lookup(p, t[:], int(idx))
	p.CNeg(sign(digit))
}

// abs returns the absolute value of x.
func abs(x int32) int32 {
	mask := x >> 31
	return (x + mask) ^ mask
}

// sign returns the sign bit of x.
func sign(x int32) uint {
	return uint(x>>31) & 1
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (c curve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return c.ScalarMult(c.Gx, c.Gy, k)
}

// Inverse computes the inverse of k modulo the order N. Satisfies the
// crypto/ecdsa.invertable interface.
func (curve) Inverse(k *big.Int) *big.Int {
	var (
		K   scalar
		inv scalar
	)

	K.SetInt(k)
	scalarinv(&inv, &K)
	return inv.Int()
}
`), nil

	case "tmpl/shortw/curve_test.go":
		return []byte(`// CodeGenerationWarning

package shortw

import (
	"math/big"
	"testing"
)

var (
	cur = CURVENAME()
	ref = curvename.CurveParams
)

func TestCurveAddRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		gx, gy := cur.Add(x1, y1, x2, y2)
		ex, ey := ref.Add(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveAddAsDouble(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		gx, gy := cur.Add(x, y, x, y)
		ex, ey := ref.Double(x, y)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveAddNegative(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		nx := new(big.Int).Set(x)
		ny := new(big.Int).Neg(y)

		gx, gy := cur.Add(x, y, nx, ny)
		zero := new(big.Int)

		EqualInt(t, "x", zero, gx)
		EqualInt(t, "y", zero, gy)
	}
}

func TestCurveDoubleRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		gx, gy := ref.Double(x, y)
		ex, ey := cur.Double(x, y)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarMultRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)
		x, y := RandPoint(t)

		gx, gy := cur.ScalarMult(x, y, k.Bytes())
		ex, ey := ref.ScalarMult(x, y, k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarBaseMultRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		gx, gy := cur.ScalarBaseMult(k.Bytes())
		ex, ey := ref.ScalarBaseMult(k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveInverseRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		got := cur.Inverse(k)

		expect := new(big.Int).Set(k)
		expect.ModInverse(expect, ref.N)

		EqualInt(t, "inv", expect, got)
	}
}

func BenchmarkScalarMult(b *testing.B) {
	x, y := RandPoint(b)
	K := RandScalarNonZero(b)
	k := K.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cur.ScalarMult(x, y, k)
	}
}

func BenchmarkScalarBaseMult(b *testing.B) {
	K := RandScalarNonZero(b)
	k := K.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cur.ScalarBaseMult(k)
	}
}
`), nil

	case "tmpl/shortw/recode.go":
		return []byte(`// CodeGenerationWarning

package shortw

import (
	"math/bits"
	"unsafe"
)

// References:
//
//	[msrecclibcode]   Microsoft Research. MSR Elliptic Curve Cryptography Library. 2014.
//	                  https://www.microsoft.com/en-us/research/project/msr-elliptic-curve-cryptography-library/
//	[msrecclibpaper]  Joppe W. Bos, Craig Costello, Patrick Longa and Michael Naehrig. Selecting
//	                  Elliptic Curves for Cryptography: An Efficiency and Security Analysis.
//	                  Cryptology ePrint Archive, Report 2014/130. 2014.
//	                  https://eprint.iacr.org/2014/130

const words = scalarsize / 8

// uint64s provides a view of k as an array of uint64 words.
func (k *scalar) uint64s() *[words]uint64 {
	return (*[words]uint64)(unsafe.Pointer(k))
}

// FixedWindowRecode recodes the odd scalar k into a signed fixed window
// representation with digits in the set {±1, ±3, ..., ±(2^(w-1)-1)}.
func (k *scalar) FixedWindowRecode() []int32 {
	// Implementation follows [msrecclibpaper] Algorithm 6.
	const (
		w    = ConstW                  // window parameter
		r    = ConstBitSize            // bit size
		t    = (r + (w - 2)) / (w - 1) // length of the window representation
		mask = (1 << w) - 1            // w-bit mask
		val  = 1 << (w - 1)            // 2ʷ⁻¹
	)

	digits := make([]int32, t+1)
	K := *k

	// Step 2: for i = 0 to (t-1)
	for i := 0; i < t; i++ {
		// Step 3: k_i = ( k mod 2ʷ ) - 2ʷ⁻¹
		digits[i] = int32(K[0]&mask) - val

		// Step 4: k = (k - k_i) / 2ʷ⁻¹
		K.SubInt32(digits[i])
		K.Rsh(w - 1)
	}

	// Step 5: k_t = k
	digits[t] = int32(K[0])

	return digits
}

// ConvertToOdd negates k if it is even. Returns whether the scalar was even.
func (k *scalar) ConvertToOdd() (even uint) {
	even = uint(k[0]&1) ^ 1
	var n scalar
	scalarneg(&n, k)
	scalarcmov(k, &n, even)
	return
}

// SubInt32 subtracts signed integer v from k.
func (k *scalar) SubInt32(v int32) {
	kw := k.uint64s()
	uv := uint64(v)
	var borrow uint64
	kw[0], borrow = bits.Sub64(kw[0], uv, 0)
	borrow &= (uv >> 63) ^ 1
	for i := 1; i < words; i++ {
		kw[i], borrow = bits.Sub64(kw[i], 0, borrow)
	}
}

// Rsh shifts the scalar k right by s.
func (k *scalar) Rsh(s uint) {
	kw := k.uint64s()
	for i := 0; i+1 < words; i++ {
		kw[i] = (kw[i] >> s) | (kw[i+1] << (64 - s))
	}
	kw[words-1] >>= s
}
`), nil

	case "tmpl/shortw/recode_test.go":
		return []byte(`// CodeGenerationWarning

package shortw

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestScalarFixedWindowRecode(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandOddScalar(t)

		var K scalar
		K.SetIntRaw(k)
		digits := K.FixedWindowRecode()

		// Verify all digits are odd.
		for i, digit := range digits {
			if (digit & 1) != 1 {
				t.Fatalf("digit %d is not odd", i)
			}
		}

		// Confirm the sum is correct.
		x := new(big.Int)
		for i := len(digits) - 1; i >= 0; i-- {
			x.Lsh(x, ConstW-1)
			x.Add(x, big.NewInt(int64(digits[i])))
		}

		if k.Cmp(x) != 0 {
			t.Logf("     k = %x", k)
			t.Logf("digits = %d", digits)
			t.Logf("   got = %x", x)
			t.FailNow()
		}
	}
}

func TestScalarSubInt(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x := RandScalarNonZero(t)
		v := rand.Int31n(64) - 32

		// Compute subtraction via scalar type.
		var k scalar
		k.SetIntRaw(x)
		k.SubInt32(v)
		got := k.IntRaw()

		// Compute expectation.
		expect := new(big.Int).Sub(x, new(big.Int).SetInt64(int64(v)))

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}

func TestScalarRsh(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x := RandScalarNonZero(t)
		s := uint(1 + rand.Intn(63))

		// Compute shift via scalar type.
		var k scalar
		k.SetIntRaw(x)
		k.Rsh(s)
		got := k.IntRaw()

		// Compute expectation.
		expect := new(big.Int).Rsh(x, s)

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}
`), nil

	case "tmpl/shortw/stubs.go":
		return []byte(`package shortw

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
	X, Y big.Int
}

func NewAffine(x, y *big.Int) *Affine {
	a := new(Affine)
	a.X.Set(x)
	a.Y.Set(y)
	return a
}

func (p *Affine) Set(q *Affine) {
	p.X.Set(&q.X)
	p.Y.Set(&q.Y)
}

func (p *Affine) Coordinates() (X, Y *big.Int) {
	return new(big.Int).Set(&p.X), new(big.Int).Set(&p.Y)
}

func (a *Affine) Jacobian() *Jacobian {
	j := &Jacobian{}
	j.a.Set(a)
	return j
}

func (a *Affine) Projective() *Projective {
	p := &Projective{}
	p.a.Set(a)
	return p
}

// Jacobian is a stub jacobian point type.
type Jacobian struct {
	a Affine
}

func (p *Jacobian) Set(q *Jacobian) {
	p.a.Set(&q.a)
}

func (p *Jacobian) Affine() *Affine {
	return &p.a
}

func (p *Jacobian) CMov(q *Jacobian, c uint) {
	if c != 0 {
		p.Set(q)
	}
}

func (p *Jacobian) CNeg(c uint) {
	if c != 0 {
		p.Neg()
	}
}

func (p *Jacobian) Neg() {
	y := new(big.Int).Neg(&p.a.Y)
	y.Mod(y, curvename.P)
	p.a.Y.Set(y)
}

func (p *Jacobian) Add(q, r *Jacobian) {
	x, y := curvename.Params().Add(&q.a.X, &q.a.Y, &r.a.X, &r.a.Y)
	p.a.X.Set(x)
	p.a.Y.Set(y)
}

func (p *Jacobian) Double(q *Jacobian) {
	x, y := curvename.Params().Double(&q.a.X, &q.a.Y)
	p.a.X.Set(x)
	p.a.Y.Set(y)
}

func (j *Jacobian) Projective() *Projective {
	p := &Projective{}
	p.a.Set(&j.a)
	return p
}

// Projective is a stub projective point type.
type Projective struct {
	a Affine
}

func (p *Projective) Affine() *Affine {
	return &p.a
}

func (p *Projective) CNeg(c uint) {
	if c != 0 {
		p.Neg()
	}
}

func (p *Projective) Neg() {
	y := new(big.Int).Neg(&p.a.Y)
	y.Mod(y, curvename.P)
	p.a.Y.Set(y)
}

func (p *Projective) CompleteAdd(q, r *Projective) {
	x, y := curvename.Params().Add(&q.a.X, &q.a.Y, &r.a.X, &r.a.Y)
	p.a.X.Set(x)
	p.a.Y.Set(y)
}

// lookup position idx in tbl.
func lookup(p *Jacobian, tbl []Jacobian, idx int) {
	p.Set(&tbl[idx])
}

// scalarsize is the size of a scalar field element in bytes.
const scalarsize = ConstBitSize / 8

// scalar is a stub scalar field element type.
type scalar [scalarsize]byte

// TODO(mbm): use ec3 itself to codegen the scalar type stub
// This will reduce duplication and help retain compatibility.

func (k *scalar) SetInt(x *big.Int) {
	k.SetIntRaw(x)
}

func (k *scalar) SetIntRaw(x *big.Int) {
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
	return k.IntRaw()
}

// IntRaw converts to a big integer.
func (k *scalar) IntRaw() *big.Int {
	// Endianness swap.
	var be scalar
	for i := 0; i < scalarsize; i++ {
		be[scalarsize-1-i] = k[i]
	}
	// Build big.Int.
	return new(big.Int).SetBytes(be[:])
}

// SetBytesRaw interprets b as the bytes of a big-endian unsigned integer, and sets k to that value.
func (k *scalar) SetBytesRaw(b []byte) {
	k.SetIntRaw(new(big.Int).SetBytes(b))
}

func scalarcmov(z, x *scalar, c uint) {
	if c != 0 {
		*z = *x
	}
}

func scalarneg(z, x *scalar) {
	neg := new(big.Int).Neg(x.IntRaw())
	neg.Mod(neg, curvename.N)
	z.SetIntRaw(neg)
}

func scalarinv(z, x *scalar) {
	inv := new(big.Int).ModInverse(x.Int(), curvename.N)
	z.SetInt(inv)
}
`), nil

	case "tmpl/shortw/stubs_test.go":
		return []byte(`package shortw

const ConstNumTrials = 128
`), nil

	case "tmpl/shortw/util_test.go":
		return []byte(`// CodeGenerationWarning

package shortw

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func RandScalarNonZero(tb testing.TB) *big.Int {
	tb.Helper()
	N := curvename.Params().N
	for {
		k, err := rand.Int(rand.Reader, N)
		if err != nil {
			tb.Fatal(err)
		}
		if k.Sign() == 0 {
			continue
		}
		return k
	}
}

func RandOddScalar(tb testing.TB) *big.Int {
	tb.Helper()
	k := RandScalarNonZero(tb)
	N := curvename.Params().N
	if k.Bit(0) == 0 {
		k.Neg(k).Mod(k, N)
	}
	return k
}

func RandPoint(tb testing.TB) (x, y *big.Int) {
	tb.Helper()
	k := RandScalarNonZero(tb)
	return curvename.Params().ScalarBaseMult(k.Bytes())
}

func EqualInt(t *testing.T, name string, expect, got *big.Int) {
	t.Helper()
	if got.Cmp(expect) != 0 {
		t.Logf("   got %x", got)
		t.Logf("expect %x", expect)
		t.Fatalf("%s: not equal", name)
	}
}
`), nil

	default:
		return nil, fmt.Errorf("unknown asset %s", name)
	}
}
