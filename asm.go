// +build ignore

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

// NewIntFromMem builds a multi-precision integer referencing the k 64-bit limbs
// at memory address m.
func NewIntFromMem(m operand.Mem, k int) Int {
	x := make(Int, k)
	for i := 0; i < k; i++ {
		x[i] = m.Offset(8 * i)
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
	k := l / 64

	// Prepare a zero register.
	zero := Zero64()

	// Note that for the Crandall prime we have
	//
	//	2ⁿ - c = 0 (mod p)
	//	2ⁿ = c (mod p)
	//
	// However n may not be on a limb boundary, so we actually need the identity
	//
	//	2ˡ = 2ˡ⁻ⁿ * c (mod p)
	//
	// We will call this quantity d. It will be required for reductions later.

	d := (1 << uint(l-n)) * p.C
	dreg := build.GP64()
	build.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

	// Add y into x.
	build.ADDQ(y[0], x[0]) // TODO(mbm): can we replace this with `ADCX`? need to ensure the carry flag is 0
	for i := 1; i < k; i++ {
		build.ADCXQ(y[i], x[i])
	}

	// Both inputs are < 2ˡ so the result is < 2ˡ⁺¹.
	// If the last addition caused a carry into the l'th bit we need to perform a reduction.
	// Prepare the value we will add in to perform the reduction. The addend may be
	// zero or d depending on the carry.
	addend := build.GP64()
	build.MOVQ(zero, addend)
	build.CMOVQCS(dreg, addend)

	// Now add the addend into x.
	build.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
	for i := 1; i < k; i++ {
		build.ADCXQ(zero, x[i])
	}

	// We have added d into the low l bits. Therefore the result is less than 2ˡ + d.
	// But note that it could still be 2ˡ or higher, so we need to perform a
	// second reduction.

	// As before, the addend is either 0 or d depending on the carry from the last add.
	build.MOVQ(zero, addend)
	build.CMOVQCS(dreg, addend)

	// This time we only need to perform one add. The result must be less than 2ˡ + 2*d,
	// therefore provided 2*d does not exceed the size of a limb we can be sure there
	// will be no carry.
	build.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
}

// Mul does a full multiply z = x*y.
func Mul(z, x, y Int) {
	acc := make([]operand.Op, len(z))
	zero := build.GP64()

	for j := 0; j < len(y); j++ {
		build.Commentf("y[%d]", j)
		build.MOVQ(y[j], reg.RDX)
		build.XORQ(zero, zero) // clears flags
		carryinto := [2]int{-1, -1}
		for i := 0; i < len(x); i++ {
			k := i + j
			build.Commentf("x[%d] * y[%d] -> z[%d]", i, j, k)

			// Determine where the results should go.
			var product [2]operand.Op
			var add [2]bool
			for b := 0; b < 2; b++ {
				if acc[k+b] == nil {
					acc[k+b] = build.GP64()
					product[b] = acc[k+b]
				} else {
					product[b] = build.GP64()
					add[b] = true
				}
			}

			// Do the multiply.
			build.MULXQ(x[i], product[0], product[1])

			// Do the adds.
			if add[0] {
				build.ADCXQ(product[0], acc[k])
				carryinto[0] = k + 1
			}
			if add[1] {
				build.ADOXQ(product[1], acc[k+1])
				carryinto[1] = k + 2
			}
		}

		if carryinto[0] > 0 {
			build.ADCXQ(zero, acc[carryinto[0]])
		}
		if carryinto[1] > 0 {
			build.ADOXQ(zero, acc[carryinto[1]])
		}

		//
		build.MOVQ(acc[j], z[j])
	}

	for j := len(y); j < len(z); j++ {
		build.MOVQ(acc[j], z[j])
	}
}

// add builds an addition function.
func add(p Crandall) {
	build.TEXT("Add"+p.Slug(), build.NOSPLIT, "func(x, y *[32]byte)")

	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := NewIntLimb64(4)
	for i := 0; i < 4; i++ {
		build.MOVQ(xb.Offset(8*i), x[i])
	}

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := NewIntLimb64(4)
	for i := 0; i < 4; i++ {
		build.MOVQ(yb.Offset(8*i), y[i])
	}

	Add(x, y, p)

	for i := 0; i < 4; i++ {
		build.MOVQ(x[i], xb.Offset(8*i))
	}

	build.RET()
}

// mul builds a multiplication function.
func mul(p Crandall) {
	build.TEXT("Mul", build.NOSPLIT, "func(z *[64]byte, x, y *[32]byte)")

	zb := operand.Mem{Base: build.Load(build.Param("z"), build.GP64())}
	z := NewIntFromMem(zb, 8)

	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := NewIntFromMem(yb, 4)

	Mul(z, x, y)

	build.RET()

}

func main() {
	p := Crandall{N: 255, C: 19} // Curve25519
	//p := Crandall{N: 383, C: 187} // Curve383187

	add(p)
	mul(p)

	build.Generate()
}
