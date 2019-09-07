package op3

import (
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
		case ast.Cond:
			expr = ast.Cond{
				X: renamevariable(e.X, replacements),
				C: renamevariable(e.C, replacements),
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

// LiveSet maintains a set of live variables.
type LiveSet map[ast.Variable]bool

// NewLiveSet constructs an empty set of live variables.
func NewLiveSet() LiveSet {
	return make(LiveSet)
}

// MarkLive records all variables in vs as live.
func (l LiveSet) MarkLive(vs ...ast.Variable) {
	for _, v := range vs {
		l[v] = true
	}
}

// Update the live set based on the assignment. In liveness analysis, program
// assignments should be processed in reverse.
func (l LiveSet) Update(a ast.Assignment) {
	// Kill the variable that's written.
	delete(l, a.LHS)

	// Input variables are live.
	inputs := ast.Variables(a.RHS.Inputs())
	l.MarkLive(inputs...)
}

// Pare down the given program to only the operations required to produce given
// outputs.
func Pare(p *ast.Program, outputs []ast.Variable) (*ast.Program, error) {
	// This is essentially liveness analysis for a single basic block.

	// Initially, the required outputs are live.
	live := NewLiveSet()
	live.MarkLive(outputs...)

	// Process the program in reverse order.
	n := len(p.Assignments)
	required := make([]ast.Assignment, 0, n)
	for i := n - 1; i >= 0; i-- {
		a := p.Assignments[i]

		// If the variable written to is live, then this operation is required.
		if live[a.LHS] {
			required = append(required, a)
		}

		// Update liveness.
		live.Update(a)
	}

	// Required assignments list was created in reverse order.
	for l, r := 0, len(required)-1; l < r; l, r = l+1, r-1 {
		required[l], required[r] = required[r], required[l]
	}

	return &ast.Program{
		Assignments: required,
	}, nil
}
