package op3

import (
	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

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

// CheckSSA verifies that every variable is written once.
func CheckSSA(p *ast.Program) error {
	seen := map[ast.Variable]bool{}
	for _, a := range p.Assignments {
		v := a.LHS
		if seen[v] {
			return xerrors.Errorf("variable %s written more than once", v)
		}
	}
	return nil
}
