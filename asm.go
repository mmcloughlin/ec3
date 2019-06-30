// +build ignore

package main

import (
	"strconv"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"

	"github.com/mmcloughlin/ec3/prime"
)

// Slug returns a concise name for p.
func Slug(p prime.Crandall) string {
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

// Zero64 returns a 64-bit register initialized to zero.
func Zero64() reg.Register {
	zero := build.GP64()
	build.XORQ(zero, zero)
	return zero
}

// AddModP adds y into x modulo p.
//	x ≡ x + y (mod p)
func AddModP(x, y Int, p prime.Crandall) {
	n := p.Bits()
	l := NextMultiple(n, 64)
	k := l / 64

	// Prepare a zero register.
	zero := Zero64()

	// Note that for the Crandall prime we have
	//
	//	2ⁿ - c ≡ 0 (mod p)
	//	2ⁿ ≡ c (mod p)
	//
	// However n may not be on a limb boundary, so we actually need the identity
	//
	//	2ˡ ≡ 2ˡ⁻ⁿ * c (mod p)
	//
	// We will call this quantity d. It will be required for reductions later.

	// TODO(mbm): refactor d computation
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
	// TODO(mbm): assert d is within an acceptable range
	build.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
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

// ReduceDouble computes z congruent to x modulo p. Let the element size be 2ˡ.
// This function assumes x < 2²ˡ and produces z < 2ˡ. Note that z is not
// guaranteed to be less than p.
func ReduceDouble(z, x Int, p prime.Crandall) {
	// TODO(mbm): helpers for limb size computations
	n := p.Bits()
	l := NextMultiple(n, 64)
	k := l / 64

	// Prepare a zero register.
	zero := Zero64()

	// Compute the reduction additive d.
	d := (1 << uint(l-n)) * p.C
	dreg := build.GP64()
	build.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

	// Stage 1: upper bound 2²ˡ → 2ˡ + d*2ˡ.
	//
	// At this point x is a 2*l bit quantity which we'll view as two l-bit halfs.
	//
	//	x = H ∥ L = 2ˡ * H + L ≡ d*H + L (mod p)
	//
	// Therefore the first reduction stage multiplies the top limbs by d and
	// accumulates the result into the low limbs. Note the result will have an
	// additional limb.

	// Multiply r = d*H.
	r := NewIntLimb64(k + 1)
	build.MOVQ(dreg, reg.RDX)
	build.XORQ(r[0], r[0]) // also clears flags
	for i := 0; i < k; i++ {
		lo := build.GP64()
		build.MULXQ(x[i+k], lo, r[i+1])
		build.ADCXQ(lo, r[i])
	}

	// Add r += x.
	for i := 0; i < k; i++ {
		build.ADOXQ(x[i], r[i])
	}
	build.ADOXQ(zero, r[k])

	// Stage 2: (d+1)*2ˡ → 2ˡ + (d+1)*d
	//
	// Currently r has one too many limbs, so we need to reduce again. The value in
	// the top limb is ⩽ d. When we reduce we have to multiply by d again, so the
	// result cannot exceed d^2. Provided d is small, the result will not exceed a single limb.

	// TODO(mbm): assert d is within an acceptable range

	top := r[k]
	build.IMULQ(dreg, top) // clears flags
	build.ADCXQ(top, r[0])
	for i := 1; i < k; i++ {
		build.ADCXQ(zero, r[i])
	}

	// Stage 3: finish
	//
	// It is still possible that the final add carried, in which case we need one final
	// add to complete the reduction.
	addend := build.GP64()
	build.MOVQ(zero, addend)
	build.CMOVQCS(dreg, addend)
	build.ADDQ(addend, r[0])

	// Write out the result.
	for i := 0; i < k; i++ {
		build.MOVQ(r[i], z[i])
	}
}

// addmod builds a modular addition function.
func addmod(p prime.Crandall) {
	build.TEXT("Add"+Slug(p), build.NOSPLIT, "func(x, y *[32]byte)")

	// TODO(mbm): helper for loading integer from memory
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

	AddModP(x, y, p)

	for i := 0; i < 4; i++ {
		build.MOVQ(x[i], xb.Offset(8*i))
	}

	build.RET()
}

// mul builds a multiplication function.
func mul() {
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

// mulmod builds a modular multiplication function.
func mulmod(p prime.Crandall) {
	build.TEXT("Mul"+Slug(p), build.NOSPLIT, "func(z *[32]byte, x, y *[32]byte)")

	// Load arguments.
	zb := operand.Mem{Base: build.Load(build.Param("z"), build.GP64())}
	z := NewIntFromMem(zb, 4)

	xb := operand.Mem{Base: build.Load(build.Param("x"), build.GP64())}
	x := NewIntFromMem(xb, 4)

	yb := operand.Mem{Base: build.Load(build.Param("y"), build.GP64())}
	y := NewIntFromMem(yb, 4)

	// Perform multiplication.
	// TODO(mbm): is it possible to store the intermediate result in registers?
	mb := build.AllocLocal(8 * 8)
	m := NewIntFromMem(mb, 8)

	Mul(m, x, y)

	// Reduce.
	build.Comment("Reduction.")
	ReduceDouble(z, m, p)

	build.RET()
}

func main() {
	// Multi-precision.
	mul()

	// Fp25519
	p := prime.Crandall{N: 255, C: 19}
	addmod(p)
	mulmod(p)

	build.Generate()
}
