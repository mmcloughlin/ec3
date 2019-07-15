package op3

import "github.com/mmcloughlin/ec3/efd/op3/ast"

// Variables returns all variables used in the program.
func Variables(p *ast.Program) []ast.Variable {
	// Build a set of all variables.
	set := map[ast.Variable]bool{}
	for _, a := range p.Assignments {
		for _, v := range ast.Variables(a.RHS.Inputs()) {
			set[v] = true
		}
	}

	// Convert to slice.
	vs := make([]ast.Variable, 0, len(set))
	for v := range set {
		vs = append(vs, v)
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
	case ast.Inv, ast.Neg, ast.Add, ast.Sub:
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
