package addchain

import (
	"math/big"
	"testing"

	"github.com/mmcloughlin/ec3/prime"
)

func TestRuns(t *testing.T) {
	seqalg := NewContinuedFractions(DichotomicStrategy{})
	a := NewRunsAlgorithm(seqalg)

	p := prime.NISTP256.Int()
	q := new(big.Int).Sub(p, big.NewInt(3))
	c, err := a.FindChain(q)
	if err != nil {
		t.Fatal(err)
	}

	for i, x := range c {
		t.Logf("%3d: %x", i+1, x)
	}

	if err := c.Produces(q); err != nil {
		t.Fatal(err)
	}
}

func TestRunsEnsemble(t *testing.T) {
	seqalgs := []SequenceAlgorithm{
		// Continued fractions algorithms.
		NewContinuedFractions(BinaryStrategy{}),
		NewContinuedFractions(BinaryStrategy{Parity: 1}),
		NewContinuedFractions(DichotomicStrategy{}),

		// Heuristics algorithms.
		NewHeuristicAlgorithm(UseFirstHeuristic{
			Halving{},
			DeltaLargest{},
		}),
		NewHeuristicAlgorithm(UseFirstHeuristic{
			Halving{},
			Approximation{},
		}),
	}

	as := []ChainAlgorithm{}
	for _, seqalg := range seqalgs {
		a := Optimized{NewRunsAlgorithm(seqalg)}
		as = append(as, a)
	}

	p := prime.NISTP256.Int()
	q := new(big.Int).Sub(p, big.NewInt(3))

	for _, a := range as {
		t.Log(a)
		c, err := a.FindChain(q)
		if err != nil {
			t.Fatal(err)
		}

		/*
			for i, x := range c {
				t.Logf("%3d: %x", i+1, x)
			}
		*/

		if err := c.Produces(q); err != nil {
			t.Fatal(err)
		}

		p, err := c.Program()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("adds=%d doubles=%d total=%d", p.Adds(), p.Doubles(), p.Adds()+p.Doubles())
	}
}
