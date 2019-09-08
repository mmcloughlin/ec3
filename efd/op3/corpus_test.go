package op3

import (
	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

// WithInputSet returns a predicate for matching formulae with the given input set.
func WithInputSet(vs ...ast.Variable) efd.Predicate {
	return func(f *efd.Formula) bool {
		p := f.Program
		if p == nil {
			return false
		}

		inputs := InputSet(p)
		if len(inputs) != len(vs) {
			return false
		}

		for _, v := range vs {
			if _, ok := inputs[v]; !ok {
				return false
			}
		}

		return true
	}
}

// WithOutputs returns a predicate for matching formulae containing the given variables.
func WithOutputs(vs ...ast.Variable) efd.Predicate {
	return func(f *efd.Formula) bool {
		p := f.Program
		if p == nil {
			return false
		}

		outputs := VariableSet(Variables(p))

		for _, v := range vs {
			if _, ok := outputs[v]; !ok {
				return false
			}
		}

		return true
	}
}

// Corpus returns a suite of test programs.
func Corpus() map[string]*ast.Program {
	corpus := map[string]*ast.Program{}

	// Include everything from the EFD.
	for _, f := range efd.Select(efd.WithProgram) {
		corpus[f.ID] = f.Program
	}

	// Conditionals do not appear in the EFD. Include a simple program that uses them.
	corpus["cond"] = &ast.Program{
		Assignments: []ast.Assignment{
			{
				LHS: ast.Variable("a"),
				RHS: ast.Cond{X: ast.Variable("b"), C: ast.Variable("c")},
			},
		},
	}

	return corpus
}
