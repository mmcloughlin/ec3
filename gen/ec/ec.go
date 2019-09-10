package ec

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/internal/container/disjointset"
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

func (Representation) private() {}

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

func (r Representation) Equals(other Representation) bool {
	return reflect.DeepEqual(r, other)
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
	// Restrict to variables output by this function.
	outputs := []ast.Variable{}
	for _, out := range f.Outputs() {
		for v := range out.Variables() {
			outputs = append(outputs, v)
		}
	}

	// Reduce formula given required output variables.
	p, err := op3.Pare(f.Formula, outputs)
	if err != nil {
		return nil, err
	}

	// Ensure the program is robust to potential alias sets.
	aliases := f.AliasSets()
	q := op3.AliasCorrect(p, aliases, outputs, op3.Temporaries())

	// Finally, reduce the program to primitives.
	return op3.Lower(q)
}

// AliasSets returns groups of variable names with a may-alias relationship,
// meaning there is a possibility they are pointers to the same memory
// locations.
func (f Function) AliasSets() [][]ast.Variable {
	// Build sets of aliases using a disjoint-set structure.
	d := disjointset.New()
	params := f.Parameters()
	n := len(params)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			sets := params[i].AliasSets(params[j])
			for _, set := range sets {
				for _, v := range set {
					d.Union(string(set[0]), string(v))
				}
			}
		}
	}

	// Recover sets from the disjoint-set structure.
	sets := map[string][]ast.Variable{}
	for name := range f.Variables() {
		s := d.Find(string(name))
		sets[s] = append(sets[s], name)
	}

	// Transform into alias sets.
	aliases := [][]ast.Variable{}
	for _, set := range sets {
		if len(set) > 1 {
			aliases = append(aliases, set)
		}
	}

	return aliases
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

// ParametersWithAction returns all parameters supporting action a.
func (f Function) ParametersWithAction(a Action) []Parameter {
	params := []Parameter{}
	for _, param := range f.Parameters() {
		if param.Action().Contains(a) {
			params = append(params, param)
		}
	}
	return params
}

// Inputs returns all read parameters.
func (f Function) Inputs() []Parameter {
	return f.ParametersWithAction(R)
}

// Outputs returns all write parameters.
func (f Function) Outputs() []Parameter {
	return f.ParametersWithAction(W)
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

	// Verifiy that all inputs have corresponding variables in the function.
	variables := f.Variables()
	for _, input := range op3.Inputs(prog) {
		if _, ok := variables[input]; !ok {
			p.SetError(xerrors.Errorf("no variable defined for program input %s", input))
			return
		}
	}

	// Setup mapping from formula variables to code, and allocate any necessary
	// temporaries.
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

	// Output sorted for reproducability.
	sort.Strings(tmps)
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
