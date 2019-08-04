package ec

import (
	"fmt"
	"go/types"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Component interface {
	private()
}

// TODO(mbm): Type and Function have similarities with corresponding go/types structs. Use them instead?

type Type struct {
	Name        string
	ElementType types.Type
	Coordinates []string
}

func (Type) private() {}

type Parameter struct {
	Name string
	Type Type
}

type Function struct {
	Name     string
	Receiver *Parameter
	Params   []*Parameter
	Results  []*Parameter
	Formula  *efd.Formula
}

func (Function) private() {}

// PrimitiveProgram lower the formula program to a primitive.
func (f Function) PrimitiveProgram() (*ast.Program, error) {
	original := f.Formula.Program
	if original == nil {
		return nil, xerrors.Errorf("function %s: missing op3 program", f.Name)
	}
	return op3.Lower(original)
}

// Outputs returns all output parameters. The receiver is considered an output, if present.
func (f Function) Outputs() []*Parameter {
	outputs := []*Parameter{}
	if f.Receiver != nil {
		outputs = append(outputs, f.Receiver)
	}
	outputs = append(outputs, f.Results...)
	return outputs
}

// Parameters returns all input and output parameters.
func (f Function) Parameters() []*Parameter {
	params := append([]*Parameter{}, f.Params...)
	return append(params, f.Outputs()...)
}

// Symbols returns the set of names defined by the function parameters.
func (f Function) Symbols() map[string]bool {
	symbols := map[string]bool{}
	for _, p := range f.Parameters() {
		symbols[p.Name] = true
	}
	return symbols
}

// Variables builds a map from program variable names to the Go code that
// references their corresponding function parameters.
func (f Function) Variables() map[ast.Variable]string {
	variables := map[ast.Variable]string{}
	n := 1
	for _, p := range f.Parameters() {
		for _, v := range p.Type.Coordinates {
			name := ast.Variable(fmt.Sprintf("%s%d", v, n))
			code := fmt.Sprintf("&%s.%s", p.Name, v)
			variables[name] = code
		}
		n++
	}
	return variables
}

type Config struct {
	PackageName string
	Field       fp.Config
	Components  []Component
}

func Package(cfg Config) (gen.Files, error) {
	fs := gen.Files{}

	// Point operations.
	b, err := Point(cfg)
	if err != nil {
		return nil, err
	}

	fs.Add("point.go", b)

	return fs, nil
}

type point struct {
	Config
	gocode.Generator
}

func Point(cfg Config) ([]byte, error) {
	p := &point{
		Config:    cfg,
		Generator: gocode.NewGenerator(),
	}
	return p.Generate()
}

func (p *point) Generate() ([]byte, error) {
	p.CodeGenerationWarning(gen.GeneratedBy)
	p.Package(p.PackageName)

	for _, component := range p.Components {
		switch c := component.(type) {
		case Type:
			p.typ(c)
		case Function:
			p.function(c)
		default:
			return nil, errutil.UnexpectedType(c)
		}
	}

	return p.Formatted()
}

func (p *point) typ(t Type) {
	p.Linef("type %s struct {", t.Name)
	for _, c := range t.Coordinates {
		p.Linef("\t%s %s", c, t.ElementType)
	}
	p.Linef("}")
}

func (p *point) function(f Function) {
	// Determine program.
	prog, err := f.PrimitiveProgram()
	if err != nil {
		p.SetError(err)
		return
	}

	// Function header.
	p.Printf("func ")
	if f.Receiver != nil {
		p.tuple([]*Parameter{f.Receiver})
	}
	p.Printf("%s", f.Name)
	p.tuple(f.Params)
	p.tuple(f.Results)
	p.EnterBlock()

	// Setup mapping from formula variables to code, and allocate any necessary
	// temporaries.
	variables := f.Variables()
	defined := f.Symbols()

	p.Linef("var (")
	for _, v := range op3.Variables(prog) {
		if _, ok := variables[v]; ok {
			continue
		}
		name := string(v)
		if _, conflict := defined[name]; conflict {
			name += "_"
		}
		p.Linef("%s %s", name, p.Field.Type())
		variables[v] = fmt.Sprintf("&%s", name)
	}
	p.Linef(")")

	// Generate program.
	p.NL()
	for _, a := range prog.Assignments {
		// TODO(mbm): ugly duplication in the switch statement below
		switch e := a.RHS.(type) {
		case ast.Pow:
			if e.N != 2 {
				p.SetError(errutil.AssertionFailure("power expected to be square"))
				return
			}
			p.Linef("Sqr(%s, %s)", variables[a.LHS], variables[e.X])
		case ast.Mul:
			vx, okx := e.X.(ast.Variable)
			vy, oky := e.Y.(ast.Variable)
			if !okx || !oky {
				p.SetError(errutil.AssertionFailure("operands should be variables"))
			}
			p.Linef("Mul(%s, %s, %s)", variables[a.LHS], variables[vx], variables[vy])
		case ast.Sub:
			vx, okx := e.X.(ast.Variable)
			vy, oky := e.Y.(ast.Variable)
			if !okx || !oky {
				p.SetError(errutil.AssertionFailure("operands should be variables"))
			}
			p.Linef("Sub(%s, %s, %s)", variables[a.LHS], variables[vx], variables[vy])
		case ast.Add:
			vx, okx := e.X.(ast.Variable)
			vy, oky := e.Y.(ast.Variable)
			if !okx || !oky {
				p.SetError(errutil.AssertionFailure("operands should be variables"))
			}
			p.Linef("Add(%s, %s, %s)", variables[a.LHS], variables[vx], variables[vy])
		default:
			p.SetError(errutil.UnexpectedType(e))
			return
		}
	}

	p.LeaveBlock()
}

func (p *point) tuple(params []*Parameter) {
	if len(params) == 0 {
		return
	}
	p.Printf("(")
	for i, param := range params {
		if i > 0 {
			p.Printf(", ")
		}
		p.Printf("%s *%s", param.Name, param.Type.Name)
	}
	p.Printf(")")
}
