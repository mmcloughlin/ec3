package p256

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

var curve = elliptic.P256()

func RandPoint(t *testing.T) (x, y *big.Int) {
	t.Helper()
	params := curve.Params()
	k, err := rand.Int(rand.Reader, params.N)
	if err != nil {
		t.Fatal(err)
	}
	return curve.ScalarBaseMult(k.Bytes())
}

func TestAffineRoundTrip(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x, y := RandPoint(t)
		a := NewAffine(x, y)
		gx, gy := a.Coordinates()
		EqualInt(t, "x", x, gx)
		EqualInt(t, "y", y, gy)
	}
}

func TestAffineJacobianRoundTrip(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x, y := RandPoint(t)
		a := NewAffine(x, y)
		j := NewFromAffine(a)
		a2 := j.Affine()
		gx, gy := a2.Coordinates()
		EqualInt(t, "x", x, gx)
		EqualInt(t, "y", y, gy)
	}
}

func AddPoints(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	a2 := NewAffine(x2, y2)
	j1 := NewFromAffine(a1)
	j2 := NewFromAffine(a2)
	s := new(Jacobian)
	s.Add(j1, j2)
	return s.Affine().Coordinates()
}

func DoublePoint(x1, y1 *big.Int) (x, y *big.Int) {
	a1 := NewAffine(x1, y1)
	j1 := NewFromAffine(a1)
	d := new(Jacobian)
	d.Double(j1)
	return d.Affine().Coordinates()
}

func TestAddPoints(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		ex, ey := curve.Add(x1, y1, x2, y2)
		gx, gy := AddPoints(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestDoublePoint(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x1, y1 := RandPoint(t)

		ex, ey := curve.Double(x1, y1)
		gx, gy := DoublePoint(x1, y1)

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
