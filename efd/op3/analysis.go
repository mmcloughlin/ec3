package op3

import (
	"math/bits"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Variables returns all variables used in the program.
func Variables(p *ast.Program) []ast.Variable {
	seen := map[ast.Variable]bool{}
	vs := []ast.Variable{}
	for _, a := range p.Assignments {
		for _, v := range ast.Variables(a.RHS.Inputs()) {
			if !seen[v] {
				vs = append(vs, v)
				seen[v] = true
			}
		}
	}
	return vs
}

// Inputs returns input variables for the given program.
func Inputs(p *ast.Program) []ast.Variable {
	// Inputs are variables that are read before they are written.
	inputs := map[ast.Variable]bool{}
	written := map[ast.Variable]bool{}
	for _, a := range p.Assignments {
		for _, v := range ast.Variables(a.RHS.Inputs()) {
			if !written[v] {
				inputs[v] = true
			}
		}
		written[a.LHS] = true
	}

	// Convert to slice.
	vs := make([]ast.Variable, 0, len(inputs))
	for input := range inputs {
		vs = append(vs, input)
	}

	return vs
}

// IsSSA reports whether every variable is written once.
func IsSSA(p *ast.Program) bool {
	seen := map[ast.Variable]bool{}
	for _, a := range p.Assignments {
		v := a.LHS
		if seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

// RenameVariables applies the given variable replacements to the program p.
func RenameVariables(p *ast.Program, replacements map[ast.Variable]ast.Variable) *ast.Program {
	r := &ast.Program{}
	for _, a := range p.Assignments {
		var expr ast.Expression
		switch e := a.RHS.(type) {
		case ast.Pow:
			expr = ast.Pow{
				X: renamevariable(e.X, replacements),
				N: e.N,
			}
		case ast.Inv:
			expr = ast.Inv{X: renameoperand(e.X, replacements)}
		case ast.Mul:
			expr = ast.Mul{
				X: renameoperand(e.X, replacements),
				Y: renameoperand(e.Y, replacements),
			}
		case ast.Neg:
			expr = ast.Neg{X: renameoperand(e.X, replacements)}
		case ast.Add:
			expr = ast.Add{
				X: renameoperand(e.X, replacements),
				Y: renameoperand(e.Y, replacements),
			}
		case ast.Sub:
			expr = ast.Sub{
				X: renameoperand(e.X, replacements),
				Y: renameoperand(e.Y, replacements),
			}
		case ast.Variable:
			expr = renamevariable(e, replacements)
		case ast.Constant:
			expr = e
		default:
			panic(errutil.UnexpectedType(e))
		}
		r.Assignments = append(r.Assignments, ast.Assignment{
			LHS: renamevariable(a.LHS, replacements),
			RHS: expr,
		})
	}
	return r
}

func renameoperand(op ast.Operand, replacements map[ast.Variable]ast.Variable) ast.Operand {
	v, ok := op.(ast.Variable)
	if !ok {
		return op
	}
	return renamevariable(v, replacements)
}

func renamevariable(v ast.Variable, replacements map[ast.Variable]ast.Variable) ast.Variable {
	if r, ok := replacements[v]; ok {
		return r
	}
	return v
}

// Pare down the given program to only the operations required to produce given
// outputs.
func Pare(p *ast.Program, outputs []ast.Variable) (*ast.Program, error) {
	// This is essentially liveness analysis for a single basic block.

	// Initially, the required outputs are live.
	live := map[ast.Variable]bool{}
	for _, output := range outputs {
		live[output] = true
	}

	// Process the program in reverse order.
	n := len(p.Assignments)
	required := make([]ast.Assignment, 0, n)
	for i := n - 1; i >= 0; i-- {
		a := p.Assignments[i]

		// If the variable written to is live, then this operation is required.
		if live[a.LHS] {
			required = append(required, a)
		}

		// Kill the variable that's written.
		delete(live, a.LHS)

		// Input variables are live.
		for _, v := range ast.Variables(a.RHS.Inputs()) {
			live[v] = true
		}
	}

	// Required assignments list was created in reverse order.
	for l, r := 0, len(required)-1; l < r; l, r = l+1, r-1 {
		required[l], required[r] = required[r], required[l]
	}

	return &ast.Program{
		Assignments: required,
	}, nil
}

// IsPrimitive reports whether p consists of only primitive operations.
func IsPrimitive(p *ast.Program) bool {
	for _, a := range p.Assignments {
		if !IsPrimitiveExpression(a.RHS) {
			return false
		}
	}
	return true
}

// IsPrimitiveExpression reports whether expr is a primitive operation. By
// "primitive", we mean that the operation directly translates to an efficient
// field operation.
//
// 1. Powers other than square are not considered primitive. They must be
// decomposed into multiplies and squares.
//
// 2. Constant multiplies are not primitive. Although they could be executed as
// multiplies, they are typically best replaced with additions.
func IsPrimitiveExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case ast.Inv, ast.Neg, ast.Add, ast.Sub, ast.Variable, ast.Constant:
		return true
	case ast.Pow:
		return e.N == 2
	case ast.Mul:
		_, xvar := e.X.(ast.Variable)
		_, yvar := e.Y.(ast.Variable)
		return xvar && yvar
	default:
		return false
	}
}

// Lower p to primitive expressions.
func Lower(p *ast.Program) (*ast.Program, error) {
	l := &ast.Program{}
	for _, a := range p.Assignments {
		as, err := lower(a)
		if err != nil {
			return nil, err
		}
		l.Assignments = append(l.Assignments, as...)
	}
	return l, nil
}

func lower(a ast.Assignment) ([]ast.Assignment, error) {
	if IsPrimitiveExpression(a.RHS) {
		return []ast.Assignment{a}, nil
	}

	// Non-primitive instructions are raising to a constant power x^c or
	// multiplying by a constant c*x. In either case lowering the instructions
	// reduces to an addition chain for the constant c. In this case the constants
	// we actually see are either small or powers of two; in particular, naive
	// addition chain algrithms will given us the optimal answer. Therefore we will
	// apply the left-to-right binary algorithm. This also has the benefit that
	// every expression in the result either just involves the previous result, or
	// the original input x, so we don't have to worry about allocating
	// temporaries.

	// Depending on power or multiply, the define the interpretation of "add one"
	// and "double" when building the constant.
	var c uint
	var first, add1, dbl ast.Expression
	switch e := a.RHS.(type) {
	case ast.Pow:
		c = uint(e.N)
		first = ast.Pow{X: e.X, N: 2}
		add1 = ast.Mul{X: a.LHS, Y: e.X}
		dbl = ast.Pow{X: a.LHS, N: 2}
	case ast.Mul:
		x, ok := e.X.(ast.Constant)
		if !ok {
			return nil, errutil.AssertionFailure("expected constant multiplier")
		}
		c = uint(x)
		first = ast.Add{X: e.Y, Y: e.Y}
		add1 = ast.Add{X: a.LHS, Y: e.Y}
		dbl = ast.Add{X: a.LHS, Y: a.LHS}
	default:
		return nil, errutil.UnexpectedType(e)
	}

	// Left-to-right binary algorithm on the constant c.
	//
	// Note that since the first step is special, we start the loop from the second
	// highest bit.
	exprs := []ast.Expression{first}
	for i := bits.Len(c) - 2; i >= 0; i-- {
		if (c & (1 << uint(i))) != 0 {
			exprs = append(exprs, add1)
		}
		if i > 0 {
			exprs = append(exprs, dbl)
		}
	}

	// Convert to assignments.
	as := []ast.Assignment{}
	for _, expr := range exprs {
		as = append(as, ast.Assignment{
			LHS: a.LHS,
			RHS: expr,
		})
	}

	return as, nil
}
