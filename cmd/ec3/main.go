package main

import (
	"crypto/elliptic"
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
	"github.com/mmcloughlin/ec3/gen/curve"
	"github.com/mmcloughlin/ec3/gen/ec"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/gen/name"
	"github.com/mmcloughlin/ec3/prime"
)

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	directory = flags.String("dir", "", "directory to write to")

	inverse       = flags.String("inv", "", "addition chain for field inversion")
	scalarinverse = flags.String("scalarinv", "", "addition chain for scalar field inversion")

	repr = flags.String("repr", "", "curve representation")
)

func main() {
	flags.Parse(os.Args[1:])

	// Load inversion chains.
	if *inverse == "" {
		log.Fatal("must provide addition chain for inversion")
	}
	p, err := acc.LoadFile(*inverse)
	if err != nil {
		log.Fatal(err)
	}

	if *scalarinverse == "" {
		log.Fatal("must provide addition chain for scalar inversion")
	}
	scalarinvp, err := acc.LoadFile(*scalarinverse)
	if err != nil {
		log.Fatal(err)
	}

	// Build file set.
	// fs := fp25519(p)
	fs := p256(p, scalarinvp)

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

func p256(p, scalarinvp *ir.Program) gen.Files {
	params := elliptic.P256().Params()

	// Field config.
	fieldcfg := fp.Config{
		Field:        mont.New(prime.NISTP256),
		InverseChain: p,

		PackageName:     "p256",
		ElementTypeName: "Elt",
		FilenamePrefix:  "fp",
		Scheme:          name.Plain,
	}

	fieldfiles, err := fp.Package(fieldcfg)
	if err != nil {
		log.Fatal(err)
	}

	// Scalar field config.
	scalarcfg := fp.Config{
		Field:        mont.New(prime.NewOther(params.N)),
		InverseChain: scalarinvp,

		PackageName:     "p256",
		ElementTypeName: "scalar",
		FilenamePrefix:  "scalar",
		Scheme: name.CompositeScheme(
			name.Prefixed("scalar"),
			name.LowerCase,
		),
	}

	scalarfiles, err := fp.Package(scalarcfg)
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

	affine := ec.Representation{
		Name:        "Affine",
		ElementType: fieldcfg.Type(),
		Coordinates: affinecoords,
	}

	jacobian := ec.Representation{
		Name:        "Jacobian",
		ElementType: fieldcfg.Type(),
		Coordinates: r.Variables,
	}

	// TODO(mbm): automatically generate conversion formulae
	fromaffine := ec.Function{
		Name: "NewFromAffine",
		Params: []ec.Parameter{
			ec.Point("a", ec.R, affine, 1),
		},
		Results: []ec.Parameter{
			ec.Point("p", ec.W, jacobian, 3),
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
		Receiver: ec.Point("p", ec.R, jacobian, 1),
		Results: []ec.Parameter{
			ec.Point("a", ec.W, affine, 3),
		},
		Formula: scalef.Program,
	}

	// TODO(mbm): automatically generate cmov formulae
	cmov := ec.Function{
		Name:     "CMov",
		Receiver: ec.Point("p", ec.W, jacobian, 3),
		Params: []ec.Parameter{
			ec.Point("q", ec.R, jacobian, 1),
			ec.Condition("c", ec.R),
		},
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
		Name:     "CNeg",
		Receiver: ec.Point("p", ec.RW, jacobian, 1, 3),
		Params: []ec.Parameter{
			ec.Condition("c", ec.R),
		},
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
		Receiver: ec.Point("p", ec.W, jacobian, 3),
		Params: []ec.Parameter{
			ec.Point("q", ec.R, jacobian, 1),
			ec.Point("r", ec.R, jacobian, 2),
		},
		Formula: addf.Program,
	}

	dblf := efd.LookupFormula("g1p/shortw/jacobian-3/doubling/dbl-2001-b")
	if dblf == nil {
		log.Fatal("unknown formula")
	}

	dbl := ec.Function{
		Name:     "Double",
		Receiver: ec.Point("p", ec.W, jacobian, 3),
		Params: []ec.Parameter{
			ec.Point("q", ec.R, jacobian, 1),
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

	// Curve operations.
	shortw := curve.ShortWeierstrass{
		PackageName: "p256",
		Params:      params,
		ShortName:   "p256",
	}

	curvefiles, err := shortw.Generate()
	if err != nil {
		log.Fatal(err)
	}

	// Merge and output.
	return gen.Merge(fieldfiles, scalarfiles, pointfiles, curvefiles)
}
