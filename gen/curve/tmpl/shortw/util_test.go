// CodeGenerationWarning

package shortw

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func RandScalarNonZero(tb testing.TB) *big.Int {
	tb.Helper()
	N := curvename.Params().N
	for {
		k, err := rand.Int(rand.Reader, N)
		if err != nil {
			tb.Fatal(err)
		}
		if k.Sign() == 0 {
			continue
		}
		return k
	}
}

func RandOddScalar(tb testing.TB) *big.Int {
	tb.Helper()
	k := RandScalarNonZero(tb)
	N := curvename.Params().N
	if k.Bit(0) == 0 {
		k.Neg(k).Mod(k, N)
	}
	return k
}

func RandPoint(tb testing.TB) (x, y *big.Int) {
	tb.Helper()
	k := RandScalarNonZero(tb)
	return curvename.Params().ScalarBaseMult(k.Bytes())
}

func EqualInt(t *testing.T, name string, expect, got *big.Int) {
	t.Helper()
	if got.Cmp(expect) != 0 {
		t.Logf("   got %x", got)
		t.Logf("expect %x", expect)
		t.Fatalf("%s: not equal", name)
	}
}
