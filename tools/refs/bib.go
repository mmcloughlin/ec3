package main

import "github.com/nickng/bibtex"

// DatabaseBibTex generates a BibTeX database from a references database. Only
// entries with an ID are included.
func DatabaseBibTex(db *Database) *bibtex.BibTex {
	bib := bibtex.NewBibTex()
	for _, ref := range db.References {
		if e := ReferenceBibEntry(ref); e != nil {
			bib.AddEntry(e)
		}
	}
	return bib
}

// ReferenceBibEntry returns a BibTeX entry corresponding to r, if possible. The
// reference ID is used as the citation name, and the type field is used as the
// BibTeX type (defaulting to misc).
func ReferenceBibEntry(r *Reference) *bibtex.BibEntry {
	// Skip those without IDs.
	if r.ID == "" {
		return nil
	}

	// Type defaults to misc.
	t := r.Type
	if t == "" {
		t = "misc"
	}

	// Build the entry.
	e := bibtex.NewBibEntry(t, r.ID)

	// Populate fields.
	fields := map[string]string{}
	for k, v := range r.Fields {
		fields[k] = v
	}

	fields["title"] = r.Title
	fields["url"] = r.URL
	fields["author"] = r.Author

	for k, v := range fields {
		e.AddField(k, bibtex.NewBibConst(v))
	}

	return e
}
