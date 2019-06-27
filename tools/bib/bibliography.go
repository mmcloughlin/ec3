package main

import (
	"os"

	"github.com/nickng/bibtex"
)

// Bibliography is a collection of references.
type Bibliography struct {
	Entries []*bibtex.BibEntry
}

// ReadBibliography reads entries from the given BiBTeX file.
func ReadBibliography(path string) (*Bibliography, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := bibtex.Parse(f)
	if err != nil {
		return nil, err
	}

	return &Bibliography{
		Entries: b.Entries,
	}, nil
}

// Lookup reference with the given key.
func (b *Bibliography) Lookup(key string) *bibtex.BibEntry {
	for _, e := range b.Entries {
		if e.CiteName == key {
			return e
		}
	}
	return nil
}
