package addchain

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/ec3/prime"

	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/test"
)

func TestChainAlgorithms(t *testing.T) {
	as := []ChainAlgorithm{
		BinaryRightToLeft{},

		// Dictionary-based algorithms.
		NewDictAlgorithm(
			SlidingWindow{K: 4},
			NewContinuedFractions(DichotomicStrategy{}),
		),
		NewDictAlgorithm(
			FixedWindow{K: 7},
			NewContinuedFractions(BinaryStrategy{}),
		),

		// Runs algorithm.
		NewRunsAlgorithm(NewContinuedFractions(DichotomicStrategy{})),
	}
	for _, a := range as {
		t.Run(a.String(), ChainAlgorithmSuite(a))
	}
}

func ChainAlgorithmSuite(a ChainAlgorithm) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("powers_of_two", CheckPowersOfTwo(a, 100))
		t.Run("binary_runs", CheckBinaryRuns(a, 32))
		t.Run("random_int64", CheckRandomInt64s(a))
		t.Run("primes", CheckPrimes(a))
	}
}

func CheckPowersOfTwo(a ChainAlgorithm, e uint) func(t *testing.T) {
	return func(t *testing.T) {
		n := big.NewInt(1)
		for i := uint(0); i <= e; i++ {
			AssertChainAlgorithmProduces(t, a, n)
			n.Lsh(n, 1)
		}
	}
}

func CheckBinaryRuns(a ChainAlgorithm, n uint) func(t *testing.T) {
	return func(t *testing.T) {
		for i := uint(1); i <= n; i++ {
			r := bigint.Pow2(i)
			r.Sub(r, bigint.One())
			AssertChainAlgorithmProduces(t, a, r)
		}
	}
}

func CheckRandomInt64s(a ChainAlgorithm) func(t *testing.T) {
	return test.Trials(func(t *testing.T) bool {
		r := rand.Int63n(math.MaxInt64)
		n := big.NewInt(r)
		AssertChainAlgorithmProduces(t, a, n)
		return true
	})
}

func CheckPrimes(a ChainAlgorithm) func(t *testing.T) {
	// Prepare primes in a random order.
	ps := []*big.Int{}
	for _, p := range prime.Distinguished {
		ps = append(ps, p.Int())
	}
	rand.Shuffle(len(ps), func(i, j int) { ps[i], ps[j] = ps[j], ps[i] })

	return test.Trials(func(t *testing.T) bool {
		AssertChainAlgorithmProduces(t, a, ps[0])
		ps = ps[1:]
		return len(ps) > 0
	})
}
