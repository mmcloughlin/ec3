package addchain

import (
	"fmt"
	"math/big"
)

type RunsAlgorithm struct {
	seqalg SequenceAlgorithm
}

func NewRunsAlgorithm(a SequenceAlgorithm) *RunsAlgorithm {
	return &RunsAlgorithm{
		seqalg: a,
	}
}

func (a RunsAlgorithm) String() string {
	return fmt.Sprintf("runs(%s)", a.seqalg)
}

func (a RunsAlgorithm) FindChain(n *big.Int) (Chain, error) {
	/*
		// Find the runs in n.
		d := RunLength{T: 0}
		sum := d.Decompose(n)
		runs := sum.Dictionary()

		// Treat the run lengths themselves as a sequence to be solved.
		lengths := []*big.Int{}
		for _, run := range runs {
			length := int64(run.BitLen())
			lengths = append(lengths, big.NewInt(length))
		}

		// Delegate to the sequence algorithm for a solution.
		lc, err := a.seqalg.FindSequence(lengths)
		if err != nil {
			return nil, err
		}
	*/

	return nil, nil
}
