package ec

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

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

type Representation struct {
	Name        string
	ElementType types.Type
	Coordinates []string
}

func (r Representation) Type() types.Type {
	fields := []*types.Var{}
	for _, coord := range r.Coordinates {
		field := types.NewField(token.NoPos, nil, coord, r.ElementType, false)
		fields = append(fields, field)
	}
	s := types.NewStruct(fields, nil)
	name := types.NewTypeName(token.NoPos, nil, r.Name, nil)
	return types.NewNamed(name, s, nil)
}

func (Representation) private() {}

type Variable interface {
	Pointer() string
	Value() string
}

type value string

func (v value) Pointer() string { return "&" + string(v) }
func (v value) Value() string   { return string(v) }

type pointer string

func (p pointer) Pointer() string { return string(p) }
func (p pointer) Value() string   { return "*" + string(p) }

type Parameter interface {
	Name() string
	Type() types.Type
	Variables() map[ast.Variable]Variable
}

// basic is a parameter representing a single variable.
type basic struct {
	name string
	typ  types.Type
	v    Variable
}

// Condition returns a condition variable parameter.
func Condition(name string) Parameter {
	return basic{
		name: name,
		typ:  types.Typ[types.Uint],
		v:    value(name),
	}
}

func (b basic) Name() string { return b.name }

func (b basic) Type() types.Type { return b.typ }

func (b basic) Variables() map[ast.Variable]Variable {
	return map[ast.Variable]Variable{
		ast.Variable(b.name): b.v,
	}
}

type point struct {
	name     string
	repr     Representation
	indicies []int
}

func Point(name string, r Representation, indicies ...int) Parameter {
	return point{
		name:     name,
		repr:     r,
		indicies: indicies,
	}
}

func (p point) Name() string { return p.name }

func (p point) Type() types.Type {
	return types.NewPointer(p.repr.Type())
}

func (p point) Variables() map[ast.Variable]Variable {
	vars := map[ast.Variable]Variable{}
	for _, idx := range p.indicies {
		for _, coord := range p.repr.Coordinates {
			name := ast.Variable(fmt.Sprintf("%s%d", coord, idx))
			code := fmt.Sprintf("%s.%s", p.Name(), coord)
			vars[name] = value(code)
		}
	}
	return vars
}

type Function struct {
	Name     string
	Receiver Parameter
	Params   []Parameter
	Results  []Parameter
	Formula  *ast.Program
}

func (Function) private() {}

// Program returns the program to be implemented by this function.
func (f Function) Program() (*ast.Program, error) {
	// Restrict to variables used in this function.
	outputs := []ast.Variable{}
	for v := range f.Variables() {
		outputs = append(outputs, v)
	}

	p, err := op3.Pare(f.Formula, outputs)
	if err != nil {
		return nil, err
	}

	return op3.Lower(p)
}

// Parameters returns all parameters.
func (f Function) Parameters() []Parameter {
	params := []Parameter{}
	if f.Receiver != nil {
		params = append(params, f.Receiver)
	}
	params = append(params, f.Params...)
	params = append(params, f.Results...)
	return params
}

// Symbols returns the set of names defined by the function parameters.
func (f Function) Symbols() map[string]bool {
	symbols := map[string]bool{}
	for _, p := range f.Parameters() {
		symbols[p.Name()] = true
	}
	return symbols
}

