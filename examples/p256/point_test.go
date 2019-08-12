package p256

import (
	"crypto/elliptic"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

var curve = elliptic.P256()

func RandPoint() (x, y *big.Int) {
	var k [32]byte
	rand.Read(k[:])
	return curve.ScalarBaseMult(k[:])
}

func TestAffineRoundTrip(t *testing.T) {
	x, y := RandPoint()
	a := NewAffine(x, y)
	gx, gy := a.Coordinates()
	EqualInt(t, "x", x, gx)
	EqualInt(t, "y", y, gy)
}

func TestAffineJacoboanRoundTrip(t *testing.T) {
	x, y := RandPoint()
	a := NewAffine(x, y)
	j := NewFromAffine(a)
	a2 := j.Affine()
	gx, gy := a2.Coordinates()
	EqualInt(t, "x", x, gx)
	EqualInt(t, "y", y, gy)
}

func AddPoints(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	a1 := NewAffine(x1, y1)
	a2 := NewAffine(x2, y2)
	j1 := NewFromAffine(a1)
	j2 := NewFromAffine(a2)
	s := new(Jacobian)
	s.Add(j1, j2)
	return s.Affine().Coordinates()
}

func TestAddPoints(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x1, y1 := RandPoint()
		x2, y2 := RandPoint()

		ex, ey := curve.Add(x1, y1, x2, y2)
		gx, gy := AddPoints(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func EqualInt(t *testing.T, name string, expect, got *big.Int) {
	if !bigint.Equal(expect, got) {
		t.Logf("   got %x", got)
		t.Logf("expect %x", expect)
		t.Fatalf("%s: not equal", name)
	}
}
