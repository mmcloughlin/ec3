package verif

import (
	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/arith/eval"
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/z3"
)

// Evaluator can evaluate arithmetic programs over bit-vector integers.
type Evaluator struct {
	eval *eval.Evaluator
}

func NewEvaluator(ctx *z3.Context, n uint) *Evaluator {
	proc := newprocessor(ctx, n)
	return &Evaluator{
		eval: eval.NewEvaluator(proc),
	}
}

func (e *Evaluator) SetRegister(r ir.Register, x *z3.BV) {
	e.eval.SetRegister(r, x)
}

func (e *Evaluator) Register(r ir.Register) (*z3.BV, error) {
	v, err := e.eval.Register(r)
	if err != nil {
		return nil, err
	}
	return v.(*z3.BV), nil
}

func (e *Evaluator) Execute(p *ir.Program) error {
	return e.eval.Execute(p)
}

type processor struct {
	*z3.BVSort

	ctx  *z3.Context
	errs errutil.Errors
}

func newprocessor(ctx *z3.Context, n uint) *processor {
	return &processor{
		BVSort: ctx.BVSort(n),
		ctx:    ctx,
	}
}

func (p *processor) Errors() []error { return p.errs }

func (p *processor) Const(x uint64, n uint) eval.Value { return p.ctx.BVSort(n).Uint64(x) }

func (p *processor) ADD(x, y, cin eval.Value) (sum, cout eval.Value) {
	n := p.Bits()

	// Extend all operands to n+1 bits.
	xext := p.word(x).ZeroExt(1)
	yext := p.word(y).ZeroExt(1)
	cext := p.bit(cin).ZeroExt(n)

	// Compute extended sum and extract the carry.
	sumext := xext.Add(yext).Add(cext)
	sum = sumext.Extract(n-1, 0)
	cout = sumext.Extract(n, n)
	return
}

func (p *processor) SUB(x, y, bin eval.Value) (diff, bout eval.Value) {
	// Implement subtraction using the adder.
	noty := p.word(y).Not()
	cin := p.bit(bin).Not()
	sum, cout := p.ADD(x, noty, cin)
	diff = sum
	bout = p.bit(cout).Not()
	return
}

func (p *processor) MUL(x, y eval.Value) (hi, lo eval.Value) {
	n := p.Bits()
	xext := p.word(x).ZeroExt(n)
	yext := p.word(y).ZeroExt(n)
	z := xext.Mul(yext)
	hi = z.Extract(2*n-1, n)
	lo = z.Extract(n-1, 0)
	return
}

func (p *processor) SHL(x eval.Value, s uint) eval.Value {
	return p.word(x).Shl(p.BVSort.Uint64(uint64(s)))
}

func (p *processor) SHR(x eval.Value, s uint) eval.Value {
	return p.word(x).LogicShr(p.BVSort.Uint64(uint64(s)))
}

func (p *processor) word(v eval.Value) *z3.BV {
	return p.bv(v, p.Bits())
}

func (p *processor) bit(v eval.Value) *z3.BV {
	return p.bv(v, 1)
}

func (p *processor) bv(v eval.Value, n uint) *z3.BV {
	x, ok := v.(*z3.BV)
	if !ok {
		p.errs.Add(errutil.UnexpectedType(v))
		return p.zero(n)
	}
	if x.Bits() != n {
		p.errs.Add(xerrors.Errorf("incorrect bit-vector size: got %d require %d", x.Bits(), n))
		return p.zero(n)
	}
	return x
}

func (p *processor) zero(n uint) *z3.BV {
	return p.ctx.BVSort(n).Uint64(0)
}
