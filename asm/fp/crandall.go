package fp

import (
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"

	"github.com/mmcloughlin/ec3/asm"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

// Crandall generates arithmetic modulo a Crandall prime.
type Crandall struct {
	P prime.Crandall
}

// ElementBits returns the number of bits used to represent a field element.
// This will be larger than the size of the prime if it's not on a word boundary.
func (f Crandall) ElementBits() int {
	n := f.P.Bits()
	return ints.NextMultiple(n, 64)
}

// ElementSize returns the number of bytes used to represent a field element.
func (f Crandall) ElementSize() int {
	return f.ElementBits() / 8
}

// Limbs returns the number of 64-bit limbs required for a field element.
func (f Crandall) Limbs() int {
	return f.ElementBits() / 64
}

// ReductionMultiplier returns the value an element must be multiplied by to reduce it modulo p upon overflow of the element size.
// Note that for the Crandall prime we have
//
//	2ⁿ - c ≡ 0 (mod p)
//	2ⁿ ≡ c (mod p)
//
// However n may not be on a limb boundary, so we actually need the identity
//
//	2ˡ ≡ 2ˡ⁻ⁿ * c (mod p)
//
// We'll call this the reduction multiplier.
func (f Crandall) ReductionMultiplier() uint32 {
	n := f.P.Bits()
	l := f.ElementBits()
	// TODO(mbm): check for overflow
	return uint32((1 << uint(l-n)) * f.P.C)
}

// Add adds y into x modulo p.
//	x ≡ x + y (mod p)
func (f Crandall) Add(ctx *build.Context, x, y mp.Int) {
	k := f.Limbs()

	// Prepare a zero register.
	zero := asm.Zero64(ctx)

	// Load reduction multiplier.
	d := f.ReductionMultiplier()
	dreg := ctx.GP64()
	ctx.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

	// Add y into x.
	ctx.ADDQ(y[0], x[0]) // TODO(mbm): can we replace this with `ADCX`? need to ensure the carry flag is 0
	for i := 1; i < k; i++ {
		ctx.ADCXQ(y[i], x[i])
	}

	// Both inputs are < 2ˡ so the result is < 2ˡ⁺¹.
	// If the last addition caused a carry into the l'th bit we need to perform a reduction.
	// Prepare the value we will add in to perform the reduction. The addend may be
	// zero or d depending on the carry.
	addend := ctx.GP64()
	ctx.MOVQ(zero, addend)
	ctx.CMOVQCS(dreg, addend)

	// Now add the addend into x.
	ctx.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
	for i := 1; i < k; i++ {
		ctx.ADCXQ(zero, x[i])
	}

	// We have added d into the low l bits. Therefore the result is less than 2ˡ + d.
	// But note that it could still be 2ˡ or higher, so we need to perform a
	// second reduction.

	// As before, the addend is either 0 or d depending on the carry from the last add.
	ctx.MOVQ(zero, addend)
	ctx.CMOVQCS(dreg, addend)

	// This time we only need to perform one add. The result must be less than 2ˡ + 2*d,
	// therefore provided 2*d does not exceed the size of a limb we can be sure there
	// will be no carry.
	// TODO(mbm): assert d is within an acceptable range
	ctx.ADDQ(addend, x[0]) // TODO(mbm): replace with ADCX?
}

// ReduceDouble computes z congruent to x modulo p. Let the element size be 2ˡ.
// This function assumes x < 2²ˡ and produces z < 2ˡ. Note that z is not
// guaranteed to be less than p.
func (f Crandall) ReduceDouble(ctx *build.Context, z, x mp.Int) {
	k := f.Limbs()

	// Prepare a zero register.
	zero := asm.Zero64(ctx)

	// Compute the reduction multiplier.
	d := f.ReductionMultiplier()
	dreg := ctx.GP64()
	ctx.MOVQ(operand.U32(d), dreg) // TODO(mbm): is U32 right?

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
	r := mp.NewIntLimb64(ctx, k+1)
	ctx.MOVQ(dreg, reg.RDX)
	ctx.XORQ(r[0], r[0]) // also clears flags
	for i := 0; i < k; i++ {
		lo := ctx.GP64()
		ctx.MULXQ(x[i+k], lo, r[i+1])
		ctx.ADCXQ(lo, r[i])
	}

	// Add r += x.
	for i := 0; i < k; i++ {
		ctx.ADOXQ(x[i], r[i])
	}
	ctx.ADOXQ(zero, r[k])

	// Stage 2: (d+1)*2ˡ → 2ˡ + (d+1)*d
	//
	// Currently r has one too many limbs, so we need to reduce again. The value in
	// the top limb is ⩽ d. When we reduce we have to multiply by d again, so the
	// result cannot exceed d^2. Provided d is small, the result will not exceed a single limb.

	// TODO(mbm): assert d is within an acceptable range

	top := r[k]
	ctx.IMULQ(dreg, top) // clears flags
	ctx.ADCXQ(top, r[0])
	for i := 1; i < k; i++ {
		ctx.ADCXQ(zero, r[i])
	}

	// Stage 3: finish
	//
	// It is still possible that the final add carried, in which case we need one final
	// add to complete the reduction.
	addend := ctx.GP64()
	ctx.MOVQ(zero, addend)
	ctx.CMOVQCS(dreg, addend)
	ctx.ADDQ(addend, r[0])

	// Write out the result.
	for i := 0; i < k; i++ {
		ctx.MOVQ(r[i], z[i])
	}
}
