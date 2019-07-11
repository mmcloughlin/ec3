package efd

import "github.com/mmcloughlin/ec3/efd/op3/ast"

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

type Representation struct {
}
