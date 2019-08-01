package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mmcloughlin/ec3/addchain/acc"
	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/asm/fp/crandall"
	"github.com/mmcloughlin/ec3/asm/fp/mont"
	"github.com/mmcloughlin/ec3/efd"
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
	//fs := fp25519(p)
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

	f := efd.LookupFormula("g1p/shortw/jacobian-3/addition/add-2007-bl")
	if f == nil {
		log.Fatalf("unknown formula")
	}

	pointcfg := ec.Config{
		Field:          fieldcfg,
		Representation: r,
		Functions: []ec.Function{
			{
				Name:       "Add",
				Parameters: []string{"q", "r"},
				Formula:    f,
			},
		},
		PackageName: "p256",
		TypeName:    "Jacobian",
	}

	pointfiles, err := ec.Package(pointcfg)
	if err != nil {
		log.Fatal(err)
	}

	// Merge and output.
	return gen.Merge(fieldfiles, pointfiles)
}
