package mp

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

// Int represents a multi-precision integer.
type Int []operand.Op

// NewInt builds an empty integer with k limbs.
func NewInt(k int) Int {
	return make(Int, k)
}

// NewIntLimb64 builds multi-precision integer with k 64-bit limbs.
func NewIntLimb64(k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = build.GP64()
	}
	return x
}

// NewIntFromMem builds a multi-precision integer referencing the k 64-bit limbs
// at memory address m.
func NewIntFromMem(m operand.Mem, k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = m.Offset(8 * i)
	}
	return x
}

// Mul does a full multiply z = x*y.
func Mul(z, x, y Int) {
	// TODO(mbm): multi-precision multiply is ugly

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
