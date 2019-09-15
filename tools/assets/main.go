// Command assets embeds static assets into Go source code.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Command line flags.
var (
	output   = flag.String("output", "", "path to output file (default stdout)")
	pkg      = flag.String("pkg", "", "package name")
	funcname = flag.String("func", "Asset", "function name")
)

func main() {
	log.SetPrefix("assets: ")
	log.SetFlags(0)

	flag.Parse()

	process(flag.Args())
}

func process(patterns []string) {
	filenames, err := expand(patterns)
	if err != nil {
		log.Fatal(err)
	}

	assets, err := LoadAssets(filenames)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{
		PackageName: *pkg,
		GeneratedBy: "assets",
		Function:    *funcname,
	}

	b, err := Generate(cfg, assets)
	if err != nil {
		log.Fatal(err)
	}

	if *output != "" {
		err = ioutil.WriteFile(*output, b, 0644)
	} else {
		_, err = os.Stdout.Write(b)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// expand filename patterns into a full list of files.
func expand(patterns []string) ([]string, error) {
	filenames := []string{}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, matches...)
	}
	return filenames, nil
}
