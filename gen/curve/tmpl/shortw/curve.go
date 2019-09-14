// CodeGenerationWarning

package shortw

import (
	"crypto/elliptic"
	"crypto/subtle"
	"math/big"
)

// References:
//
//	[msrecclibpaper]  Joppe W. Bos, Craig Costello, Patrick Longa and Michael Naehrig. Selecting
//	                  Elliptic Curves for Cryptography: An Efficiency and Security Analysis.
//	                  Cryptology ePrint Archive, Report 2014/130. 2014.
//	                  https://eprint.iacr.org/2014/130

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
	j1 := a1.Jacobian()
	j2 := a2.Jacobian()
	s := new(Jacobian)
	s.Add(j1, j2)
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

	// Step 7: Recode k to (k_t, ..., k_0) using Algorithm 6.
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

	// Step 19: Q = Q ⊕ s_0 * P[(|k_0| − 1)/2]
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
	for i := range t {
		p.CMov(&t[i], uint(subtle.ConstantTimeEq(int32(i), idx)))
	}
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
