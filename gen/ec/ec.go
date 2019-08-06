package ec

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
	"github.com/mmcloughlin/ec3/internal/ints"
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
	Name          string
	Receiver      *Parameter
	WriteReceiver bool
	Params        []*Parameter
	Results       []*Parameter
	Formula       *efd.Formula
}

func (Function) private() {}

// Program returns the program to be implemented by this function.
func (f Function) Program() (*ast.Program, error) {
	// Fetch from the formula.
	original := f.Formula.Program
	if original == nil {
		return nil, xerrors.Errorf("function %s: missing op3 program", f.Name)
	}

	// Restrict to variables used in this function.
	outputs := []ast.Variable{}
	log.Print(f.Formula.Tag, f.Variables())
	for v := range f.Variables() {
		outputs = append(outputs, v)
	}

	p, err := op3.Pare(original, outputs)
	if err != nil {
		return nil, err
	}

	return op3.Lower(p)
}

// Inputs returns all input parameters.
func (f Function) Inputs() []*Parameter {
	inputs := []*Parameter{}
	if f.Receiver != nil && !f.WriteReceiver {
		inputs = append(inputs, f.Receiver)
	}
	inputs = append(inputs, f.Params...)
	return inputs
}

// Outputs returns all output parameters. The receiver is considered an output, if present.
func (f Function) Outputs() []*Parameter {
	outputs := []*Parameter{}
	if f.Receiver != nil && f.WriteReceiver {
		outputs = append(outputs, f.Receiver)
	}
	outputs = append(outputs, f.Results...)
	return outputs
}

// Parameters returns all input and output parameters.
func (f Function) Parameters() []*Parameter {
	return append(f.Inputs(), f.Outputs()...)
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
	// Assign indicies.
	byindex := map[int]*Parameter{}
	n := 1
	for _, p := range f.Inputs() {
		byindex[n] = p
		n++
	}

	n = ints.Max(n, 3)
	for _, p := range f.Outputs() {
		byindex[n] = p
		n++
	}

	// Create variable map.
	variables := map[ast.Variable]string{}
	for n, p := range byindex {
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

	p.Import("math/big")

	for _, component := range p.Components {
		p.NL()
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
	// Type declaration.
	p.Linef("type %s struct {", t.Name)
	for _, c := range t.Coordinates {
		p.Linef("\t%s %s", c, t.ElementType)
	}
	p.Linef("}")

	// Constructor.
	p.NL()
	p.Printf("func New%s(%s *big.Int) *%s", t.Name, strings.Join(t.Coordinates, ", "), t.Name)
	p.EnterBlock()
	p.Linef("p := new(%s)", t.Name)
	for _, v := range t.Coordinates {
		p.Linef("p.%s.SetInt(%s)", v, v)
	}
	p.Linef("return p")
	p.LeaveBlock()

	// Conversion to big.Ints.
	p.NL()
	p.Printf("func (p *%s) Coordinates() (%s *big.Int)", t.Name, strings.Join(t.Coordinates, ", "))
	p.EnterBlock()
	for _, v := range t.Coordinates {
		p.Linef("%s = p.%s.Int()", v, v)
	}
	p.Linef("return")
	p.LeaveBlock()
}

func (p *point) function(f Function) {
	// Determine program.
	prog, err := f.Program()
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
	if len(f.Results) > 0 {
		p.tuple(f.Results)
	}
	p.EnterBlock()

	// Setup return variables.
	for _, r := range f.Results {
		p.Linef("%s = new(%s)", r.Name, r.Type.Name)
	}

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
		case ast.Inv:
			x, ok := e.X.(ast.Variable)
			if !ok {
				p.SetError(errutil.AssertionFailure("operand should be variable"))
			}
			p.Linef("Inv(%s, %s)", variables[a.LHS], variables[x])
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

	if len(f.Results) > 0 {
		p.Linef("return")
	}
	p.LeaveBlock()
}

func (p *point) tuple(params []*Parameter) {
	p.Printf("(")
	for i, param := range params {
		if i > 0 {
			p.Printf(", ")
		}
		p.Printf("%s *%s", param.Name, param.Type.Name)
	}
	p.Printf(")")
}
