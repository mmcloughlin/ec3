// +build ignore

package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

// Command-line flags.
var (
	typename = flag.String("type", "", "type name")
	input    = flag.String("input", "", "path to input file")
	output   = flag.String("output", "", "path to output file (default stdout)")
)

func main() {
	log.SetPrefix("wrap: ")
	log.SetFlags(0)

	flag.Parse()

	// Parse directives from input.
	if *input == "" {
		log.Fatal("no input file")
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := &parser{
		Leader: "//wrap",
	}

	ws, err := p.Reader(f)
	if err != nil {
		log.Fatal(err)
	}

	// Add common methods.
	ws = append(ws, common...)

	// Sanity checks.
	if errs := Lint(ws); len(errs) > 0 {
		for _, err := range errs {
			log.Print(err)
		}
		log.Fatal("failing due to lint errors")
	}

	// Setup output writer.
	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = f
	}

	// Generate.
	g := &generator{
		Generator:   gocode.NewGenerator(),
		wrappers:    ws,
		defaulttype: *typename,
	}

	b, err := g.Generate()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

type Identifier struct {
	Name     string
	Type     string
	Variadic bool
}

type Function struct {
	Identifier
	Parameters []Identifier
}

type Wrapper struct {
	Doc []string
	Go  Function
	C   Function
}

// common methods to be added to all types.
var common = []*Wrapper{
	// Equality.
	{
		Doc: []string{
			"Eq returns x == y.",
		},
		Go: Function{
			Identifier: Identifier{Name: "Eq", Type: "Bool"},
			Parameters: []Identifier{
				{Name: "x"},
				{Name: "y"},
			},
		},
		C: Function{
			Identifier: Identifier{Name: "Z3_mk_eq"},
			Parameters: []Identifier{
				{Name: "x"},
				{Name: "y"},
			},
		},
	},
	// Distinct.
	{
		Doc: []string{
			"Distinct returns a predicate representing whether all input parameters are distinct.",
		},
		Go: Function{
			Identifier: Identifier{Name: "Distinct", Type: "Bool"},
			Parameters: []Identifier{
				{Name: "x"},
				{Name: "y", Variadic: true},
			},
		},
		C: Function{
			Identifier: Identifier{Name: "Z3_mk_distinct"},
			Parameters: []Identifier{
				{Name: "y", Variadic: true},
			},
		},
	},
	// If-then-else.
	{
		Doc: []string{
			"ITE returns x if c else y.",
		},
		Go: Function{
			Identifier: Identifier{Name: "ITE"},
			Parameters: []Identifier{
				{Name: "x"},
				{Name: "c", Type: "*Bool"},
				{Name: "y"},
			},
		},
		C: Function{
			Identifier: Identifier{Name: "Z3_mk_ite"},
			Parameters: []Identifier{
				{Name: "c"},
				{Name: "x"},
				{Name: "y"},
			},
		},
	},
}

// Lint runs some sanity checks on a set of wrappers.
func Lint(ws []*Wrapper) errutil.Errors {
	var errs errutil.Errors

	// Check for multiple wrappers of the same function.
	funcs := map[string]bool{}
	for _, w := range ws {
		name := w.C.Name
		if funcs[name] {
			errs.Addf("function %q is wrapped multiple times", name)
		}
		funcs[name] = true
	}

	// Check documentation.
	blacklist := []string{"t1", "t2"}
	for _, w := range ws {
		if len(w.Doc) == 0 {
			errs.Addf("function %q is undocumented", w.Go.Name)
			continue
		}
		if !strings.HasPrefix(w.Doc[0], w.Go.Name+" ") {
			errs.Addf("function %q documentation should start with %q", w.Go.Name, w.Go.Name)
		}
		for _, word := range blacklist {
			for _, line := range w.Doc {
				if strings.Contains(line, word) {
					errs.Addf("function %q documentation contains blacklisted word %q", w.Go.Name, word)
					break
				}
			}
		}
	}

	// Expect correspondance with the Z3 function names.
	abbrev := strings.NewReplacer(
		"Reduce", "red",
		"Logic", "l",
		"Arith", "a",
	)
	for _, w := range ws {
		goname := strings.ToLower(abbrev.Replace(w.Go.Name))
		cname := strings.ToLower(strings.ReplaceAll(w.C.Name, "_", ""))
		if !strings.HasSuffix(cname, goname) {
			errs.Addf("function name %q expected to match suffix of %q", w.Go.Name, w.C.Name)
		}
	}

	return errs
}

// parser parses directives from source code.
type parser struct {
	Leader string

	wrappers []*Wrapper
	cur      *Wrapper
}

// Reader parses wrappers from source code.
func (p *parser) Reader(r io.Reader) ([]*Wrapper, error) {
	p.reset()
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		// Non-directive line.
		if !strings.HasPrefix(line, p.Leader) {
			p.finish()
			continue
		}

		// Directive.
		if p.cur == nil {
			p.cur = &Wrapper{}
		}

		name, fields, err := p.directive(line)
		if err != nil {
			return nil, err
		}

		switch name {
		case "doc":
			err = p.doc(fields)
		case "unary":
			err = p.unary(fields)
		case "binary":
			err = p.binary(fields)
		case "go":
			err = p.golang(fields)
		case "c":
			err = p.c(fields)
		default:
			return nil, xerrors.Errorf("unknown directive type %q", name)
		}

		if err != nil {
			return nil, err
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	p.finish()

	return p.wrappers, nil
}

func (p *parser) directive(line string) (string, []string, error) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return "", nil, errutil.AssertionFailure("empty directive line")
	}

	parts := strings.Split(fields[0], ":")
	if parts[0] != p.Leader {
		return "", nil, errutil.AssertionFailure("expected directive leader")
	}

	return parts[1], fields[1:], nil
}

