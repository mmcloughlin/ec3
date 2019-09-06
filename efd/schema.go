package efd

import (
	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

var baseurl = "https://hyperelliptic.org/EFD"

type Shape struct {
	Collection string
	ID         string
	Tag        string
	Class      string

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
	Collection string
	ID         string
	Tag        string
	Class      string
	Shape      *Shape

	Name       string
	Assume     []string
	Parameters []string
	Variables  []string
	Satisfying []string
}

type Formula struct {
	Collection     string
	ID             string
	Tag            string
	Class          string
	Shape          *Shape
	Representation *Representation
	Operation      string
	URL            string

	Source     string
	AppliesTo  string
	Assume     []string
	Compute    []string
	Parameters []string
	Program    *ast.Program
}

// AllParameters returns all parameter variables for this formula. This merges
// the parameters from the shape, representation and formula itself.
func (f Formula) AllParameters() []string {
	params := []string{}
	params = append(params, f.Shape.Parameters...)
	params = append(params, f.Representation.Parameters...)
	params = append(params, f.Parameters...)
	return params
}
