package ec3

import (
	"crypto/rand"
	"testing"

	"github.com/cloudflare/circl/math/fp25519"
)

//go:generate go run asm.go -out fp25519.s -stubs stub.go

func TestAdd(t *testing.T) {
	trials := 1 << 8
	for trial := 0; trial < trials; trial++ {
		var x, y, expect fp25519.Elt
		rand.Read(x[:])
		rand.Read(y[:])

		var xb, yb [32]byte
		copy(xb[:], x[:])
		copy(yb[:], y[:])

		fp25519.Add(&expect, &x, &y)
		Add25519(&xb, &yb)

		if xb != expect {
			t.Logf(" trial = %d", trial)
			t.Logf("   got = %x", xb)
			t.Logf("expect = %x", expect)
			t.Fail()
		}
	}
}
