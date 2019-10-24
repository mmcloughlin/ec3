package ir

import (
	"strconv"
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
