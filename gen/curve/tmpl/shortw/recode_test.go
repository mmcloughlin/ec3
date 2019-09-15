// CodeGenerationWarning

package shortw

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestScalarFixedWindowRecode(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		k := RandOddScalar(t)

		var K scalar
		K.SetIntRaw(k)
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
		x := RandScalarNonZero(t)
		v := rand.Int31n(64) - 32

		// Compute subtraction via scalar type.
		var k scalar
		k.SetIntRaw(x)
		k.SubInt32(v)
		got := k.IntRaw()

		// Compute expectation.
		expect := new(big.Int).Sub(x, new(big.Int).SetInt64(int64(v)))

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}

func TestScalarRsh(t *testing.T) {
	for trial := 0; trial < ConstNumTrials; trial++ {
		x := RandScalarNonZero(t)
		s := uint(1 + rand.Intn(63))

		// Compute shift via scalar type.
		var k scalar
		k.SetIntRaw(x)
		k.Rsh(s)
		got := k.IntRaw()

		// Compute expectation.
		expect := new(big.Int).Rsh(x, s)

		if got.Cmp(expect) != 0 {
			t.FailNow()
		}
	}
}
