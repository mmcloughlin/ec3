package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/prime"
)

func main() {
	cfg := fp.Config{
		PackageName:     "fp25519",
		Prime:           prime.P25519,
		ElementTypeName: "Elt",
	}

	b, err := fp.Package(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range b.Files {
		fmt.Printf("## `%s`\n", f.Path)
		fmt.Printf("```\n")
		b, err := f.Source.Generate()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stdout.Write(b); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("```\n")
	}
}
