// Code generated by ec3. DO NOT EDIT.

package fp25519

// Size of a field element in bytes.
const Size = 32

// Elt is a field element.
type Elt [32]uint8

// Sqr computes z = x^2 (mod p).
func Sqr(z *Elt, x *Elt) {
	Mul(z, x, x)
}

// Inv computes z = 1/x (mod p).
func Inv(z *Elt, x *Elt) {
	var t [3]Elt
	Sqr(z, x)
	Mul(z, x, z)
	Sqr(&t[0], z)
	for s := 1; s < 2; s++ {
		Sqr(&t[0], &t[0])
	}
	Mul(&t[0], z, &t[0])
	Sqr(&t[1], &t[0])
	for s := 1; s < 4; s++ {
		Sqr(&t[1], &t[1])
	}
	Mul(&t[0], &t[0], &t[1])
	for s := 0; s < 2; s++ {
		Sqr(&t[0], &t[0])
	}
	Mul(&t[0], z, &t[0])
	Sqr(&t[1], &t[0])
	for s := 1; s < 10; s++ {
		Sqr(&t[1], &t[1])
	}
	Mul(&t[1], &t[0], &t[1])
	for s := 0; s < 10; s++ {
		Sqr(&t[1], &t[1])
	}
	Mul(&t[1], &t[0], &t[1])
	Sqr(&t[2], &t[1])
	for s := 1; s < 30; s++ {
		Sqr(&t[2], &t[2])
	}
	Mul(&t[1], &t[1], &t[2])
	Sqr(&t[2], &t[1])
	for s := 1; s < 60; s++ {
		Sqr(&t[2], &t[2])
	}
	Mul(&t[1], &t[1], &t[2])
	Sqr(&t[2], &t[1])
	for s := 1; s < 120; s++ {
		Sqr(&t[2], &t[2])
	}
	Mul(&t[1], &t[1], &t[2])
	for s := 0; s < 10; s++ {
		Sqr(&t[1], &t[1])
	}
	Mul(&t[0], &t[0], &t[1])
	for s := 0; s < 2; s++ {
		Sqr(&t[0], &t[0])
	}
	Mul(&t[0], x, &t[0])
	for s := 0; s < 3; s++ {
		Sqr(&t[0], &t[0])
	}
	Mul(z, z, &t[0])
}
