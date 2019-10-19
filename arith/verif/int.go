package verif

import "github.com/mmcloughlin/ec3/internal/z3"

// Limbs breaks x into k-bit limbs.
func Limbs(x *z3.BV, k uint) []*z3.BV {
	limbs := []*z3.BV{}
	for l := uint(0); l < x.Bits(); l += k {
		limb := x.Extract(l+k-1, l)
		limbs = append(limbs, limb)
	}
	return limbs
}

// FromLimbs builds a full integer from component limbs.
func FromLimbs(limbs []*z3.BV) *z3.BV {
	x := limbs[0]
	for i := 1; i < len(limbs); i++ {
		x = limbs[i].Concat(x)
	}
	return x
}
