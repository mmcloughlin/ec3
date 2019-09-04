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
