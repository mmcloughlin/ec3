// Package eval provides an evaluator for arithmetic intermediate representation.
package eval

import (
	"math/bits"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Evaluator evaluates an arithmetic program in 64-bit.
type Evaluator struct {
	mem  map[string]uint64
	errs errutil.Errors
}

// NewEvaluator constructs an evaluator with empty state.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		mem: map[string]uint64{},
	}
}

// Register returns the value in the given register.
func (e *Evaluator) Register(r ir.Register) (uint64, error) {
	defer e.reseterror()
	return e.register(r), e.errs.Err()
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
		sum, cout := bits.Add64(x, y, cin)
		e.setregister(i.Sum, sum)
		e.setflag(i.CarryOut, cout)
	case ir.SUB:
		x := e.operand(i.X)
		y := e.operand(i.Y)
		bin := e.flag(i.BorrowIn)
		diff, bout := bits.Sub64(x, y, bin)
		e.setregister(i.Diff, diff)
		e.setflag(i.BorrowOut, bout)
	case ir.MUL:
		x := e.operand(i.X)
		y := e.operand(i.Y)
		hi, lo := bits.Mul64(x, y)
		e.setregister(i.High, hi)
		e.setregister(i.Low, lo)
	case ir.SHL:
		x := e.operand(i.X)
		e.setregister(i.Result, x<<i.Shift)
	case ir.SHR:
		x := e.operand(i.X)
		e.setregister(i.Result, x>>i.Shift)
	default:
		return errutil.UnexpectedType(i)
	}
	return e.errs.Err()
}

// operand returns the value of the given operand.
func (e *Evaluator) operand(operand ir.Operand) uint64 {
	switch op := operand.(type) {
	case ir.Register:
		return e.register(op)
	case ir.Constant:
		return uint64(op)
	default:
		e.adderror(errutil.UnexpectedType(op))
		return 0
	}
}

// register loads the value in register r.
func (e *Evaluator) register(r ir.Register) uint64 {
	return e.load(string(r))
}

// flag loads the value in the given flag operand.
func (e *Evaluator) flag(op ir.Operand) uint64 {
	b := e.operand(op)
	e.assertwidth(b, 1)
	return b
}

// load named value from memory.
func (e *Evaluator) load(name string) uint64 {
	x, ok := e.mem[name]
	if !ok {
		e.errorf("operand %q undefined", name)
		return 0
	}
	return x
}

// setregister sets register r to value x.
func (e *Evaluator) setregister(r ir.Register, x uint64) {
	e.store(string(r), x)
}

// setflag sets register r to the flab bit b.
func (e *Evaluator) setflag(r ir.Register, b uint64) {
	e.assertwidth(b, 1)
	e.setregister(r, b)
}

// store x at name.
func (e *Evaluator) store(name string, x uint64) {
	e.mem[name] = x
}

// assertwidth sets an error if x is not a w-bit value.
func (e *Evaluator) assertwidth(x uint64, w uint) {
	if (x >> w) != 0 {
		e.errorf("expected to be %d-bit value", w)
	}
}

// reseterror sets internal error state to nil.
func (e *Evaluator) reseterror() { e.errs = nil }

// adderror adds an error to the internal error list.
func (e *Evaluator) adderror(err error) { e.errs.Add(err) }

// errorf is a convenience for adding a formatted error to the internal list.
func (e *Evaluator) errorf(format string, args ...interface{}) {
	e.adderror(xerrors.Errorf(format, args...))
}
