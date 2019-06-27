package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Commandline flags.
var (
	bibfile = flag.String("bib", "", "bibliography file")
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

func main() {
	log.SetPrefix("bib: ")
	log.SetFlags(0)

	flag.Parse()

	b, err := ReadBibliography(*bibfile)
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

	out, err := s.Bytes(b)
	if err != nil {
		log.Fatal(err)
	}

	if *write {
		err = ioutil.WriteFile(filename, out, 0644)
	} else {
		_, err = os.Stdout.Write(out)
	}

	if err != nil {
		log.Fatal(err)
	}
}