// Variables builds a map from program variable names to the Go code that
// references their corresponding function parameters.
func (f Function) Variables() map[ast.Variable]Variable {
	variables := map[ast.Variable]Variable{}
	for _, p := range f.Parameters() {
		for name, v := range p.Variables() {
			variables[name] = v
		}
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
	b, err := PointOperations(cfg)
	if err != nil {
		return nil, err
	}

	fs.Add("point.go", b)

	return fs, nil
}

type pointops struct {
	Config
	gocode.Generator
}

func PointOperations(cfg Config) ([]byte, error) {
	p := &pointops{
		Config:    cfg,
		Generator: gocode.NewGenerator(),
	}
	return p.Generate()
}

func (p *pointops) Generate() ([]byte, error) {
	p.CodeGenerationWarning(gen.GeneratedBy)
	p.Package(p.PackageName)

	p.Import("math/big")

	for _, component := range p.Components {
		p.NL()
		switch c := component.(type) {
		case Representation:
			p.representation(c)
		case Function:
			p.function(c)
		default:
			return nil, errutil.UnexpectedType(c)
		}
	}

	return p.Formatted()
}

func (p *pointops) representation(r Representation) {
	// Type declaration.
	p.Linef("type %s struct {", r.Name)
	for _, c := range r.Coordinates {
		p.Linef("\t%s %s", c, r.ElementType)
	}
	p.Linef("}")

	// Constructor.
	p.NL()
	p.Printf("func New%s(%s *big.Int) *%s", r.Name, strings.Join(r.Coordinates, ", "), r.Name)
	p.EnterBlock()
	p.Linef("p := new(%s)", r.Name)
	for _, v := range r.Coordinates {
		p.Linef("p.%s.SetInt(%s)", v, v)
	}
	for _, v := range r.Coordinates {
		p.encode("&p." + v)
	}
	p.Linef("return p")
	p.LeaveBlock()

	// Set from another point.
	p.NL()
	p.Printf("func (p *%s) Set(q *%s)", r.Name, r.Name)
	p.EnterBlock()
	p.Linef("*p = *q")
	p.LeaveBlock()

	// Conversion to big.Ints.
	p.NL()
	p.Printf("func (p *%s) Coordinates() (%s *big.Int)", r.Name, strings.Join(r.Coordinates, ", "))
	p.EnterBlock()
	prefix := "p."
	if p.Field.Montgomery() {
		prefix = "d"
		p.Linef("var d%s %s", strings.Join(r.Coordinates, ", d"), p.Field.Type())
		for _, v := range r.Coordinates {
			p.Linef("Decode(&d%s, &p.%s)", v, v)
		}
	}
	for _, v := range r.Coordinates {
		p.Linef("%s = %s%s.Int()", v, prefix, v)
	}
	p.Linef("return")
	p.LeaveBlock()
}

func (p *pointops) function(f Function) {
	// Determine program.
	prog, err := f.Program()
	if err != nil {
		p.SetError(err)
		return
	}

	// Function header.
	p.Printf("func ")
	if f.Receiver != nil {
		p.tuple([]Parameter{f.Receiver})
	}
	p.Printf("%s", f.Name)
	p.tuple(f.Params)
	if len(f.Results) > 0 {
		p.tuple(f.Results)
	}
	p.EnterBlock()

	// Setup return variables.
	for _, r := range f.Results {
		if ptr, ok := r.Type().(*types.Pointer); ok {
			p.Linef("%s = new(%s)", r.Name(), ptr.Elem())
		}
	}

	// Setup mapping from formula variables to code, and allocate any necessary
	// temporaries.
	variables := f.Variables()
	defined := f.Symbols()
	tmps := []string{}

	for _, v := range op3.Variables(prog) {
		if _, ok := variables[v]; ok {
			continue
		}
		name := string(v)
		if _, conflict := defined[name]; conflict {
			name += "_"
		}
		tmps = append(tmps, name)
		variables[v] = value(name)
	}

	p.declare(tmps)

	// Generate program.
	for _, a := range prog.Assignments {
		switch e := a.RHS.(type) {
		case ast.Variable:
			p.Linef("%s = %s", variables[a.LHS].Value(), variables[e].Value())
		case ast.Constant:
			p.constant(variables[a.LHS], int(e))
		case ast.Pow:
			if e.N != 2 {
				p.SetError(errutil.AssertionFailure("power expected to be square"))
				return
			}
			p.call("Sqr", a.LHS, e, variables)
		case ast.Inv:
			p.call("Inv", a.LHS, e, variables)
		case ast.Mul:
			p.call("Mul", a.LHS, e, variables)
		case ast.Sub:
			p.call("Sub", a.LHS, e, variables)
		case ast.Add:
			p.call("Add", a.LHS, e, variables)
		case ast.Neg:
			p.call("Neg", a.LHS, e, variables)
		case ast.Cond:
			p.Linef("CMov(%s, %s, %s)", variables[a.LHS].Pointer(), variables[e.X].Pointer(), variables[e.C].Value())
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

func (p *pointops) declare(vars []string) {
	switch len(vars) {
	case 0:
		return
	case 1:
		p.Linef("var %s %s", vars[0], p.Field.Type())
	default:
		p.Linef("var (")
		for _, name := range vars {
			p.Linef("%s %s", name, p.Field.Type())
		}
		p.Linef(")")
		p.NL()
	}
}

func (p *pointops) call(fn, lhs ast.Variable, expr ast.Expression, vars map[ast.Variable]Variable) {
	p.Printf("%s(%s", fn, vars[lhs].Pointer())
	for _, operand := range expr.Inputs() {
		v, ok := operand.(ast.Variable)
		if !ok {
			p.SetError(errutil.AssertionFailure("operand must be variable"))
			return
		}
		p.Printf(", %s", vars[v].Pointer())
	}
	p.Linef(")")
}

func (p *pointops) tuple(params []Parameter) {
	args := []string{}
	for _, param := range params {
		args = append(args, fmt.Sprintf("%s %s", param.Name(), param.Type()))
	}
	p.Printf("(%s)", strings.Join(args, ", "))
}

func (p *pointops) constant(v Variable, x int) {
	p.Linef("%s.SetInt64(%d)", v.Value(), x)
	p.encode(v.Pointer())
}

func (p *pointops) encode(v string) {
	if p.Field.Montgomery() {
		p.Linef("Encode(%s, %s)", v, v)
	}
}
