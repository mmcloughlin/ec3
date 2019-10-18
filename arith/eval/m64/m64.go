// Package m64 provides a 64-bit arithmetic evaluator.
package m64

import (
	"math/bits"

	"github.com/mmcloughlin/ec3/arith/eval"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Word is a 64-bit machine word.
type Word uint64

// Bits returns the number of bits required to represent x.
func (x Word) Bits() uint {
	return uint(bits.Len64(uint64(x)))
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

// Uint64 returns x as a processor word.
func (Processor) Uint64(x uint64) eval.Value { return Word(x) }

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
	if x, ok := v.(Word); ok {
		return uint64(x)
	}
	p.errs.Add(errutil.UnexpectedType(v))
	return 0
}
