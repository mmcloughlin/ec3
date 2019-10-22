package name

import (
	"fmt"
)

// VariableGenerator is a method of generating new variable names.
type VariableGenerator interface {
	// MarkUsed marks a variable as used, ensuring it will not be returned by any
	// future call to New().
	MarkUsed(...string)

	// New generates a new unique variable name.
	New() string
}

type variablegenerator struct {
	used map[string]bool
	next func() string
}

// NewVariableGeneratorFunc builds a VariableGenerator based on a function that
// produces a sequence of possible variable names.
func NewVariableGeneratorFunc(next func() string) VariableGenerator {
	return &variablegenerator{
		used: make(map[string]bool),
		next: next,
	}
}

func (g *variablegenerator) MarkUsed(vs ...string) {
	for _, v := range vs {
		g.used[v] = true
	}
}

func (g *variablegenerator) New() string {
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
	return NewVariableGeneratorFunc(func() string {
		v := fmt.Sprintf(format, i)
		i++
		return string(v)
	})
}

// Temporaries generates temporary variables in a standard form.
func Temporaries() VariableGenerator {
	return IndexedVariables("t%d")
}
