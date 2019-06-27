package main

import (
	"flag"
	"log"
	"os"
)

// Commandline flags.
var (
	bib = flag.String("bib", "", "bibliography file")
)

func main() {
	log.SetPrefix("bib: ")
	log.SetFlags(0)

	flag.Parse()

	b, err := ReadBibliography(*bib)
	if err != nil {
		log.Fatal(err)
	}

	for _, filename := range flag.Args() {
		process(filename, b)
	}
}

func process(filename string, b *Bibliography) {
	s, err := ParseFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Validate(b); err != nil {
		log.Fatal(err)
	}

	if err = s.Write(os.Stdout, b); err != nil {
		log.Fatal(err)
	}
}
