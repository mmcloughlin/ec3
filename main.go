package main

import (
	"fmt"
	"strconv"

	"github.com/mmcloughlin/avo/reg"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
)

// Crandall represents a prime of the form 2ⁿ - c. Named after Richard E. Crandall.
type Crandall struct {
	N int
	C int
}

// Bits returns the number of bits required to represent p.
func (p Crandall) Bits() int {
	return p.N
}

func (p Crandall) String() string {
	return fmt.Sprintf("2^%d - %d", p.N, p.C)
}

func (p Crandall) Slug() string {
	return strconv.Itoa(p.N) + strconv.Itoa(p.C)
}

// BitsToQuadWords returns the number of 64-bit quad-words required to hold bits.
func BitsToQuadWords(bits int) int {
	return (bits + 63) / 64
}

// NextMultiple returns the next multiple of n greater than or equal to a.
func NextMultiple(a, n int) int {
	a += n - 1
	a -= a % n
	return a - (a % n)
}

// Int represents a multi-precision integer.
type Int []operand.Op

// NewIntLimb64 builds multi-precision integer with k 64-bit limbs.
func NewIntLimb64(k int) Int {
	x := make(Int, k)
	for i := 0; i < k; i++ {
		x[i] = build.GP64()
	}
	return x
}

// Zero64 returns a 64-bit register initialized to zero.
func Zero64() reg.Register {
	zero := build.GP64()
	build.XORQ(zero, zero)
	return zero
}

// Add y into x modulo p.
//	x = x + y (mod p)
func Add(x, y Int, p Crandall) {
	n := p.Bits()
	l := NextMultiple(n, 64)

	// Add y into x.
	for i := 0; i < l; i++ {
		build.ADCXQ(y[i], x[i])
	}

	// Both inputs are < 2ˡ so the result is < 2^{l+1}.
	// If the last addition caused a carry into the l'th bit we need to perform a reduction.
	// Note that for the Crandall prime we have
	//
	//	2ⁿ - c = 0 (mod p)
	//	2ⁿ = c (mod p)
	//
	// However n may not be on a limb boundary, so we actually need the identity
	//
	//	2ˡ = 2^{l-n} * c (mod p)

	d := (1 << uint(l-n)) * p.C

	// The addend may be zero or d depending on the carry.
	// Initialize to zero and conditionally move d into it.
	addend := Zero64()
	dreg := build.GP64()
	build.MOVQ(operand.U16(d), dreg) // TODO(mbm): is U16 right?
	build.CMOVQCS(dreg, addend)

	// Now add the addend into x.
	build.ADCXQ(addend, x[0])
	zero := Zero64()
	for i := 1; i < l; i++ {
		build.ADCXQ(zero, x[i])
	}

	// We have added d into the low l bits. Therefore the result is less than 2ˡ + d.
	// But note that it could still be 2ˡ or higher, so we need to perform a
	// second reduction.

	// As before, the addend is either 0 or d depending on the carry from the last add.
	addend = Zero64()
	build.CMOVQCS(dreg, addend)

	// This time we only need to perform one add. The result must be less than 2ˡ + 2*d,
	// therefore provided 2*d does not exceed the size of a limb we can be sure there
	// will be no carry.
	build.ADCXQ(addend, x[0])
}

func main() {
	p := Crandall{N: 255, C: 19}
	fmt.Println(p)
}
