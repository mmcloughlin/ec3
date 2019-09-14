package p256

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

func ScalarR() *big.Int {
	r := bigint.Pow2(256)
	r.Mod(r, scalarp)
	return r
}

func ScalarRInv() *big.Int {
	rinv := ScalarR()
	rinv.ModInverse(rinv, scalarp)
	return rinv
}

func RandScalarElt() scalar {
	var r scalar
	rand.Read(r[:])

	x := IntFromBytesLittleEndian(r[:])
	x.Mul(x, ScalarR())
	x.Mod(x, scalarp)

	BytesFromIntLittleEndian(r[:], x)
	return r
}

func ExpectScalarMul(x, y scalar) scalar {
	xi := IntFromBytesLittleEndian(x[:])
	yi := IntFromBytesLittleEndian(y[:])

	zi := new(big.Int).Mul(xi, yi)
	zi.Mul(zi, ScalarRInv())
	zi.Mod(zi, scalarp)

	var z scalar
	BytesFromIntLittleEndian(z[:], zi)
	return z
}

func TestScalarMul(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x := RandScalarElt()
		y := RandScalarElt()

		var got scalar
		scalarmul(&got, &x, &y)

		expect := ExpectScalarMul(x, y)

		if got != expect {
			t.Logf(" trial = %d", trial)
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestScalarMulEdgeCases(t *testing.T) {
	// Values that might exercise edge cases.
	interesting := map[string]*big.Int{
		"0":     big.NewInt(0),
		"1":     big.NewInt(1),
		"R":     ScalarR(),
		"(1/R)": ScalarRInv(),
		"(p-1)": new(big.Int).Sub(scalarp, big.NewInt(1)),
	}

	// Test cases for every pair of interesting numbers.
	type testcase struct {
		Name string
		X, Y *big.Int
	}
	cases := []testcase{}
	for nx, x := range interesting {
		for ny, y := range interesting {
			cases = append(cases, testcase{
				Name: fmt.Sprintf("%s*%s", nx, ny),
				X:    new(big.Int).Set(x),
				Y:    new(big.Int).Set(y),
			})
		}
	}

	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			// Montgomery encode.
			c.X.Mul(c.X, ScalarR())
			c.X.Mod(c.X, scalarp)

			c.Y.Mul(c.Y, ScalarR())
			c.Y.Mod(c.Y, scalarp)

			// Compute expectation.
			expect := new(big.Int).Mul(c.X, c.Y)
			expect.Mul(expect, ScalarRInv())
			expect.Mod(expect, scalarp)

			// What do we get?
			var x, y, z scalar
			BytesFromIntLittleEndian(x[:], c.X)
			BytesFromIntLittleEndian(y[:], c.Y)
			scalarmul(&z, &x, &y)
			got := IntFromBytesLittleEndian(z[:])

			if expect.Cmp(got) != 0 {
				t.Logf("     x = %x", c.X)
				t.Logf("     y = %x", c.Y)
				t.Logf("   got = %x", got)
				t.Logf("expect = %x", expect)
				t.Fail()
			}
		})
	}
}
