package ec

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Function struct {
	Name       string
	Parameters []string
	Formula    *efd.Formula
}

type Config struct {
	Field          fp.Config
	Representation *efd.Representation
	Functions      []Function

	PackageName string
	TypeName    string
}

func (c Config) Receiver() string { return "p" }

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

	// Type definition.
	p.Linef("type %s struct {", p.TypeName)
	for _, v := range p.Representation.Variables {
		p.Linef("\t%s %s", v, p.Field.Type())
	}
	p.Linef("}")

	// Formulae.
	for _, fn := range p.Functions {
		p.function(fn)
	}

	return p.Formatted()
}

func (p *point) function(fn Function) {
	// Lower the formula program to a primitive.
	original := fn.Formula.Program
	if original == nil {
		p.SetError(xerrors.Errorf("function %s: missing op3 program", fn.Name))
		return
	}

	prog, err := op3.Lower(original)
	if err != nil {
		p.SetError(err)
		return
	}

	// Setup input and outputs.
	points := append([]string{}, fn.Parameters...)
	points = append(points, p.Receiver())

	variables := map[ast.Variable]string{}
	n := 1
	for _, point := range points {
		for _, v := range p.Representation.Variables {
			name := ast.Variable(fmt.Sprintf("%s%d", v, n))
			code := fmt.Sprintf("&%s.%s", point, v)
			variables[name] = code
		}
		n++
	}

	// Function header.
	p.Printf("func (%s *%s) %s(%s *%s)", p.Receiver(), p.TypeName, fn.Name, strings.Join(fn.Parameters, ", "), p.TypeName)
	p.EnterBlock()

	// Allocate temporaries.
	p.Linef("var (")
	for _, v := range op3.Variables(prog) {
		if _, ok := variables[v]; ok {
			continue
		}
		p.Linef("%s %s", v, p.Field.Type())
		variables[v] = fmt.Sprintf("&%s", v)
	}
	p.Linef(")")

	// Generate program.
	p.NL()
	for _, a := range prog.Assignments {
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
