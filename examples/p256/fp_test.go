package p256

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

func NumTrials() int {
	if testing.Short() {
		return 1 << 8
	}
	return 1 << 20
}

var p *big.Int

func init() {
	p, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
}

func R() *big.Int {
	r := bigint.Pow2(256)
	r.Mod(r, p)
	return r
}

func RandElt() Elt {
	var r Elt
	rand.Read(r[:])

	x := IntFromBytesLittleEndian(r[:])
	x.Mul(x, R())
	x.Mod(x, p)

	BytesFromIntLittleEndian(r[:], x)
	return r
}

func RInv() *big.Int {
	rinv := R()
	rinv.ModInverse(rinv, p)
	return rinv
}

func ExpectAdd(x, y Elt) Elt {
	xi := IntFromBytesLittleEndian(x[:])
	yi := IntFromBytesLittleEndian(y[:])
	zi := new(big.Int).Add(xi, yi)
	zi.Mod(zi, p)
	var z Elt
	BytesFromIntLittleEndian(z[:], zi)
	return z
}

func ExpectSub(x, y Elt) Elt {
	xi := IntFromBytesLittleEndian(x[:])
	yi := IntFromBytesLittleEndian(y[:])
	zi := new(big.Int).Sub(xi, yi)
	zi.Mod(zi, p)
	var z Elt
	BytesFromIntLittleEndian(z[:], zi)
	return z
}

func ExpectMul(x, y Elt) Elt {
	xi := IntFromBytesLittleEndian(x[:])
	yi := IntFromBytesLittleEndian(y[:])

	zi := new(big.Int).Mul(xi, yi)
	zi.Mul(zi, RInv())
	zi.Mod(zi, p)

	var z Elt
	BytesFromIntLittleEndian(z[:], zi)
	return z
}

func TestAdd(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, y, got Elt
		rand.Read(x[:])
		rand.Read(y[:])

		Add(&got, &x, &y)
		expect := ExpectAdd(x, y)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestSub(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, y, got Elt
		rand.Read(x[:])
		rand.Read(y[:])

		Sub(&got, &x, &y)
		expect := ExpectSub(x, y)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestMul(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x := RandElt()
		y := RandElt()

		var got Elt
		Mul(&got, &x, &y)

		expect := ExpectMul(x, y)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestMulEdgeCases(t *testing.T) {
	// Values that might exercise edge cases.
	interesting := map[string]*big.Int{
		"0":     big.NewInt(0),
		"1":     big.NewInt(1),
		"R":     R(),
		"(1/R)": RInv(),
		"(p-1)": new(big.Int).Sub(p, big.NewInt(1)),
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
			t.Logf("     x = %x", c.X)
			t.Logf("     y = %x", c.Y)

			// Montgomery encode.
			c.X.Mul(c.X, R())
			c.X.Mod(c.X, p)

			c.Y.Mul(c.Y, R())
			c.Y.Mod(c.Y, p)

			// Compute expectation.
			expect := new(big.Int).Mul(c.X, c.Y)
			expect.Mul(expect, RInv())
			expect.Mod(expect, p)

			// What do we get?
			var x, y, z Elt
			BytesFromIntLittleEndian(x[:], c.X)
			BytesFromIntLittleEndian(y[:], c.Y)
			Mul(&z, &x, &y)
			got := IntFromBytesLittleEndian(z[:])

			if expect.Cmp(got) != 0 {
				t.Logf("   got = %x", got)
				t.Logf("expect = %x", expect)
				t.Fail()
			}
		})
	}
}

func IntFromBytesLittleEndian(b []byte) *big.Int {
	bigendian := append([]byte{}, b...)
	ReverseBytes(bigendian)
	return new(big.Int).SetBytes(bigendian)
}

func BytesFromIntLittleEndian(b []byte, x *big.Int) []byte {
	xb := x.Bytes()
	ReverseBytes(xb)
	copy(b, xb)
	return b
}

func ReverseBytes(b []byte) {
	for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
		b[l], b[r] = b[r], b[l]
	}
}
