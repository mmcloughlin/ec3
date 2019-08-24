package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mmcloughlin/ec3/addchain/acc"
	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/asm/fp/crandall"
	"github.com/mmcloughlin/ec3/asm/fp/mont"
	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/ec"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/prime"
)

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	directory = flags.String("dir", "", "directory to write to")
	inverse   = flags.String("inv", "", "addition chain for field inversion")

	repr = flags.String("repr", "", "curve representation")
)

func main() {
	flags.Parse(os.Args[1:])

	// Load inversion chain.
	if *inverse == "" {
		log.Fatal("must provide addition chain for inversion")
	}
	p, err := acc.LoadFile(*inverse)
	if err != nil {
		log.Fatal(err)
	}

	// Build file set.
	// fs := fp25519(p)
	fs := p256(p)

	for _, f := range fs {
		fmt.Printf("## `%s`\n", f.Path)
		fmt.Printf("```\n")
		if _, err := os.Stdout.Write(f.Source); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("```\n")
	}

	if *directory != "" {
		log.Printf("writing to %s", *directory)
		if err := fs.Output(*directory); err != nil {
			log.Fatal(err)
		}
	}
}

func fp25519(p *ir.Program) gen.Files {
	cfg := fp.Config{
		Field:        crandall.New(prime.P25519),
		InverseChain: p,

		PackageName:     "fp25519",
		ElementTypeName: "Elt",
	}

	fs, err := fp.Package(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return fs
}

func p256(p *ir.Program) gen.Files {
	// Field config.
	fieldcfg := fp.Config{
		Field:        mont.New(prime.NISTP256),
		InverseChain: p,

		PackageName:     "p256",
		ElementTypeName: "Elt",
	}

	fieldfiles, err := fp.Package(fieldcfg)
	if err != nil {
		log.Fatal(err)
	}

	// Point config.
	r := efd.LookupRepresentation(*repr)
	if r == nil {
		log.Fatalf("unknown representation %q", *repr)
	}

	affinecoords := []string{}
	for _, v := range r.Shape.Coordinates {
		affinecoords = append(affinecoords, strings.ToUpper(v))
	}

	affine := ec.Type{
		Name:        "Affine",
		ElementType: fieldcfg.Type(),
		Coordinates: affinecoords,
	}

	jacobian := ec.Type{
		Name:        "Jacobian",
		ElementType: fieldcfg.Type(),
		Coordinates: r.Variables,
	}

	// TODO(mbm): automatically generate conversion formulae
	fromaffine := ec.Function{
		Name: "NewFromAffine",
		Params: []*ec.Parameter{
			{Name: "a", Type: affine, Action: ec.R},
		},
		Results: []*ec.Parameter{
			{Name: "p", Type: jacobian, Action: ec.W},
		},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "X3", RHS: ast.Variable("X1")},
				{LHS: "Y3", RHS: ast.Variable("Y1")},
				{LHS: "Z3", RHS: ast.Constant(1)},
			},
		},
	}

	scalef := efd.LookupFormula("g1p/shortw/jacobian-3/scaling/z")
	if scalef == nil {
		log.Fatalf("unknown formula")
	}

	toaffine := ec.Function{
		Name:     "Affine",
		Receiver: &ec.Parameter{Name: "p", Type: jacobian, Action: ec.R},
		Results: []*ec.Parameter{
			{Name: "a", Type: affine, Action: ec.W},
		},
		Formula: scalef.Program,
	}

	// TODO(mbm): automatically generate cmov formulae
	cmov := ec.Function{
		Name:     "CMov",
		Receiver: &ec.Parameter{Name: "p", Type: jacobian, Action: ec.W},
		Params: []*ec.Parameter{
			{Name: "q", Type: jacobian, Action: ec.R},
		},
		Conditions: []string{"c"},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "X3", RHS: ast.Cond{X: ast.Variable("X1"), C: ast.Variable("c")}},
				{LHS: "Y3", RHS: ast.Cond{X: ast.Variable("Y1"), C: ast.Variable("c")}},
				{LHS: "Z3", RHS: ast.Cond{X: ast.Variable("Z1"), C: ast.Variable("c")}},
			},
		},
	}

	// TODO(mbm): automatically generate negation formulae
	cneg := ec.Function{
		Name:       "CNeg",
		Receiver:   &ec.Parameter{Name: "p", Type: jacobian, Action: ec.RW},
		Conditions: []string{"c"},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "t", RHS: ast.Neg{X: ast.Variable("Y1")}},
				{LHS: "Y3", RHS: ast.Cond{X: ast.Variable("t"), C: ast.Variable("c")}},
			},
		},
	}

	addf := efd.LookupFormula("g1p/shortw/jacobian-3/addition/add-2007-bl")
	if addf == nil {
		log.Fatal("unknown formula")
	}

	add := ec.Function{
		Name:     "Add",
		Receiver: &ec.Parameter{Name: "p", Type: jacobian, Action: ec.W},
		Params: []*ec.Parameter{
			{Name: "q", Type: jacobian, Action: ec.R},
			{Name: "r", Type: jacobian, Action: ec.R},
		},
		Formula: addf.Program,
	}

	dblf := efd.LookupFormula("g1p/shortw/jacobian-3/doubling/dbl-2001-b")
	if dblf == nil {
		log.Fatal("unknown formula")
	}

	dbl := ec.Function{
		Name:     "Double",
		Receiver: &ec.Parameter{Name: "p", Type: jacobian, Action: ec.W},
		Params: []*ec.Parameter{
			{Name: "q", Type: jacobian, Action: ec.R},
		},
		Formula: dblf.Program,
	}

	pointcfg := ec.Config{
		PackageName: "p256",
		Field:       fieldcfg,
		Components: []ec.Component{
			affine,
			jacobian,
			fromaffine,
			toaffine,
			cmov,
			cneg,
			add,
			dbl,
		},
	}

	pointfiles, err := ec.Package(pointcfg)
	if err != nil {
		log.Fatal(err)
	}

	// Merge and output.
	return gen.Merge(fieldfiles, pointfiles)
}
