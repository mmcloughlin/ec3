// Package ir defines an intermediate representation for multi-precision arithmetic.
package ir

// Operand is an interface for instruction operands.
type Operand interface {
	operand()
}

// Register is a machine word operand.
type Register string

func (Register) operand() {}

// Constant is a constant operand.
type Constant uint64

func (Constant) operand() {}

// Program is a sequence of instructions.
type Program struct {
	Instructions []Instruction
}

// Instruction in the intermediate representation.
type Instruction interface {
	instruction()
}

// MOV is a move instruction.
type MOV struct {
	Source      Operand
	Destination Register
}

func (MOV) instruction() {}

// ADD is an add with carry instruction.
type ADD struct {
	X        Operand
	Y        Operand
	CarryIn  Operand
	Sum      Register
	CarryOut Register
}

func (ADD) instruction() {}

// SUB is an subtract with borrow instruction.
type SUB struct {
	X         Operand
	Y         Operand
	BorrowIn  Operand
	Result    Register
	BorrowOut Register
}

func (SUB) instruction() {}

// MUL is a multiply instruction providing lower and upper parts of the result.
type MUL struct {
	X     Operand
	Y     Operand
	Upper Register
	Lower Register
}

func (MUL) instruction() {}

// SHL is a shift left instruction.
type SHL struct {
	X      Operand
	Shift  Constant
	Result Register
}

func (SHL) instruction() {}

// SHR is a shift right instruction.
type SHR struct {
	X      Operand
	Shift  Constant
	Result Register
}

func (SHR) instruction() {}
