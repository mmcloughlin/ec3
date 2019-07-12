package main

import (
	"flag"
	"log"
	"os"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/print"
)

var (
	class = flag.String("class", "g1p", "class of curve")
	shape = flag.String("shape", "shortw", "curve shape")
	repr  = flag.String("repr", "jacobian-3", "representation")
	op    = flag.String("op", "addition", "operation")
)

func main() {
	flag.Parse()

	// Prepare filters.
	predicates := []efd.Predicate{}
	if *class != "" {
		predicates = append(predicates, efd.WithClass(*class))
	}
	if *shape != "" {
		predicates = append(predicates, efd.WithShape(*shape))
	}
	if *repr != "" {
		predicates = append(predicates, efd.WithRepresentation(*repr))
	}
	if *op != "" {
		predicates = append(predicates, efd.WithOperation(*op))
	}

	// Get and print list of selected formulae.
	fs := efd.Formulae(predicates...)

	p := &printer{
		TabWriter: print.NewTabWriter(os.Stdout, 1, 4, 4, ' ', 0),
	}
	p.formulae(fs)
	p.Flush()

	if err := p.Error(); err != nil {
		log.Fatal(err)
	}
}

type printer struct {
	*print.TabWriter
}

func (p *printer) formulae(fs []*efd.Formula) {
	for _, f := range fs {
		p.formula(f)
	}
}

func (p *printer) formula(f *efd.Formula) {
	p.field("id", f.ID)
	p.field("tag", f.Tag)
	p.field("class", f.Class)
	p.field("shape", f.Shape.Tag)
	p.field("repr", f.Representation.Tag)
	p.field("operation", f.Operation)

	p.maybe("source", f.Source)
	p.maybe("appliesto", f.AppliesTo)
	p.values("params", f.Parameters)
	p.values("assume", f.Assume)
	p.values("compute", f.Compute)
	p.program("op3", f.Program)
}

func (p *printer) field(key, value string) {
	p.Linef("%s\t%s", key, value)
}

func (p *printer) maybe(key, value string) {
	if len(value) > 0 {
		p.field(key, value)
	}
}

func (p *printer) values(key string, values []string) {
	if len(values) == 0 {
		return
	}
	p.field(key, values[0])
	for _, value := range values[1:] {
		p.field("", value)
	}
}

func (p *printer) program(key string, prog *ast.Program) {
	if prog == nil {
		return
	}
	lines := []string{}
	for _, a := range prog.Assignments {
		lines = append(lines, a.String())
	}
	p.values(key, lines)
}
