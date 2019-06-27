package main

import (
	"fmt"
	"strings"
)

// Format entry as a string.
func Format(e *Entry) (string, error) {
	// For simplicity assume author and title.
	s := FormatAuthors(e.Authors()) + "."
	s += " " + e.Fields["title"].String() + "."

	// Custom fields.
	switch e.Type {
	case "misc":
		// Required: author/editor, title, year/date
		if how, found := e.Fields["howpublished"]; found {
			s += " " + how.String() + "."
		}

	case "inproceedings":
		// Required: author, title, booktitle, year/date.
		s += " In " + e.Fields["booktitle"].String()
		if pages, found := e.Fields["pages"]; found {
			s += ", pages " + pages.String()
		}
		s += "."

	default:
		return "", fmt.Errorf("unknown entry type '%s'", e.Type)
	}

	// Look for a date.
	if year, found := e.Fields["year"]; found {
		s += " " + year.String() + "."
	}

	// Always look for a URL.
	if url, found := e.Fields["url"]; found {
		s += " " + url.String()
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
