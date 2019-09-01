package shortw

import (
	"math/big"
	"math/bits"
)

// References:
//
//	[msrecclibcode]   Microsoft Research. MSR Elliptic Curve Cryptography Library. 2014.
//	                  https://www.microsoft.com/en-us/research/project/msr-elliptic-curve-cryptography-library/
//	[msrecclibpaper]  Joppe W. Bos, Craig Costello, Patrick Longa and Michael Naehrig. Selecting
//	                  Elliptic Curves for Cryptography: An Efficiency and Security Analysis.
//	                  Cryptology ePrint Archive, Report 2014/130. 2014.
//	                  https://eprint.iacr.org/2014/130

const words = ConstBitSize / 64

type scalar [words]uint64

// SetInt sets the scalar k to the big integer value x.
func (k *scalar) SetInt(x *big.Int) {
	if x.Sign() < 0 || x.Cmp(curvename.N) >= 0 {
		x = new(big.Int).Mod(x, curvename.N)
	}

	for i := range k {
		k[i] = 0
	}

	for i, word := range x.Bits() {
		k[i] = uint64(word)
	}
}

// Int returns k as a big integer.
func (k *scalar) Int() *big.Int {
	x := new(big.Int)
	for i := words - 1; i >= 0; i-- {
		x.Lsh(x, 64)
		x.Add(x, new(big.Int).SetUint64(k[i]))
	}
	return x
}

// FixedWindowRecode recodes the odd scalar k into a signed fixed window
// representation with digits in the set {+-1, +-3, ..., +-(2^(w-1)-1)}.
func (k *scalar) FixedWindowRecode() []int {
	// Implementation follows [msrecclibpaper] Algorithm 6.
	const (
		w    = ConstW                  // window parameter
		r    = ConstBitSize            // bit size
		t    = (r + (w - 2)) / (w - 1) // length of the window representation
		mask = (1 << w) - 1            // w-bit mask
		val  = 1 << (w - 1)            // 2^{w-1}
	)

	digits := make([]int, t+1)

	// Step 2: for i = 0 to (t-1)
	for i := 0; i < t; i++ {
		// Step 3: k_i = ( k mod 2^w ) - 2^{w-1}
		v := int(k[0]&mask) - val
		digits[i] = v

		// Step 4: k = (k - k_i) / 2^{w-1}
		k.Sub(v)
		k.Rsh(w - 1)
	}

	// Step 5: k_t = k
	digits[t] = int(k[0])

	return digits
}

// Sub subtracts signed integer v from k.
func (k *scalar) Sub(v int) {
	uv := uint64(v)
	var borrow uint64
	k[0], borrow = bits.Sub64(k[0], uv, 0)
	borrow &= (uv >> 63) ^ 1
	for i := 1; i < words; i++ {
		k[i], borrow = bits.Sub64(k[i], 0, borrow)
	}
}

// Rsh shifts the scalar k right by s.
func (k *scalar) Rsh(s uint) {
	for i := 0; i+1 < words; i++ {
		k[i] = (k[i] >> s) | (k[i+1] << (64 - s))
	}
	k[words-1] >>= s
}
