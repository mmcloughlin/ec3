package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/mmcloughlin/ec3/internal/print"
)

type bibtex struct {
	w io.Writer
	print.Printer
}

func WriteBibTeX(w io.Writer, db *Database) error {
	p := &bibtex{
		w:       w,
		Printer: print.New(w),
	}
	p.database(db)
	return p.Error()
}

func (p *bibtex) database(db *Database) {
	first := true
	for _, r := range db.References {
		if !isbibtexentry(r) {
			continue
		}
		if !first {
			p.NL()
		}
		p.reference(r)
		first = false
	}
}

func (p *bibtex) reference(r *Reference) {
	// Determine bibtex type.
	t := r.Type
	if t == "" {
		t = "misc"
	}

	// Open entry.
	p.Linef("@%s{%s,", t, r.ID)

	// Build fields list.
	type field struct{ name, value string }
	fields := []field{
		{"title", r.Title},
		{"author", r.Author},
		{"url", r.URL},
	}

	extras := []field{}
	for k, v := range r.Fields {
		extras = append(extras, field{k, v})
	}

	sort.Slice(extras, func(i, j int) bool {
		return extras[i].name < extras[j].name
	})
	fields = append(fields, extras...)

	// Output fields.
	tw := tabwriter.NewWriter(p.w, 1, 4, 1, ' ', 0)
	for _, f := range fields {
		format := stringformat(f.value)
		_, err := fmt.Fprintf(tw, "    %s\t=\t"+format+",\n", f.name, f.value)
		p.SetError(err)
	}
	p.SetError(tw.Flush())

	// Close.
	p.Linef("}")
}

// isbibtexentry decides whether r can be converted to a BibTeX entry.
func isbibtexentry(r *Reference) bool {
	required := []string{
		r.ID,
		r.Title,
		r.Author,
		r.URL,
	}
	for _, v := range required {
		if v == "" {
			return false
		}
	}
	return true
}

// stringformat determines the correct formatting verb for the given BibTeX field value.
func stringformat(v string) string {
	// Numbers may be represented unquoted.
	if _, err := strconv.Atoi(v); err == nil {
		return "%s"
	}

	// Strings with certain characters must be brace quoted.
	if strings.ContainsAny(v, "\"{}") {
		return "{%s}"
	}

	// Default to quoted string.
	return "%q"
}
