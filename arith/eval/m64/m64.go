// Package m64 provides a 64-bit arithmetic evaluator.
package m64

import (
	"math/big"
	"math/bits"

	"github.com/mmcloughlin/ec3/arith/eval"
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Word is a 64-bit machine word.
type Word uint64

// Bits returns the number of bits required to represent x.
func (x Word) Bits() uint {
	return uint(bits.Len64(uint64(x)))
}

// Evaluator for arithmetic programs with 64-bit limbs.
type Evaluator struct {
	eval *eval.Evaluator
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		eval: eval.NewEvaluator(New()),
	}
}

// SetRegister sets register r to value x.
func (e *Evaluator) SetRegister(r ir.Register, x uint64) {
	e.eval.SetRegister(r, Word(x))
}

// Register returns the value in the given register.
func (e *Evaluator) Register(r ir.Register) (uint64, error) {
	v, err := e.eval.Register(r)
	if err != nil {
		return 0, err
	}
	return u64(v)
}

// SetInt sets registers to the 64-bit limbs of x.
func (e *Evaluator) SetInt(z ir.Registers, x *big.Int) {
	limbs := bigint.Uint64s(x)
	for i, limb := range limbs {
		e.SetRegister(z[i], limb)
	}
	for i := len(limbs); i < len(z); i++ {
		e.SetRegister(z[i], 0)
	}
}

// Int returns the integer represented by the 64-bit limbs in the given registers.
func (e *Evaluator) Int(z ir.Registers) (*big.Int, error) {
	words := make([]uint64, len(z))
	for i, r := range z {
		word, err := e.Register(r)
		if err != nil {
			return nil, err
		}
		words[i] = word
	}
	return bigint.FromUint64s(words), nil
}

// Execute the program p.
func (e *Evaluator) Execute(p *ir.Program) error {
	return e.eval.Execute(p)
}

// Processor is a 64-bit arithmetic evaluator.
type Processor struct {
	errs errutil.Errors
}

// New builds a new 64-bit arithmetic processor.
func New() *Processor {
	return &Processor{}
}

// Errors returns any errors encountered during execution.
func (p *Processor) Errors() []error {
	return p.errs
}

// Bits returns the word size.
func (Processor) Bits() uint { return 64 }

// Const builds an n-bit constant.
func (Processor) Const(x uint64, n uint) eval.Value { return Word(x) }

// ITE returns x if l≡r else y.
func (p *Processor) ITE(l, r, x, y eval.Value) eval.Value {
	if p.u64(l) == p.u64(r) {
		return x
	}
	return y
}

// ADD executes an add with carry instruction.
func (p *Processor) ADD(x, y, cin eval.Value) (sum, cout eval.Value) {
	s, c := bits.Add64(p.u64(x), p.u64(y), p.u64(cin))
	return Word(s), Word(c)
}

// SUB executes a subtract with borrow instruction.
func (p *Processor) SUB(x, y, bin eval.Value) (diff, bout eval.Value) {
	d, b := bits.Sub64(p.u64(x), p.u64(y), p.u64(bin))
	return Word(d), Word(b)
}

// MUL executes a multiply instruction.
func (p *Processor) MUL(x, y eval.Value) (hi, lo eval.Value) {
	h, l := bits.Mul64(p.u64(x), p.u64(y))
	return Word(h), Word(l)
}

// SHL executes a shift left instruction.
func (p *Processor) SHL(x eval.Value, s uint) eval.Value {
	return Word(p.u64(x) << s)
}

// SHR executes a shift right instruction.
func (p *Processor) SHR(x eval.Value, s uint) eval.Value {
	return Word(p.u64(x) >> s)
}

// u64 casts v to uint64.
func (p *Processor) u64(v eval.Value) uint64 {
	x, err := u64(v)
	if err != nil {
		p.errs.Add(errutil.UnexpectedType(v))
		return 0
	}
	return x
}

// u64 type asserts v to a uint64.
func u64(v eval.Value) (uint64, error) {
	if x, ok := v.(Word); ok {
		return uint64(x), nil
	}
	return 0, errutil.UnexpectedType(v)
}
