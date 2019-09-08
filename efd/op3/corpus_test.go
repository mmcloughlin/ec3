package op3

import (
	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

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
