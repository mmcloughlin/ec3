package efd

import (
	"fmt"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

var baseurl = "https://hyperelliptic.org/EFD"

type Shape struct {
	ID    string
	Tag   string
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

func (s Shape) URL() string {
	// Example: https://hyperelliptic.org/EFD/g12o/auto-hessian.html
	return fmt.Sprintf("%s/%s/auto-%s.html", baseurl, s.Class, s.Tag)
}

type Representation struct {
	ID    string
	Tag   string
	Class string
	Shape *Shape

	Name       string
	Assume     []string
	Parameters []string
	Variables  []string
	Satisfying []string
}

func (r Representation) URL() string {
	// Example: https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html
	return fmt.Sprintf("%s/%s/auto-%s-%s.html", baseurl, r.Class, r.Shape.Tag, r.Tag)
}

type Formula struct {
	ID             string
	Tag            string
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

func (f Formula) URL() string {
	// Example: https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
	return fmt.Sprintf("%s#%s-%s", f.Representation.URL(), f.Operation, f.Tag)
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
