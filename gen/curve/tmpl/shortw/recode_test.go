package shortw

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
	"testing"
)

func RandScalar(t *testing.T) *big.Int {
	t.Helper()
	N := curvename.Params().N
	for {
		k, err := rand.Int(rand.Reader, N)
		if err != nil {
			t.Fatal(err)
		}
		if k.Sign() == 0 {
			continue
		}
		return k
	}
}

func RandOddScalar(t *testing.T) *big.Int {
	t.Helper()
	k := RandScalar(t)
	N := curvename.Params().N
	if k.Bit(0) == 0 {
		k.Neg(k).Mod(k, N)
	}
	return k
}

func TestScalarIntRoundtrip(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		var k scalar
		expect := RandScalar(t)
		k.SetInt(expect)
		got := k.Int()
		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}

func TestScalarFixedWindowRecode(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandOddScalar(t)

		var K scalar
		K.SetInt(k)
		digits := K.FixedWindowRecode()

		// Verify all digits are odd.
		for i, digit := range digits {
			if (digit & 1) != 1 {
				t.Fatalf("digit %d is not odd", i)
			}
		}

		// Confirm the sum is correct.
		x := new(big.Int)
		for i := len(digits) - 1; i >= 0; i-- {
			x.Lsh(x, ConstW-1)
			x.Add(x, big.NewInt(int64(digits[i])))
		}

		if k.Cmp(x) != 0 {
			t.Logf("     k = %x", k)
			t.Logf("digits = %d", digits)
			t.Logf("   got = %x", x)
			t.FailNow()
		}
	}
}

func TestScalarSubInt(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x := RandScalar(t)
		v := mathrand.Int31n(64) - 32

		// Compute subtraction via scalar type.
		var k scalar
		k.SetInt(x)
		k.SubInt32(v)
		got := k.Int()

		// Compute expectation.
		expect := new(big.Int).Sub(x, new(big.Int).SetInt64(int64(v)))

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}

func TestScalarRsh(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x := RandScalar(t)
		s := uint(1 + mathrand.Intn(63))

		// Compute shift via scalar type.
		var k scalar
		k.SetInt(x)
		k.Rsh(s)
		got := k.Int()

		// Compute expectation.
		expect := new(big.Int).Rsh(x, s)

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}
