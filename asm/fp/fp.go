package fp

import (
	"github.com/mmcloughlin/avo/build"

	"github.com/mmcloughlin/ec3/asm/mp"
)

type Properties interface {
	// ElementBits returns the number of bits used to represent a field element.
	// This will be larger than the size of the prime if it's not on a word boundary.
	ElementBits() int

	// ElementSize returns the number of bytes used to represent a field element.
	ElementSize() int

	// Limbs returns the number of 64-bit limbs required for a field element.
	Limbs() int
}

type Field interface {
	Properties

	Build(*build.Context) Builder
}

type Builder interface {
	Properties

	// Add generates code to add y into x modulo p.
	//	x ≡ x + y (mod p)
	Add(x, y mp.Int)

	// ReduceDouble computes z congruent to x modulo p. Let the element size be 2ˡ.
	// This function assumes x < 2²ˡ and produces z < 2ˡ. Note that z is not
	// guaranteed to be less than p.
	ReduceDouble(z, x mp.Int)
}
