package ir

import (
	"math/big"
	"strconv"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

// Int represents a multi-precision integer.
type Int interface {
	Len() int
	Limb(int) Operand
	Limbs() Operands
}

// Limb returns the i'th limb of x, or zero if i ⩾ x.Len().
func Limb(x Int, i int) Operand {
	if i < x.Len() {
		return x.Limb(i)
	}
	return Zero
}

// Clone returns a copy of x.
func Clone(x Int) Operands {
	return append(Operands{}, x.Limbs()...)
}

// Extend returns a new integer with limbs appended to the end.
func Extend(x Int, limbs ...Operand) Operands {
	return append(Clone(x), limbs...)
}

// Operands is a multi-precision integer formed of operands.
type Operands []Operand

func (o Operands) Len() int           { return len(o) }
func (o Operands) Limb(i int) Operand { return o[i] }
func (o Operands) Limbs() Operands    { return o }

// Registers is a multi-precision integer formed of registers.
type Registers []Register

func NewRegisters(prefix string, k uint) Registers {
	x := make(Registers, k)
	for i := 0; i < int(k); i++ {
		x[i] = Register(prefix + strconv.Itoa(i))
	}
	return x
}

func (r Registers) Len() int           { return len(r) }
func (r Registers) Limb(i int) Operand { return r[i] }

func (r Registers) Limbs() Operands {
	operands := make(Operands, len(r))
	for i := range r {
		operands[i] = r[i]
	}
	return operands
}

// Constants is a multi-precision integer formed of constants.
type Constants []Constant

// NewConstantsFromInt represents x as s-bit integers.
func NewConstantsFromInt(x *big.Int, s uint) Constants {
	c := Constants{}
	mask := bigint.Ones(s)
	for bigint.IsNonZero(x) {
		limb := new(big.Int).And(x, mask).Uint64()
		c = append(c, Constant(limb))
		x.Rsh(x, s)
	}
	return c
}

func (c Constants) Len() int           { return len(c) }
func (c Constants) Limb(i int) Operand { return c[i] }

func (c Constants) Limbs() Operands {
	operands := make(Operands, len(c))
	for i := range c {
		operands[i] = c[i]
	}
	return operands
}
