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

func ExpectEncode(x Elt) Elt {
	xi := IntFromBytesLittleEndian(x[:])

	zi := new(big.Int).Mul(xi, R())
	zi.Mod(zi, p)

	var z Elt
	BytesFromIntLittleEndian(z[:], zi)
	return z
}

func ExpectDecode(x Elt) Elt {
	xi := IntFromBytesLittleEndian(x[:])

	zi := new(big.Int).Mul(xi, RInv())
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

		px, py := x, y
		Add(&got, &x, &y)
		expect := ExpectAdd(x, y)

		if px != x || py != y {
			t.Fatal("changed inputs")
		}

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

		px, py := x, y
		Sub(&got, &x, &y)
		expect := ExpectSub(x, y)

		if px != x || py != y {
			t.Fatal("changed inputs")
		}

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
		px, py := x, y

		var got Elt
		Mul(&got, &x, &y)

		expect := ExpectMul(x, y)

		if px != x || py != y {
			t.Fatal("changed inputs")
		}

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestSqr(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x := RandElt()
		px := x

		var got Elt
		Sqr(&got, &x)

		expect := ExpectMul(x, x)

		if px != x {
			t.Fatal("changed input")
		}

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func DebugMontMul(t *testing.T, x, y *big.Int) {
	m := new(big.Int).Mul(x, y)
	acc := m
	for i := uint(0); i < 256; i += 64 {
		t.Logf("step %3d: acc = %0129x", i, acc)

		// Step 2.1: u_i = x_i * m' (mod b)
		u := new(big.Int).Rsh(acc, i)
		u.And(u, bigint.Ones(64))

		// Step 2.2: x += u_i * m * b^i
		u.Mul(u, p)
		u.Lsh(u, i)
		acc.Add(acc, u)
	}

	t.Logf("  reduced acc = %0129x", acc)

	acc.Rsh(acc, 256)
	t.Logf("    shift acc = %0129x", acc)

	t.Logf("        cmp p = %d", acc.Cmp(p))
	if acc.Cmp(p) >= 0 {
		acc.Sub(acc, p)
	}

	t.Logf("    final acc = %0129x", acc)
}

func TestMulEdgeCases(t *testing.T) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	p1 := new(big.Int).Sub(p, one)

	interesting := map[string]*big.Int{
		"0":     zero,
		"1":     one,
		"R":     R(),
		"(1/R)": RInv(),
		"(p-1)": p1,
	}
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

			t.Logf("     x = %x (encoded)", c.X)
			t.Logf("     y = %x (encoded)", c.Y)

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

			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)

			if expect.Cmp(got) != 0 {
				DebugMontMul(t, c.X, c.Y)
				t.Fail()
			}
		})
	}
}

func TestInv(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x Elt
		rand.Read(x[:])

		// Compute via.
		var m, minv, got Elt
		Encode(&m, &x)
		Inv(&minv, &m)
		Decode(&got, &minv)

		// Expect.
		xi := IntFromBytesLittleEndian(x[:])
		xi.ModInverse(xi, p)
		var expect Elt
		BytesFromIntLittleEndian(expect[:], xi)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestInvProperty(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var xm, xminv, p, got Elt
		x := RandElt()

		Encode(&xm, &x)
		Inv(&xminv, &xm)
		Mul(&p, &xm, &xminv)
		Decode(&got, &p)

		one := Elt{1}
		if one != got {
			t.Logf("     x = %x", x)
			t.Logf("    xm = %x", xm)
			t.Logf(" xminv = %x", xminv)
			t.Logf("     p = %x", p)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", one)
			t.Fatal("invalid inverse")
		}
	}
}

func TestSetInt64(t *testing.T) {
	got := new(Elt).SetInt64(1)
	expect := &Elt{1}
	if *got != *expect {
		t.Fatal("SetInt64(1) failed")
	}
}

func TestEncode(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, got Elt
		rand.Read(x[:])

		Encode(&got, &x)
		expect := ExpectEncode(x)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestDecode(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, got Elt
		rand.Read(x[:])

		Decode(&got, &x)
		expect := ExpectDecode(x)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.FailNow()
		}
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, m, got Elt
		rand.Read(x[:])
		Encode(&m, &x)
		Decode(&got, &m)
		if got != x {
			t.Fatal("failed roundtrip")
		}
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
