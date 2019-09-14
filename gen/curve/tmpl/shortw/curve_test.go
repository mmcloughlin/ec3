package shortw

import (
	"math/big"
	"testing"
)

var (
	cur = CURVENAME()
	ref = curvename.CurveParams
)

func RandPoint(t *testing.T) (x, y *big.Int) {
	t.Helper()
	k := RandScalarNonZero(t)
	return ref.ScalarBaseMult(k.Bytes())
}

func TestAdd(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		ex, ey := ref.Add(x1, y1, x2, y2)
		gx, gy := cur.Add(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestDouble(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)

		ex, ey := cur.Double(x1, y1)
		gx, gy := ref.Double(x1, y1)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestScalarMult(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)
		x1, y1 := RandPoint(t)

		gx, gy := cur.ScalarMult(x1, y1, k.Bytes())
		ex, ey := ref.ScalarMult(x1, y1, k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestInverse(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		got := cur.Inverse(k)

		expect := new(big.Int).Set(k)
		expect.ModInverse(expect, ref.N)

		EqualInt(t, "inv", expect, got)
	}
}

func EqualInt(t *testing.T, name string, expect, got *big.Int) {
	if got.Cmp(expect) != 0 {
		t.Logf("   got %x", got)
		t.Logf("expect %x", expect)
		t.Fatalf("%s: not equal", name)
	}
}
