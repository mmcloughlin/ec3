// Package ir defines an intermediate representation for multi-precision arithmetic.
package ir

import (
	"fmt"
	"reflect"
	"strings"
)

// Operand is an interface for instruction operands.
type Operand interface {
	fmt.Stringer

	operand() // sealed
}

// Register is a machine word operand.
type Register string

func (r Register) String() string { return string(r) }

func (Register) operand() {}

// SelectRegisters selects the registers from the list of operands.
func SelectRegisters(ops []Operand) []Register {
	var rs []Register
	for _, op := range ops {
		if r, ok := op.(Register); ok {
			rs = append(rs, r)
		}
	}
	return rs
}

// Constant is a constant operand.
type Constant uint64

func (c Constant) String() string { return fmt.Sprintf("$%#x", uint64(c)) }

func (Constant) operand() {}

// Zero is the zero constant.
var Zero = Constant(0)

// Flag is a single-bit constant operand.
type Flag uint64

func (f Flag) String() string { return fmt.Sprintf("$%d", uint64(f)) }

func (Flag) operand() {}

// Program is a sequence of instructions.
type Program struct {
	Instructions []Instruction
}

func (p *Program) String() string {
	s := ""
	for _, i := range p.Instructions {
		s += FormatInstruction(i) + "\n"
	}
	return s
}

// Instruction in the intermediate representation.
type Instruction interface {
	Operands() []Operand

	instruction() // sealed
}

// FormatInstruction returns a string representation of the instruction.
func FormatInstruction(i Instruction) string {
	mneumonic := reflect.TypeOf(i).Name()
	ops := []string{}
	for _, op := range i.Operands() {
		ops = append(ops, op.String())
	}
	return mneumonic + "\t" + strings.Join(ops, ", ")
}

// MOV is a move instruction.
type MOV struct {
	Source      Operand
	Destination Register
}

func (i MOV) Operands() []Operand {
	return []Operand{i.Source, i.Destination}
}

func (MOV) instruction() {}

// CMOV is a conditional move.
type CMOV struct {
	Source      Operand
	Destination Register
	Flag        Operand
	Equals      Flag
}

func (i CMOV) Operands() []Operand {
	return []Operand{i.Source, i.Destination, i.Flag}
}

func (CMOV) instruction() {}

// ADD is an add with carry instruction.
type ADD struct {
	X        Operand
	Y        Operand
	CarryIn  Operand
	Sum      Register
	CarryOut Register
}

func (i ADD) Operands() []Operand {
	return []Operand{i.X, i.Y, i.CarryIn, i.Sum, i.CarryOut}
}

func (ADD) instruction() {}

// SUB is an subtract with borrow instruction.
type SUB struct {
	X         Operand
	Y         Operand
	BorrowIn  Operand
	Diff      Register
	BorrowOut Register
}

func (i SUB) Operands() []Operand {
	return []Operand{i.X, i.Y, i.BorrowIn, i.Diff, i.BorrowOut}
}

func (SUB) instruction() {}

// MUL is a multiply instruction providing lower and upper parts of the result.
type MUL struct {
	X    Operand
	Y    Operand
	High Register
	Low  Register
}

func (i MUL) Operands() []Operand {
	return []Operand{i.X, i.Y, i.High, i.Low}
}

func (MUL) instruction() {}

// SHL is a shift left instruction.
type SHL struct {
	X      Operand
	Shift  Constant
	Result Register
}

func (i SHL) Operands() []Operand {
	return []Operand{i.X, i.Result}
}

func (SHL) instruction() {}

// SHR is a shift right instruction.
type SHR struct {
	X      Operand
	Shift  Constant
	Result Register
}

func (i SHR) Operands() []Operand {
	return []Operand{i.X, i.Result}
}

func (SHR) instruction() {}
