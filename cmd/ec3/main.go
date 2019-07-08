package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/prime"
)

var (
	flags     = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	directory = flags.String("dir", "", "directory to write to")
)

func main() {
	flags.Parse(os.Args[1:])

	cfg := fp.Config{
		PackageName:     "fp25519",
		Prime:           prime.P25519,
		ElementTypeName: "Elt",
	}

	fs, err := fp.Package(cfg)
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
