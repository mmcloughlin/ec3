package verif

import "github.com/mmcloughlin/ec3/z3"

// Limbs breaks x into s-bit limbs.
func Limbs(x *z3.BV, s uint) []*z3.BV {
	limbs := []*z3.BV{}
	for l := uint(0); l < x.Bits(); l += s {
		limb := x.Extract(l+s-1, l)
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
