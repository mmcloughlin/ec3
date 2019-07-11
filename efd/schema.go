package efd

import "github.com/mmcloughlin/ec3/efd/op3/ast"

type Shape struct {
	ID    string
	Class string

	Name            string
	Parameters      []string
	Coordinates     []string
	A               [7]string
	Satisfying      []string
	Addition        []string
	Doubling        []string
	Negation        []string
	Neutral         []string
	FromWeierstrass []string
	ToWeierstrass   []string
}

type Representation struct {
	ID    string
	Class string
	Shape *Shape

	Name       string
	Assume     []string
	Parameters []string
	Variables  []string
	Satisfying []string
}

type Formula struct {
	ID             string
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
