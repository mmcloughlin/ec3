package addchain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/bigints"
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

	/*
		fmt.Println("lengths chain:")
		DumpChain(lc)
	*/

	// Build a dictionary chain from this.
	c, err := RunsChain(lc)
	if err != nil {
		return nil, err
	}

	/*
			fmt.Println("runs chain:")
			DumpChain(c)

		// Reduce.
		sum, c, err = PrimitiveDictionary(sum, c)
		if err != nil {
			return nil, err
		}
	*/

	// Build chain for n out of the dictionary.
	k := len(sum) - 1
	cur := bigint.Clone(sum[k].D)
	for ; k > 0; k-- {
		// Shift until the next exponent.
		for i := sum[k].E; i > sum[k-1].E; i-- {
			cur.Lsh(cur, 1)
			c.AppendClone(cur)
		}

		// Add in the dictionary term at this position.
		cur.Add(cur, sum[k-1].D)
		c.AppendClone(cur)
	}

	for i := sum[0].E; i > 0; i-- {
		cur.Lsh(cur, 1)
		c.AppendClone(cur)
	}

	// Prepare chain for returning.
	bigints.Sort(c)
	c = Chain(bigints.Unique(c))

	// DumpChain(c)

	return c, nil
}

func RunsChain(lc Chain) (Chain, error) {
	p, err := lc.Program()
	if err != nil {
		return nil, err
	}

	c := New()
	for _, op := range p {
		a, b := bigint.MinMax(lc[op.I], lc[op.J])
		if !a.IsUint64() || !b.IsUint64() {
			return nil, errors.New("values in lengths chain are too large")
		}

		la := uint(a.Uint64())
		lb := uint(b.Uint64())

		ra := bigint.Ones(la)
		rb := bigint.Ones(lb)

		cur := bigint.Clone(rb)
		for s := uint(0); s < la; s++ {
			cur.Lsh(cur, 1)
			c.AppendClone(cur)
		}

		cur.Add(cur, ra)
		c.AppendClone(cur)
	}

	return c, nil
}