func (p *parser) reset() {
	p.wrappers = nil
	p.cur = nil
}

func (p *parser) finish() {
	if p.cur == nil {
		return
	}
	p.wrappers = append(p.wrappers, p.cur)
	p.cur = nil
}

func (p *parser) doc(fields []string) error {
	line := strings.Join(fields, " ")
	p.cur.Doc = append(p.cur.Doc, line)
	return nil
}

func (p *parser) unary(fields []string) error {
	if len(fields) != 2 {
		return xerrors.New("unary directive expects 2 arguments")
	}
	p.golang([]string{fields[0], "x"})
	p.c([]string{fields[1], "x"})
	return nil
}

func (p *parser) binary(fields []string) error {
	if len(fields) != 2 {
		return xerrors.New("binary directive expects 2 arguments")
	}
	p.golang([]string{fields[0], "x", "y"})
	p.c([]string{fields[1], "x", "y"})
	return nil
}

func (p *parser) golang(fields []string) error {
	var err error
	p.cur.Go, err = function(fields)
	return err
}

func (p *parser) c(fields []string) error {
	var err error
	p.cur.C, err = function(fields)
	return err
}

func function(fields []string) (Function, error) {
	ids := identifiers(fields)
	if len(ids) < 1 {
		return Function{}, xerrors.New("expected at least one identifier")
	}
	return Function{
		Identifier: ids[0],
		Parameters: ids[1:],
	}, nil
}

func identifiers(fields []string) []Identifier {
	ids := []Identifier{}
	for _, field := range fields {
		ids = append(ids, identifier(field))
	}
	return ids
}

func identifier(field string) Identifier {
	parts := strings.Split(field, ":")
	id := Identifier{
		Name: parts[0],
	}
	if len(parts) > 1 {
		id.Type = parts[1]
	}
	if strings.HasSuffix(id.Name, "...") {
		id.Variadic = true
		id.Name = strings.TrimSuffix(id.Name, "...")
	}
	return id
}

type generator struct {
	gocode.Generator

	wrappers    []*Wrapper
	defaulttype string
}

func (g *generator) Generate() ([]byte, error) {
	g.CodeGenerationWarningSelf()
	g.Package("z3")

	g.Printf(`
/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"
	`)

	for _, w := range g.wrappers {
		g.NL()
		g.wrapper(w)
	}

	return g.Formatted()
}

func (g *generator) wrapper(w *Wrapper) {
	// Documentation.
	g.Comment(w.Doc...)
	g.Commentf("Corresponds to %s.", w.C.Name)

	// Go signature.
	receiver := w.Go.Parameters[0]
	params := w.Go.Parameters[1:]
	returns := withdefault(w.Go.Type, g.defaulttype)

	g.Printf("func (")
	g.param(receiver)
	g.Printf(") %s (", w.Go.Name)
	for i, param := range params {
		if i > 0 {
			g.Printf(", ")
		}
		g.param(param)
	}
	g.Printf(") *%s", returns)

	// Function body.
	g.EnterBlock()

	for _, param := range params {
		if param.Variadic {
			array := variadicname(param.Name)
			g.Linef("%s := []C.Z3_ast{%s.ast}", array, receiver.Name)
			g.Linef("for _, a := range %s {", param.Name)
			g.Linef("%s = append(%s, a.ast)", array, array)
			g.Linef("}")
		}
	}

	g.Linef("return &%s{", returns)
	g.Linef("value: value{")
	g.Linef("ctx: %s.ctx,", receiver.Name)
	g.Printf("ast: C.%s(%s.ctx", w.C.Name, receiver.Name)
	for _, arg := range w.C.Parameters {
		g.Printf(", ")
		g.arg(arg)
	}
	g.Linef("),")

	g.Linef("},")
	g.Linef("}")
	g.LeaveBlock()
}

func (g *generator) param(id Identifier) {
	g.Printf("%s %s%s", id.Name, ternary(id.Variadic, "...", ""), g.typ(id.Type))
}

func (g *generator) arg(id Identifier) {
	switch {
	case id.Variadic:
		array := variadicname(id.Name)
		g.Printf("C.unsigned(len(%s)), &%s[0]", array, array)
	case id.Type == "":
		g.Printf("%s.ast", id.Name)
	default:
		g.Printf("C.%s(%s)", id.Type, id.Name)
	}
}

func (g *generator) typ(name string) string {
	return withdefault(name, "*"+g.defaulttype)
}

func variadicname(name string) string {
	return name + "s"
}

func withdefault(a, b string) string {
	return ternary(a != "", a, b)
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
