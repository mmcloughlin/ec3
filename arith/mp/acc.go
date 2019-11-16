package mp

import (
	"math/big"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/ints"
)

type Accumulator struct {
	dst ir.Int        // destination
	x   ir.Operands   // current state
	c   []ir.Operands // carry bits at each position
	max *big.Int      // maximum possible value

	s   uint
	ctx *build.Context
}

func NewAccumulator(ctx *build.Context, dst ir.Int) *Accumulator {
	return &Accumulator{
		dst: dst,
		max: new(big.Int),

		s:   64,
		ctx: ctx,
	}
}

// Int returns the value in the accumulator.
func (a *Accumulator) Int() ir.Int { return a.x }

// AddAtMax adds a into the accumulator starting at limb i, assuming that the
// registers have maximum value max.
func (a *Accumulator) AddAtMax(y ir.Int, max *big.Int, i int) {
	// Add one limb at a time.
	M := bigint.Clone(max)
	for _, limb := range y.Limbs() {
		// Determine maximum value of this limb.
		maxsum := bigint.Min(bigint.Clone(M), a.maxword())
		M.Rsh(M, a.s)

		// Determine additive and it's maximim value.
		var additive, sum ir.Operand
		if x := a.get(i); x != nil {
			additive = x
			sum = x
			maxsum.Add(maxsum, a.maxlimb(i))
		} else {
			additive = ir.Constant(0)
			sum = a.dst.Limb(i)
			a.set(i, sum)
		}

		// Determine carry in.
		cin := a.carry(i)
		if cin != nil {
			maxsum.Add(maxsum, bigint.One())
		} else {
			cin = ir.Flag(0)
		}

		// Is a carry out possible?
		cout := ir.Discard
		if maxsum.BitLen() > int(a.s) {
			cout = a.ctx.Register("c")
			a.setcarry(i+1, cout)
		}

		a.ctx.ADD(additive, limb, cin, sum, cout)

		i++
	}

	// Update maximum value.
	Y := new(big.Int).Lsh(max, a.s*uint(i))
	a.max.Add(a.max, Y)
}

// AddAt adds a into the accumulator starting at limb i. Assumes the maximum possible value of y.
func (a *Accumulator) AddAt(y ir.Int, i int) {
	bits := a.s * uint(y.Len())
	a.AddAtMax(y, bigint.Ones(bits), i)
}

// AddProduct is a convenience for adding the high and low parts of a product into the accumulator.
func (a *Accumulator) AddProduct(hi, lo ir.Operand, i int) {
	maxprod := new(big.Int).Mul(a.maxword(), a.maxword())
	a.AddAtMax(ir.Operands{lo, hi}, maxprod, i)
}

func (a *Accumulator) Flush() {
	for i := 0; i < a.limbs(); i++ {
		for {
			c := a.carry(i)
			if c == nil {
				break
			}

			x := a.get(i)
			if x == nil {
				x = a.dst.Limb(i)
				a.set(i, x)
			}

			a.ctx.ADD(ir.Constant(0), ir.Constant(0), c, x, a.ctx.Register("c"))
		}
	}
}

// maxword returns the largest possible value for a single word.
func (a *Accumulator) maxword() *big.Int {
	return bigint.Ones(a.s)
}

// maxlimb returns an upper bound on the value of limb i.
func (a *Accumulator) maxlimb(i int) *big.Int {
	maxshift := new(big.Int).Rsh(a.max, a.s*uint(i))
	return bigint.Min(maxshift, a.maxword())
}

// limbs returns the required number of limbs to represent the max possible value.
func (a *Accumulator) limbs() int {
	return ints.NextMultiple(a.max.BitLen(), int(a.s)) / int(a.s)
}

func (a *Accumulator) get(i int) ir.Operand {
	if i < len(a.x) {
		return a.x[i]
	}
	return nil
}

func (a *Accumulator) set(i int, v ir.Operand) {
	for len(a.x) <= i {
		a.x = append(a.x, nil)
	}
	a.x[i] = v
}

// carry gets the carry into limb i.
func (a *Accumulator) carry(i int) ir.Operand {
	if i < len(a.c) && len(a.c[i]) > 0 {
		c := a.c[i][0]
		a.c[i] = a.c[i][1:]
		return c
	}
	return nil
}

func (a *Accumulator) setcarry(i int, c ir.Operand) {
	for len(a.c) <= i {
		a.c = append(a.c, nil)
	}
	a.c[i] = append(a.c[i], c)
}
