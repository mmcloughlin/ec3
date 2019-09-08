// Package eval provides an evaluator for op3 programs.
package eval

import (
	"math/big"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

type Evaluator struct {
	state map[ast.Variable]*big.Int
	m     *big.Int
}

// NewEvaluator builds an evaluator using arithmetic modulo m.
func NewEvaluator(m *big.Int) *Evaluator {
	return &Evaluator{
		state: make(map[ast.Variable]*big.Int),
		m:     m,
	}
}

// Load the variable v.
func (e *Evaluator) Load(v ast.Variable) (*big.Int, bool) {
	x, ok := e.state[v]
	return x, ok
}

// Store x into the variable v.
func (e *Evaluator) Store(v ast.Variable, x *big.Int) {
	e.state[v] = x
}

// Initialize the variable v to x. Errors if v is already defined.
func (e *Evaluator) Initialize(v ast.Variable, x *big.Int) error {
	if _, ok := e.Load(v); ok {
		return xerrors.Errorf("variable %q is already defined", v)
	}
	e.Store(v, x)
	return nil
}

// Execute the program p.
func (e *Evaluator) Execute(p *ast.Program) error {
	for _, a := range p.Assignments {
		if err := e.assignment(a); err != nil {
			return err
		}
	}
	return nil
}

func (e *Evaluator) assignment(a ast.Assignment) error {
	lhs := e.dst(a.LHS)
	switch expr := a.RHS.(type) {
	case ast.Pow:
		x, err := e.operands(expr.X, expr.N)
		if err != nil {
			return err
		}
		lhs.Exp(x[0], x[1], e.m)
	case ast.Inv:
		x, err := e.operand(expr.X)
		if err != nil {
			return err
		}
		lhs.ModInverse(x, e.m)
	case ast.Mul:
		x, err := e.operands(expr.X, expr.Y)
		if err != nil {
			return err
		}
		lhs.Mul(x[0], x[1])
	case ast.Neg:
		x, err := e.operand(expr.X)
		if err != nil {
			return err
		}
		lhs.Neg(x)
	case ast.Add:
		x, err := e.operands(expr.X, expr.Y)
		if err != nil {
			return err
		}
		lhs.Add(x[0], x[1])
	case ast.Sub:
		x, err := e.operands(expr.X, expr.Y)
		if err != nil {
			return err
		}
		lhs.Sub(x[0], x[1])
	case ast.Cond:
		x, err := e.operands(expr.X, expr.C)
		if err != nil {
			return err
		}
		if x[1].Sign() != 0 {
			lhs.Set(x[0])
		}
	case ast.Variable, ast.Constant:
		x, err := e.operand(expr)
		if err != nil {
			return err
		}
		lhs.Set(x)
	default:
		return errutil.UnexpectedType(expr)
	}
	lhs.Mod(lhs, e.m)
	return nil
}

func (e Evaluator) dst(v ast.Variable) *big.Int {
	if x, ok := e.Load(v); ok {
		return x
	}
	x := new(big.Int)
	e.Store(v, x)
	return x
}

func (e *Evaluator) operands(operands ...ast.Operand) ([]*big.Int, error) {
	xs := make([]*big.Int, 0, len(operands))
	for _, operand := range operands {
		x, err := e.operand(operand)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, nil
}

func (e *Evaluator) operand(operand ast.Operand) (*big.Int, error) {
	switch op := operand.(type) {
	case ast.Variable:
		x, ok := e.Load(op)
		if !ok {
			return nil, xerrors.Errorf("variable %q is not defined", op)
		}
		return x, nil
	case ast.Constant:
		return new(big.Int).SetUint64(uint64(op)), nil
	default:
		return nil, errutil.UnexpectedType(op)
	}
}
