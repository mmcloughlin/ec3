// CodeGenerationWarning

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

// Double returns 2*(x1,y1)
func (c curve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	j1 := NewFromAffine(a1)
	d := new(Jacobian)
	d.Double(j1)
	return d.Affine().Coordinates()
}

// ScalarMult returns k*(x1,y1) where k is a number in big-endian form.
func (c curve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	// Implementation follows [msrecclibpaper] Algorithm 1.

	a := NewAffine(x1, y1)
	P := NewFromAffine(a)

	// Step 1: scalar validation.
	// TODO(mbm): exit if scalar is 0.

	// 	const omega = uint(5)
	// 	k = c.reduceScalar(k)

	var K scalar
	K.SetBytes(k)

	// Step 5: odd = k mod 2
	// Step 6: if odd = 0 then k = r − k

	// 	oddK, isEvenK := c.toOdd(k)

	even := K.ConvertToOdd()

	// Step 7: Recode k to (kt, . . . , k0) using Algorithm 6.

	//
	// 	var scalar big.Int
	// 	scalar.SetBytes(oddK)
	// 	if scalar.Sign() ≡ 0 {
	// 		return new(big.Int), new(big.Int)
	// 	}
	// 	L := math.SignedDigit(&scalar, omega)

	digits := K.FixedWindowRecode()

	// Step 4: Compute P[i] = (2i + 1)P for 0 <= i < 2^{w−2}.
	const TableSize = 1 << (ConstW - 1)
	var T [TableSize]Jacobian
	T[0] = P
	_2P := P
	_2P.Double()

	// 	var Q, R jacobianPoint
	// 	TabP := newAffinePoint(Px, Py).oddMultiples(omega)
	// 	for i := len(L) - 1; i ⩾ 0; i-- {
	// 		for j := uint(0); j < omega-1; j++ {
	// 			Q.double()
	// 		}
	// 		idx := absolute(L[i]) >> 1
	// 		for j := range TabP {
	// 			R.cmov(&TabP[j], subtle.ConstantTimeEq(int32(j), idx))
	// 		}
	// 		R.cneg(int(L[i]>>31) & 1)
	// 		Q.add(&Q, &R)
	// 	}
	// 	Q.cneg(isEvenK)
	// 	return Q.toAffine().toInt()
}

// 1. if k = 0 ∨ k ≥ r then return (“error: invalid scalar”) [if: validation]
// 2. Run point validation and compute T = 4P (for Ed) using Algorithm 2 for Eb and Algorithm 3 for Ed. If “invalid” return
// (“error: invalid point”), else set P = T (for Ed). [if: validation]
// Precomputation Stage:
// 3. Fix the window width 2 ≤ w < 10 ∈ Z+.
// Recoding Stage:
// (r)/(w − 1)e and sj are the signs of
// the recoded digits.
// Evaluation Stage:
// 8. Q = stP[(|kt| − 1)/2]
// 9. for i = (t − 1) to 1
// 10. if DBLADD = true ∧ w 6= 2 then [if: algorithm variant]
// 11. Q = 2(w−2)Q (Use Alg.10)
// 12. Q = 2Q + siP[(|ki| − 1)/2] (Use Alg.11)
// 13. else
// 14. Q = 2(w−1)Q (Use Alg.10 for Eb and Alg.14 for Ed)
// 15. Q = Q + siP[(|ki| − 1)/2] (Use Alg.12 for Eb and Alg.15 for Ed)
// 16. end if
// 17. end for
// 18. Q = 2(w−1)Q (Use Alg.10 for Eb and Alg.14 for Ed)
// 19. Q = Q ⊕ s0P[(|k0| − 1)/2] (Use Alg.19 for Eb and Alg.17 for Ed)
// 20. if odd = 0 then Q = −Q [if: masked constant time]
// 21. Convert Q to affine coordinates (x, y).
// 22. return Q.
