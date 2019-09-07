package op3

import (
	"math/bits"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

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
	case ast.Inv, ast.Neg, ast.Add, ast.Sub, ast.Cond, ast.Variable, ast.Constant:
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
