package ir

// Operand is an interface for instruction operands.
type Operand interface {
	operand()
}

// Declare members of the sealed interface.
func (Register) operand() {}
func (Flag) operand()     {}
func (Constant) operand() {}

// Register is a machine word operand.
type Register string

// Flag is a boolean (single-bit) operand.
type Flag string

// Constant is a constant operand.
type Constant uint64

// MOV is a move instruction.
type MOV struct {
	Source      Operand
	Destination Register
}

// ADD is an add with carry instruction.
type ADD struct {
	X        Operand
	Y        Operand
	CarryIn  Flag
	Result   Register
	CarryOut Flag
}

// SUB is an subtract with borrow instruction.
type SUB struct {
	X         Operand
	Y         Operand
	BorrowIn  Flag
	Result    Register
	BorrowOut Flag
}

// MUL is a multiply instruction providing lower and upper parts of the result.
type MUL struct {
	X     Operand
	Y     Operand
	Upper Register
	Lower Register
}

// SHL is a shift left instruction.
type SHL struct {
	X      Operand
	Shift  Constant
	Result Register
}

// SHR is a shift right instruction.
type SHR struct {
	X      Operand
	Shift  Constant
	Result Register
}
