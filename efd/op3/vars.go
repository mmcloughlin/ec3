package op3

import (
	"fmt"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

// VariableGenerator is a method of generating new variable names.
type VariableGenerator interface {
	// MarkUsed marks a variable as used, ensuring it will not be returned by any
	// future call to New().
	MarkUsed(...ast.Variable)

	// New generates a new unique variable name.
	New() ast.Variable
}

type variablegenerator struct {
	used map[ast.Variable]bool
	next func() ast.Variable
}

// NewVariableGeneratorFunc builds a VariableGenerator based on a function that
// produces a sequence of possible variable names.
func NewVariableGeneratorFunc(next func() ast.Variable) VariableGenerator {
	return &variablegenerator{
		used: make(map[ast.Variable]bool),
		next: next,
	}
}

func (g *variablegenerator) MarkUsed(vs ...ast.Variable) {
	for _, v := range vs {
		g.used[v] = true
	}
}

func (g *variablegenerator) New() ast.Variable {
	for {
		v := g.next()
		if !g.used[v] {
			g.MarkUsed(v)
			return v
		}
	}
}

// IndexedVariables generates variables using an increasing index and
func IndexedVariables(format string) VariableGenerator {
	i := 0
	return NewVariableGeneratorFunc(func() ast.Variable {
		v := fmt.Sprintf(format, i)
		i++
		return ast.Variable(v)
	})
}

// Temporaries generates temporary variables in a standard form.
var Temporaries = IndexedVariables("t%d")
