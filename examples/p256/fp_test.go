package p256

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func NumTrials() int {
	if testing.Short() {
		return 1 << 8
	}
	return 1 << 15
}

var p *big.Int

func init() {
	p, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
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

func TestAdd(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		var x, y Elt
		rand.Read(x[:])
		rand.Read(y[:])

		got := x
		Add(&got, &y)

		expect := ExpectAdd(x, y)

		if got != expect {
			t.Logf("     x = %x", x)
			t.Logf("     y = %x", y)
			t.Logf("   got = %x", got)
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
