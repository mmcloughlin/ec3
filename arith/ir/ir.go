// Package ir defines an intermediate representation for multi-precision arithmetic.
package ir

// Operand is an interface for instruction operands.
type Operand interface {
	operand()
}

// Register is a machine word operand.
type Register string

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

func (Constant) operand() {}

// Zero is the zero constant.
var Zero = Constant(0)

// Flag is a single-bit constant operand.
type Flag uint64

func (Flag) operand() {}

// Program is a sequence of instructions.
type Program struct {
	Instructions []Instruction
}

// Instruction in the intermediate representation.
type Instruction interface {
	Operands() []Operand

	instruction()
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
