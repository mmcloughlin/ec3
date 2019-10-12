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

func NewEvaluator() *Evaluator {
	return &Evaluator{
		mem: map[string]uint64{},
	}
}

func (e *Evaluator) Register(r ir.Register) (uint64, error) {
	defer e.reseterror()
	return e.register(r), e.errs.Err()
}

func (e *Evaluator) Execute(p *ir.Program) error {
	defer e.reseterror()
	for _, i := range p.Instructions {
		if err := e.instruction(i); err != nil {
			return err
		}
	}
	return nil
}

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
	// case ir.SUB:
	// case ir.MUL:
	// case ir.SHL:
	// case ir.SHR:
	default:
		return errutil.UnexpectedType(i)
	}
	return e.errs.Err()
}

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

func (e *Evaluator) register(r ir.Register) uint64 {
	return e.load(string(r))
}

func (e *Evaluator) flag(op ir.Operand) uint64 {
	b := e.operand(op)
	e.assertwidth(b, 1)
	return b
}

func (e *Evaluator) load(name string) uint64 {
	x, ok := e.mem[name]
	if !ok {
		e.errorf("operand %q undefined", name)
		return 0
	}
	return x
}

func (e *Evaluator) setregister(r ir.Register, x uint64) {
	e.store(string(r), x)
}

func (e *Evaluator) setflag(r ir.Register, b uint64) {
	e.assertwidth(b, 1)
	e.setregister(r, b)
}

func (e *Evaluator) store(name string, x uint64) {
	e.mem[name] = x
}

func (e *Evaluator) assertwidth(x uint64, w uint) {
	if (x >> w) != 0 {
		e.errorf("expected to be %d bits wide", w)
	}
}

func (e *Evaluator) reseterror() { e.errs = nil }

func (e *Evaluator) adderror(err error) { e.errs.Add(err) }

func (e *Evaluator) errorf(format string, args ...interface{}) {
	e.adderror(xerrors.Errorf(format, args...))
}
