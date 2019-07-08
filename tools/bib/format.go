package main

import (
	"strings"

	"golang.org/x/xerrors"
)

// Format entry as a string.
func Format(e *Entry) (string, error) {
	// For simplicity assume author and title.
	s := FormatAuthors(e.Authors())
	if !strings.HasSuffix(s, ".") {
		s += "."
	}
	s += " " + e.Fields["title"].String() + "."

	// Custom fields.
	switch e.Type {
	case "misc":
		// Optional fields: author, title, howpublished, month, year, note.
		if how, found := e.Fields["howpublished"]; found {
			s += " " + how.String() + "."
		}

		if license, found := e.Fields["license"]; found {
			s += " " + license.String() + "."
		}

	case "inproceedings":
		// Required fields: author, title, booktitle, year.
		s += " In " + e.Fields["booktitle"].String()
		if pages, found := e.Fields["pages"]; found {
			s += ", pages " + pages.String()
		}
		s += "."

	case "article":
		// Required fields: author, title, journal, year.
		if journal, found := e.Fields["journal"]; found {
			s += " " + journal.String() + "."
		}

	case "inbook":
		// Required fields: author or editor, title, chapter and/or pages, publisher, year.
		s += " " + e.Fields["booktitle"].String()
		s += ", chapter " + e.Fields["chapter"].String() + "."

	case "phdthesis":
		// Required fields: author, title, school, year.
		s += " PhD thesis, " + e.Fields["school"].String() + "."

	case "mastersthesis":
		// Required fields: author, title, school, year.
		s += " Masters thesis, " + e.Fields["school"].String() + "."

	case "techreport":
		// Required fields: author, title, institution, year.
		// Optional fields: type, number, address, month, note.
		s += " Technical Report " + e.Fields["number"].String()
		s += ", " + e.Fields["institution"].String() + "."

	default:
		return "", xerrors.Errorf("unknown entry type '%s'", e.Type)
	}

	// Look for a date.
	if year, found := e.Fields["year"]; found {
		s += " " + year.String() + "."
	}

	// Always look for a URL.
	if url, found := e.Fields["url"]; found {
		s += " " + url.String()
	}

	if accessed, err := e.DateField("urldate"); err == nil {
		s += " (accessed " + accessed.Format("January _2, 2006") + ")"
	}

	return s, nil
}

// FormatAuthors formats a list of authors in a readable form.
func FormatAuthors(authors []string) string {
	n := len(authors)
	switch n {
	case 0:
		return ""
	case 1:
		return authors[0]
	default:
		return strings.Join(authors[:n-1], ", ") + " and " + authors[n-1]
	}
}

// Wrap text into lines of length at most width.
func Wrap(text string, width int) []string {
	words := strings.Fields(text)
	lines := []string{}
	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+len(word) > width {
			lines = append(lines, line)
			line = word
		} else {
			line += " " + word
		}
	}
	lines = append(lines, line)
	return lines
}
