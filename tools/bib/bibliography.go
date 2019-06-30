package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nickng/bibtex"
)

// Entry in a bibliography.
type Entry struct {
	bibtex.BibEntry
}

// Authors returns the list of authors.
func (e Entry) Authors() []string {
	authors := strings.Split(e.Fields["author"].String(), " and ")
	for i := range authors {
		authors[i] = strings.TrimSpace(authors[i])
	}
	return authors
}

// DateField parses a field as a date in ISO 8601 format.
func (e Entry) DateField(name string) (time.Time, error) {
	s, ok := e.Fields[name]
	if !ok {
		return time.Time{}, errors.New("field not found")
	}
	return time.Parse("2006-01-02", s.String())
}

// ByCiteName sorts a list of entries by their citation name.
type ByCiteName []*Entry

func (e ByCiteName) Len() int           { return len(e) }
func (e ByCiteName) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ByCiteName) Less(i, j int) bool { return e[i].CiteName < e[j].CiteName }

// Bibliography is a collection of references.
type Bibliography struct {
	Entries []*Entry
}

// ReadBibliography reads entries from the given BiBTeX file.
func ReadBibliography(path string) (*Bibliography, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bib, err := bibtex.Parse(f)
	if err != nil {
		return nil, err
	}

	// Build.
	b := &Bibliography{}
	for _, e := range bib.Entries {
		if err := b.AddEntry(&Entry{BibEntry: *e}); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// AddEntry adds an entry to the bibliography.
func (b *Bibliography) AddEntry(e *Entry) error {
	if b.Lookup(e.CiteName) != nil {
		return fmt.Errorf("key '%s' already in bibliography", e.CiteName)
	}
	b.Entries = append(b.Entries, e)
	return nil
}

// Lookup reference with the given key.
func (b *Bibliography) Lookup(key string) *Entry {
	for _, e := range b.Entries {
		if e.CiteName == key {
			return e
		}
	}
	return nil
}
