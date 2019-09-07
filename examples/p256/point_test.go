package p256

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

var (
	cur = P256()
	ref = elliptic.P256()
)

func RandPoint(t *testing.T) (x, y *big.Int) {
	t.Helper()
	k, err := rand.Int(rand.Reader, ref.Params().N)
	if err != nil {
		t.Fatal(err)
	}
	return ref.ScalarBaseMult(k.Bytes())
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
		j := a.Jacobian()
		a2 := j.Affine()
		gx, gy := a2.Coordinates()
		EqualInt(t, "x", x, gx)
		EqualInt(t, "y", y, gy)
	}
}

func TestAffineJacobianProjectiveRoundTrip(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x, y := RandPoint(t)
		a := NewAffine(x, y)
		j := a.Jacobian()
		p := j.Projective()
		a2 := p.Affine()
		gx, gy := a2.Coordinates()
		EqualInt(t, "x", x, gx)
		EqualInt(t, "y", y, gy)
	}
}

func TestAddPoints(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		ex, ey := ref.Add(x1, y1, x2, y2)
		gx, gy := cur.Add(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestDoublePoint(t *testing.T) {
	for trial := 0; trial < NumTrials(); trial++ {
		x1, y1 := RandPoint(t)

		ex, ey := cur.Double(x1, y1)
		gx, gy := ref.Double(x1, y1)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestScalarMult(t *testing.T) {
	t.Logf(" b = %x", *b)

	b2 := new(Elt).SetInt(ref.Params().B)
	Encode(b2, b2)
	t.Logf("b2 = %x", *b2)

	for trial := 0; trial < 128; trial++ {
		k := RandScalar(t)
		x1, y1 := RandPoint(t)

		gx, gy := cur.ScalarMult(x1, y1, k.Bytes())
		ex, ey := ref.ScalarMult(x1, y1, k.Bytes())

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
