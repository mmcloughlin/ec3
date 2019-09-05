// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/db"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

var (
	archive = flag.String("archive", "efd.tar.gz", "path to efd archive")
	output  = flag.String("output", "", "path to output file (default stdout)")
)

func main() {
	flag.Parse()

	// Read archive.
	a, err := db.Archive(*archive)
	if err != nil {
		log.Fatal(err)
	}

	d, err := db.Read(a)
	if err != nil {
		log.Fatal(err)
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
		Database:  d,
		Generator: gocode.NewGenerator(),
	}
	b, err := g.Generate()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

type generator struct {
	*db.Database
	gocode.Generator

	shapeidx map[string]int
	repridx  map[string]int
}

func (g *generator) Generate() ([]byte, error) {
	g.CodeGenerationWarningSelf()
	g.Package("efd")
	g.Import("github.com/mmcloughlin/ec3/efd/op3/ast")

	g.shapes()
	g.representations()
	g.formulae()

	return g.Formatted()
}

func (g *generator) shapes() {
	// Sort them to ensure reproducability.
	ss := make([]*efd.Shape, 0, len(g.Shapes))
	for _, s := range g.Shapes {
		ss = append(ss, s)
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].ID < ss[j].ID
	})

	// Output code.
	g.shapeidx = map[string]int{}
	g.NL()
	g.Printf("var shapes = []*Shape")
	g.EnterBlock()
	for i, s := range ss {
		g.shape(s)
		g.shapeidx[s.ID] = i
	}
	g.LeaveBlock()
}

func (g *generator) shape(s *efd.Shape) {
	g.EnterBlock()
	g.Linef("ID: %#v,", s.ID)
	g.Linef("Tag: %#v,", s.Tag)
	g.Linef("Class: %#v,", s.Class)
	g.Linef("Name: %#v,", s.Name)
	g.Linef("Parameters: %#v,", s.Parameters)
	g.Linef("Coordinates: %#v,", s.Coordinates)
	g.Linef("A: %#v,", s.A)
	g.Linef("Satisfying: %#v,", s.Satisfying)
	g.Linef("Addition: %#v,", s.Addition)
	g.Linef("Doubling: %#v,", s.Doubling)
	g.Linef("Negation: %#v,", s.Negation)
	g.Linef("Neutral: %#v,", s.Neutral)
	g.Linef("FromWeierstrass: %#v,", s.FromWeierstrass)
	g.Linef("ToWeierstrass: %#v,", s.ToWeierstrass)
	g.Dedent()
	g.Linef("},")
}

func (g *generator) shaperef(s *efd.Shape) string {
	idx, ok := g.shapeidx[s.ID]
	if !ok {
		g.SetError(xerrors.Errorf("unknown shape %q", s.ID))
	}
	return fmt.Sprintf("shapes[%d]", idx)
}

func (g *generator) representations() {
	// Sort them to ensure reproducability.
	rs := make([]*efd.Representation, 0, len(g.Representations))
	for _, r := range g.Representations {
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].ID < rs[j].ID
	})

	// Output.
	g.repridx = map[string]int{}
	g.NL()
	g.Printf("var representations = []*Representation")
	g.EnterBlock()
	for i, r := range rs {
		g.representation(r)
		g.repridx[r.ID] = i
	}
	g.LeaveBlock()
}

func (g *generator) representation(r *efd.Representation) {
	g.EnterBlock()
	g.Linef("ID: %#v,", r.ID)
	g.Linef("Tag: %#v,", r.Tag)
	g.Linef("Class: %#v,", r.Class)
	g.Linef("Shape: %s,", g.shaperef(r.Shape))
	g.Linef("Name: %#v,", r.Name)
	g.Linef("Assume: %#v,", r.Assume)
	g.Linef("Parameters: %#v,", r.Parameters)
	g.Linef("Variables: %#v,", r.Variables)
	g.Linef("Satisfying: %#v,", r.Satisfying)
	g.Dedent()
	g.Linef("},")
}

func (g *generator) representationref(r *efd.Representation) string {
	idx, ok := g.repridx[r.ID]
	if !ok {
		g.SetError(xerrors.Errorf("unknown representation %q", r.ID))
	}
	return fmt.Sprintf("representations[%d]", idx)
}

func (g *generator) formulae() {
	// Sort them to ensure reproducability.
	fs := make([]*efd.Formula, 0, len(g.Formulae))
	for _, f := range g.Formulae {
		fs = append(fs, f)
	}
	sort.Slice(fs, func(i, j int) bool {
		return fs[i].ID < fs[j].ID
	})

	// Output.
	g.NL()
	g.Printf("var formulae = []*Formula")
	g.EnterBlock()
	for _, f := range fs {
		g.formula(f)
	}
	g.LeaveBlock()
}

func (g *generator) formula(f *efd.Formula) {
	g.EnterBlock()
	g.Linef("ID: %#v,", f.ID)
	g.Linef("Tag: %#v,", f.Tag)
	g.Linef("Class: %#v,", f.Class)
	g.Linef("Shape: %s,", g.shaperef(f.Shape))
	g.Linef("Representation: %s,", g.representationref(f.Representation))
	g.Linef("Operation: %#v,", f.Operation)
	g.Linef("Source: %#v,", f.Source)
	g.Linef("AppliesTo: %#v,", f.AppliesTo)
	g.Linef("Assume: %#v,", f.Assume)
	g.Linef("Compute: %#v,", f.Compute)
	g.Linef("Parameters: %#v,", f.Parameters)
	g.Linef("Program: %#v,", f.Program)
	g.Dedent()
	g.Linef("},")
}
