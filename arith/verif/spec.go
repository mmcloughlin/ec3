package verif

import (
	"math/big"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/z3"
)

// Spec specifies an operation and facilitates equivalence proofs.
type Spec struct {
	sig *ir.Signature
	s   uint

	vars map[string]*z3.BV
	pre  []*z3.Bool

	ctx   *z3.Context
	model *z3.Model
}

func NewSpec(ctx *z3.Context, sig *ir.Signature, s uint) *Spec {
	return &Spec{
		sig:  sig,
		s:    s,
		vars: make(map[string]*z3.BV),
		ctx:  ctx,
	}
}

// Param returns a variable for the named parameter.
func (s *Spec) Param(name string) (*z3.BV, error) {
	if x, ok := s.vars[name]; ok {
		return x, nil
	}
	v := s.sig.Param(name)
	if v == nil {
		return nil, xerrors.Errorf("unknown parameter %q", name)
	}
	x, err := s.variable(v)
	if err != nil {
		return nil, err
	}
	s.vars[name] = x
	return x, nil
}

func (s *Spec) variable(v *ir.Var) (*z3.BV, error) {
	switch t := v.Type.(type) {
	case ir.Integer:
		sort := s.ctx.BVSort(t.K * s.s)
		return sort.Const(v.Name), nil
	default:
		return nil, errutil.UnexpectedType(t)
	}
}

// Registers returns registers mapped to the named parameter or result variable.
func (s *Spec) Registers(name string) (ir.Registers, error) {
	p := s.sig.Var(name)
	if p == nil {
		return nil, xerrors.Errorf("unknown variable %q", name)
	}
	switch t := p.Type.(type) {
	case ir.Integer:
		return ir.NewRegisters(strings.ToUpper(p.Name), t.K), nil
	default:
		return nil, errutil.UnexpectedType(t)
	}
}

// Result returns the named result.
func (s *Spec) Result(name string) (*z3.BV, error) {
	if s.sig.Result(name) == nil {
		return nil, xerrors.Errorf("unknown result %q", name)
	}

	x, ok := s.vars[name]
	if !ok {
		return nil, xerrors.Errorf("result %q not set", name)
	}

	return x, nil
}

// SetResult sets the named result.
func (s *Spec) SetResult(name string, x *z3.BV) error {
	if s.sig.Result(name) == nil {
		return xerrors.Errorf("unknown result %q", name)
	}

	if _, ok := s.vars[name]; ok {
		return xerrors.Errorf("result %q already set", name)
	}

	s.vars[name] = x
	return nil
}

// AddPrecondition sets a precondition for the function, a predicate which must
// be true prior to function execution.
func (s *Spec) AddPrecondition(cond *z3.Bool) {
	s.pre = append(s.pre, cond)
}

// Prove the program p meets the specification.
func (s *Spec) Prove(p *ir.Program) (bool, error) {
	e := NewEvaluator(s.ctx, s.s)

	// Configure inputs.
	for _, param := range s.sig.Params {
		switch t := param.Type.(type) {
		case ir.Integer:
			x, err := s.Param(param.Name)
			if err != nil {
				return false, err
			}
			regs, err := s.Registers(param.Name)
			if err != nil {
				return false, err
			}
			e.SetInt(regs, x)
		default:
			return false, errutil.UnexpectedType(t)
		}
	}

	// Execute program.
	if err := e.Execute(p); err != nil {
		return false, err
	}

	// Compare outputs.
	thms := []*z3.Bool{}
	for _, result := range s.sig.Results {
		switch t := result.Type.(type) {
		case ir.Integer:
			expect, err := s.Result(result.Name)
			if err != nil {
				return false, err
			}

			regs, err := s.Registers(result.Name)
			if err != nil {
				return false, err
			}

			got, err := e.Int(regs)
			if err != nil {
				return false, err
			}

			thm := got.Eq(expect)
			thms = append(thms, thm)
		default:
			return false, errutil.UnexpectedType(t)
		}
	}

	equiv := s.ctx.True().And(thms...)

	// Pass problem to solver.
	solver := s.ctx.SolverForLogic("QF_BV")
	defer solver.Close()

	// params := s.ctx.Params()
	// params.SetUint("threads", uint(runtime.NumCPU()))
	// solver.SetParams(params)

	for _, cond := range s.pre {
		solver.Assert(cond)
	}

	result, err := solver.Prove(equiv)
	if err != nil {
		return result, err
	}

	// Obtain a model if proving failed.
	s.model = nil
	if !result {
		s.model = solver.Model()
	}

	return result, nil
}

// Model returns the failing model, if a prior call to Prove returned false. Otherwise returns nil.
func (s *Spec) Model() *z3.Model {
	return s.model
}

// Counterexample returns a set of failing assignments, if a previous call to Prove
// returned false.
func (s *Spec) Counterexample() (map[string]*big.Int, error) {
	m := s.Model()
	if m == nil {
		return nil, xerrors.New("no counterexample available")
	}

	vs, err := m.Assignments()
	if err != nil {
		return nil, err
	}

	assign := map[string]*big.Int{}
	for name, v := range vs {
		if v == nil {
			assign[name] = nil
			continue
		}

		x, ok := v.(*z3.BV)
		if !ok {
			return nil, errutil.UnexpectedType(v)
		}

		i, ok := x.Int()
		if !ok {
			return nil, errutil.AssertionFailure("unable to convert %q to integer")
		}

		assign[name] = i
	}

	return assign, nil
}