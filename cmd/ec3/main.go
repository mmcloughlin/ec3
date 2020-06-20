package main

import (
	"crypto/elliptic"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/ir"

	"github.com/mmcloughlin/ec3/asm/fp/mont"
	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/gen/curve"
	"github.com/mmcloughlin/ec3/gen/fmla"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/name"
	"github.com/mmcloughlin/ec3/prime"
)

var (
	directory = flag.String("dir", "", "directory to write to")

	inverse       = flag.String("inv", "", "addition chain for field inversion")
	scalarinverse = flag.String("scalarinv", "", "addition chain for scalar field inversion")
)

func main() {
	flag.Parse()

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
	fs, err := p256(p, scalarinvp)
	if err != nil {
		log.Fatal(err)
	}

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

func p256(p, scalarinvp *ir.Program) (gen.Files, error) {
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
		return nil, err
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
		return nil, err
	}

	// Point config.
	shape := efd.LookupShape("g1p/shortw")
	if shape == nil {
		return nil, errors.New("unknown shape")
	}

	affinecoords := []string{}
	for _, v := range shape.Coordinates {
		affinecoords = append(affinecoords, strings.ToUpper(v))
	}

	affine := fmla.Representation{
		Name:        "Affine",
		ElementType: fieldcfg.Type(),
		Coordinates: affinecoords,
	}

	reprjac := efd.LookupRepresentation("g1p/shortw/jacobian-3")
	if reprjac == nil {
		return nil, errors.New("unknown representation")
	}

	jacobian := fmla.Representation{
		Name:        "Jacobian",
		ElementType: fieldcfg.Type(),
		Coordinates: reprjac.Variables,
	}

	reprproj := efd.LookupRepresentation("g1p/shortw/projective-3")
	if reprproj == nil {
		return nil, errors.New("unknown representation")
	}

	projective := fmla.Representation{
		Name:        "Projective",
		ElementType: fieldcfg.Type(),
		Coordinates: reprproj.Variables,
	}

	// TODO(mbm): automatically generate conversion formulae
	atoj := fmla.Function{
		Name:     "Jacobian",
		Receiver: fmla.Point("a", fmla.R, affine, 1),
		Results: []fmla.Parameter{
			fmla.Point("p", fmla.W, jacobian, 3),
		},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "X3", RHS: ast.Variable("X1")},
				{LHS: "Y3", RHS: ast.Variable("Y1")},
				{LHS: "Z3", RHS: ast.Constant(1)},
			},
		},
	}

	atop := fmla.Function{
		Name:     "Projective",
		Receiver: fmla.Point("a", fmla.R, affine, 1),
		Results: []fmla.Parameter{
			fmla.Point("p", fmla.W, projective, 3),
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
		return nil, errors.New("unknown formula")
	}

	jtoa := fmla.Function{
		Name:     "Affine",
		Receiver: fmla.Point("p", fmla.R, jacobian, 1),
		Results: []fmla.Parameter{
			fmla.Point("a", fmla.W, affine, 3),
		},
		Formula: scalef.Program,
	}

	// TODO(mbm): improve handling of conversion formulae
	jtop := fmla.Function{
		Name:     "Projective",
		Receiver: fmla.Point("p", fmla.R, jacobian, 1),
		Results: []fmla.Parameter{
			fmla.Point("q", fmla.W, projective, 3),
		},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "X3", RHS: ast.Mul{X: ast.Variable("X1"), Y: ast.Variable("Z1")}},
				{LHS: "Y3", RHS: ast.Variable("Y1")},
				{LHS: "Z3", RHS: ast.Pow{X: ast.Variable("Z1"), N: 3}},
			},
		},
	}

	pscalef := efd.LookupFormula("g1p/shortw/projective-3/scaling/z")
	if pscalef == nil {
		return nil, errors.New("unknown formula")
	}

	ptoa := fmla.Function{
		Name:     "Affine",
		Receiver: fmla.Point("p", fmla.R, projective, 1),
		Results: []fmla.Parameter{
			fmla.Point("a", fmla.W, affine, 3),
		},
		Formula: pscalef.Program,
	}

	// Lookup formula.
	lookup := fmla.Lookup{
		Name: "lookup",
		Repr: jacobian,
	}

	// TODO(mbm): automatically generate cmov formulae
	cmov := fmla.Function{
		Name:     "CMov",
		Receiver: fmla.Point("p", fmla.W, jacobian, 3),
		Params: []fmla.Parameter{
			fmla.Point("q", fmla.R, jacobian, 1),
			fmla.Condition("c", fmla.R),
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
	jcneg := fmla.Function{
		Name:     "CNeg",
		Receiver: fmla.Point("p", fmla.RW, jacobian, 1, 3),
		Params: []fmla.Parameter{
			fmla.Condition("c", fmla.R),
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
		return nil, errors.New("unknown formula")
	}

	add := fmla.NewAsmFunctionDefault(fmla.Function{
		Name:     "Add",
		Receiver: fmla.Point("p", fmla.W, jacobian, 3),
		Params: []fmla.Parameter{
			fmla.Point("q", fmla.R, jacobian, 1),
			fmla.Point("r", fmla.R, jacobian, 2),
		},
		Formula: addf.Program,
	})

	dblf := efd.LookupFormula("g1p/shortw/jacobian-3/doubling/dbl-2001-b")
	if dblf == nil {
		return nil, errors.New("unknown formula")
	}

	dbl := fmla.NewAsmFunctionDefault(fmla.Function{
		Name:     "Double",
		Receiver: fmla.Point("p", fmla.W, jacobian, 3),
		Params: []fmla.Parameter{
			fmla.Point("q", fmla.R, jacobian, 1),
		},
		Formula: dblf.Program,
	})

	b := fmla.Constant{
		VariableName: "b",
		ElementType:  fieldcfg.Type(),
		Value:        params.B,
	}

	pcneg := fmla.Function{
		Name:     "CNeg",
		Receiver: fmla.Point("p", fmla.RW, projective, 1, 3),
		Params: []fmla.Parameter{
			fmla.Condition("c", fmla.R),
		},
		Formula: &ast.Program{
			Assignments: []ast.Assignment{
				{LHS: "t", RHS: ast.Neg{X: ast.Variable("Y1")}},
				{LHS: "Y3", RHS: ast.Cond{X: ast.Variable("t"), C: ast.Variable("c")}},
			},
		},
	}

	compaddf := efd.LookupFormula("g1p/shortw/projective-3/addition/add-2015-rcb")
	if compaddf == nil {
		return nil, errors.New("unknown formula")
	}

	compadd := fmla.NewAsmFunctionDefault(fmla.Function{
		Name:     "CompleteAdd",
		Receiver: fmla.Point("p", fmla.W, projective, 3),
		Params: []fmla.Parameter{
			fmla.Point("q", fmla.R, projective, 1),
			fmla.Point("r", fmla.R, projective, 2),
		},
		Globals: []fmla.Parameter{b},
		Formula: compaddf.Program,
	})

	pointcfg := fmla.Config{
		PackageName: "p256",
		Field:       fieldcfg,
		Components: []fmla.Component{
			// Constants.
			b,

			// Affine representation.
			affine,
			atoj,
			atop,

			// Jacobian representation.
			jacobian,
			jtoa,
			jtop,
			lookup,
			cmov,
			jcneg,
			add,
			dbl,

			// Projective representation.
			projective,
			ptoa,
			pcneg,
			compadd,
		},
	}

	pointfiles, err := fmla.Package(pointcfg)
	if err != nil {
		return nil, err
	}

	// Curve operations.
	shortw := curve.ShortWeierstrass{
		PackageName: "p256",
		Params:      params,
		ShortName:   "p256",
	}

	curvefiles, err := shortw.Generate()
	if err != nil {
		return nil, err
	}

	// Merge and output.
	return gen.Merge(fieldfiles, scalarfiles, pointfiles, curvefiles), nil
}
