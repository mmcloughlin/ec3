package efd

import "github.com/mmcloughlin/ec3/efd/op3/ast"

type Representation struct {
	Class string
	Shape *Shape

	Name       string
	Assume     []string
	Parameters []string
	Variables  []string
	Satisfying []string
}

type Formula struct {
	Class          string
	Shape          *Shape
	Representation *Representation
	Operation      string

	Source     string
	AppliesTo  string
	Assume     []string
	Compute    []string
	Parameters []string
	Program    *ast.Program
}

type Shape struct {
}
