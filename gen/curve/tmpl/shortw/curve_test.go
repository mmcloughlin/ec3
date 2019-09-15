// CodeGenerationWarning

package shortw

import (
	"math/big"
	"testing"
)

var (
	cur = CURVENAME()
	ref = curvename.CurveParams
)

func TestCurveAddRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x1, y1 := RandPoint(t)
		x2, y2 := RandPoint(t)

		gx, gy := cur.Add(x1, y1, x2, y2)
		ex, ey := ref.Add(x1, y1, x2, y2)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveAddAsDouble(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		gx, gy := cur.Add(x, y, x, y)
		ex, ey := ref.Double(x, y)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveAddNegative(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		nx := new(big.Int).Set(x)
		ny := new(big.Int).Neg(y)

		gx, gy := cur.Add(x, y, nx, ny)
		zero := new(big.Int)

		EqualInt(t, "x", zero, gx)
		EqualInt(t, "y", zero, gy)
	}
}

func TestCurveDoubleRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x, y := RandPoint(t)

		gx, gy := ref.Double(x, y)
		ex, ey := cur.Double(x, y)

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarMultRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)
		x, y := RandPoint(t)

		gx, gy := cur.ScalarMult(x, y, k.Bytes())
		ex, ey := ref.ScalarMult(x, y, k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveScalarBaseMultRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		gx, gy := cur.ScalarBaseMult(k.Bytes())
		ex, ey := ref.ScalarBaseMult(k.Bytes())

		EqualInt(t, "x", ex, gx)
		EqualInt(t, "y", ey, gy)
	}
}

func TestCurveInverseRand(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandScalarNonZero(t)

		got := cur.Inverse(k)

		expect := new(big.Int).Set(k)
		expect.ModInverse(expect, ref.N)

		EqualInt(t, "inv", expect, got)
	}
}
