// Package eval provides an evaluator for arithmetic intermediate representation.
package eval

import (
	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Value is a value stored in a register.
type Value interface {
	// Bits returns the number of bits required to represent the value.
	Bits() uint
}

// Processor is an implementation of the arithmetic instruction set.
type Processor interface {
	Bits() uint
	Const(x uint64, n uint) Value
	ADD(x, y, cin Value) (sum, cout Value)
	SUB(x, y, bin Value) (diff, bout Value)
	MUL(x, y Value) (hi, lo Value)
	SHL(x Value, s uint) Value
	SHR(x Value, s uint) Value
	Errors() []error
}

// Evaluator evaluates an arithmetic program.
type Evaluator struct {
	proc Processor
	mem  map[string]Value
	errs errutil.Errors
}

// NewEvaluator constructs an evaluator with empty state.
func NewEvaluator(proc Processor) *Evaluator {
	return &Evaluator{
		proc: proc,
		mem:  map[string]Value{},
	}
}

// SetRegister sets register r to value x.
func (e *Evaluator) SetRegister(r ir.Register, x Value) {
	e.setregister(r, x)
}

// Register returns the value in the given register.
func (e *Evaluator) Register(r ir.Register) (Value, error) {
	defer e.reseterror()
	return e.register(r), e.err()
}

// Execute the program p.
func (e *Evaluator) Execute(p *ir.Program) error {
	defer e.reseterror()
	for _, i := range p.Instructions {
		if err := e.instruction(i); err != nil {
			return err
		}
	}
	return nil
}

// instruction executes a single instruction.
func (e *Evaluator) instruction(inst ir.Instruction) error {
	switch i := inst.(type) {
	case ir.MOV:
		e.setregister(i.Destination, e.operand(i.Source))
	case ir.ADD:
		x := e.operand(i.X)
		y := e.operand(i.Y)
		cin := e.flag(i.CarryIn)
		sum, cout := e.proc.ADD(x, y, cin)
		e.setregister(i.Sum, sum)
		e.setflag(i.CarryOut, cout)
	case ir.SUB:
		x := e.operand(i.X)
		y := e.operand(i.Y)
		bin := e.flag(i.BorrowIn)
		diff, bout := e.proc.SUB(x, y, bin)
		e.setregister(i.Diff, diff)
		e.setflag(i.BorrowOut, bout)
	case ir.MUL:
		x := e.operand(i.X)
		y := e.operand(i.Y)
		hi, lo := e.proc.MUL(x, y)
		e.setregister(i.High, hi)
		e.setregister(i.Low, lo)
	case ir.SHL:
		x := e.operand(i.X)
		e.setregister(i.Result, e.proc.SHL(x, uint(i.Shift)))
	case ir.SHR:
		x := e.operand(i.X)
		e.setregister(i.Result, e.proc.SHR(x, uint(i.Shift)))
	default:
		return errutil.UnexpectedType(i)
	}
	return e.err()
}

// operand returns the value of the given operand.
func (e *Evaluator) operand(operand ir.Operand) Value {
	switch op := operand.(type) {
	case ir.Register:
		return e.register(op)
	case ir.Constant:
		return e.proc.Const(uint64(op), e.proc.Bits())
	case ir.Flag:
		return e.proc.Const(uint64(op), 1)
	default:
		e.adderror(errutil.UnexpectedType(op))
		return nil
	}
}

// register loads the value in register r.
func (e *Evaluator) register(r ir.Register) Value {
	return e.load(string(r))
}

// flag loads the value in the given flag operand.
func (e *Evaluator) flag(op ir.Operand) Value {
	b := e.operand(op)
	e.assertwidth(b, 1)
	return b
}

// load named value from memory.
func (e *Evaluator) load(name string) Value {
	x, ok := e.mem[name]
	if !ok {
		e.errorf("operand %q undefined", name)
		return nil
	}
	return x
}

// setregister sets register r to value x.
func (e *Evaluator) setregister(r ir.Register, x Value) {
	e.store(string(r), x)
}

// setflag sets register r to the flab bit b.
func (e *Evaluator) setflag(r ir.Register, b Value) {
	e.assertwidth(b, 1)
	e.setregister(r, b)
}

// store x at name.
func (e *Evaluator) store(name string, x Value) {
	e.mem[name] = x
}

// assertwidth sets an error if x is not a w-bit value.
func (e *Evaluator) assertwidth(x Value, w uint) {
	if x.Bits() > w {
		e.errorf("expected to be %d-bit value", w)
	}
}

// err returns any accumulated errors.
func (e *Evaluator) err() error {
	var errs errutil.Errors
	errs.Add(e.errs...)
	errs.Add(e.proc.Errors()...)
	return errs.Err()
}

// reseterror sets internal error state to nil.
func (e *Evaluator) reseterror() { e.errs = nil }

// adderror adds an error to the internal error list.
func (e *Evaluator) adderror(err error) { e.errs.Add(err) }

// errorf is a convenience for adding a formatted error to the internal list.
func (e *Evaluator) errorf(format string, args ...interface{}) {
	e.adderror(xerrors.Errorf(format, args...))
}
