package fp25519

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"

	"github.com/cloudflare/circl/math/fp25519"
)

func NumTrials() int {
	if testing.Short() {
		return 1 << 8
	}
	return 1 << 15
}

func TestMul(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, y, expect fp25519.Elt
		rand.Read(x[:])
		rand.Read(y[:])

		var xb, yb, got Elt
		copy(xb[:], x[:])
		copy(yb[:], y[:])

		fp25519.Mul(&expect, &x, &y)
		Mul(&got, &xb, &yb)

		if !bytes.Equal(got[:], expect[:]) {
			t.Logf(" trial = %d", trial)
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", expect)
			t.Fail()
		}
	}
}

func TestAdd(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, y, expect fp25519.Elt
		rand.Read(x[:])
		rand.Read(y[:])

		var xb, yb Elt
		copy(xb[:], x[:])
		copy(yb[:], y[:])

		fp25519.Add(&expect, &x, &y)
		Add(&xb, &yb)

		if !bytes.Equal(xb[:], expect[:]) {
			t.Logf(" trial = %d", trial)
			t.Logf("   got = %x", xb)
			t.Logf("expect = %x", expect)
			t.Fail()
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
