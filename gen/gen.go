package gen

import "go/format"

// GeneratedBy is the name used in code generation warnings.
const GeneratedBy = "ec3"

// Interface for a code generator.
type Interface interface {
	Generate() ([]byte, error)
}

// Func allows a standalone function to implement Interface.
type Func func() ([]byte, error)

// Generate calls f.
func (f Func) Generate() ([]byte, error) {
	return f()
}

// Formatted builds a generator that runs Go formatting on the result of g.
func Formatted(g Interface) Interface {
	return Func(func() ([]byte, error) {
		b, err := g.Generate()
		if err != nil {
			return nil, err
		}
		return format.Source(b)
	})
}
