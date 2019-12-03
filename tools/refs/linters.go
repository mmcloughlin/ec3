package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Linter checks a database for errors.
type Linter interface {
	Lint(*Database) []error
}

// LinterFunc adapts a function to the Linter interface.
type LinterFunc func(*Database) []error

// Lint calls f.
func (f LinterFunc) Lint(db *Database) []error {
	return f(db)
}

// ConcatLinter builds a composite linter from the given linter.
func ConcatLinter(linters ...Linter) Linter {
	return LinterFunc(func(db *Database) []error {
		var errs []error
		for _, l := range linters {
			if e := l.Lint(db); len(e) > 0 {
				errs = append(errs, e...)
			}
		}
		return errs
	})
}

// ReferenceLinterFunc adapts a function that lints a single reference to a
// linter for an entire database.
type ReferenceLinterFunc func(*Reference) []error

// Lint calls f on every reference in the database and concatencates the results.
func (f ReferenceLinterFunc) Lint(db *Database) []error {
	var errs []error
	for _, r := range db.References {
		if e := f(r); len(e) > 0 {
			errs = append(errs, e...)
		}
	}
	return errs
}

// RequireURL checks that a reference has a URL field.
func RequireURL(r *Reference) []error {
	if r.URL != "" {
		return nil
	}
	return singleerror("missing url")
}

// RequireTitle checks that a reference has a title.
func RequireTitle(r *Reference) []error {
	if r.Title != "" {
		return nil
	}
	return singleerror("missing title")
}

// ValidURL checks that a reference has a valid URL.
func ValidURL(r *Reference) []error {
	_, err := url.Parse(r.URL)
	if err != nil {
		return []error{err}
	}
	return nil
}

// AuthorPeriod confirms the author field does not end in a period.
func AuthorPeriod(r *Reference) []error {
	if strings.HasSuffix(r.Author, ".") {
		return singleerror("author %q ends with period", r.Author)
	}
	return nil
}

// IACRCanonical ensures that IACR eprint references conform to a canonical form.
func IACRCanonical(r *Reference) []error {
	u, err := url.Parse(r.URL)
	if err != nil {
		return []error{err}
	}

	const (
		ShortHost = "ia.cr"
		Host      = "eprint.iacr.org"
	)

	if u.Host == ShortHost {
		return singleerror("url %s: prefer host %s", u, Host)
	}

	if u.Host != Host {
		return nil
	}

	if u.Scheme != "https" {
		return singleerror("url %s: require https", u)
	}

	match, err := regexp.MatchString(`^/\d{4}/\d+$`, u.Path)
	if err != nil {
		return []error{err}
	}

	if !match {
		return singleerror("url %s: link to the report page", u)
	}

	return nil
}

// DisallowHost builds a linter that errors for URLs with the given host.
func DisallowHost(host string) Linter {
	return ReferenceLinterFunc(func(r *Reference) []error {
		u, err := url.Parse(r.URL)
		if err != nil {
			return []error{err}
		}

		if u.Host == host {
			return singleerror("host %s disallowed", host)
		}

		return nil
	})
}

// DuplicateURLs checks for duplicate URLs on the reference database.
func DuplicateURLs(db *Database) []error {
	count := map[string]int{}
	for _, ref := range db.References {
		count[ref.URL]++
	}

	var errs []error
	for u, n := range count {
		if n > 1 {
			errs = append(errs, fmt.Errorf("url %s occurs %d times", u, n))
		}
	}

	return errs
}

// CheckNewlines confirms fields do not have newlines in them.
func CheckNewlines(r *Reference) (errs []error) {
	fields := []string{r.Title, r.Author, r.Note}
	for _, field := range fields {
		if strings.Contains(field, "\n") {
			err := fmt.Errorf("%q contains newline", field)
			errs = append(errs, err)
		}
	}
	return
}

// CheckSectionTags verifies that every section tag is defined.
func CheckSectionTags(db *Database) (errs []error) {
	// Collect valid section IDs.
	defined := map[string]bool{}
	for _, section := range db.Sections {
		defined[section.ID] = true
	}

	// Check each reference.
	undef := map[string]bool{}
	for _, ref := range db.References {
		s := ref.Section
		if s != "" && !defined[s] {
			undef[s] = true
		}
	}

	for name := range undef {
		err := fmt.Errorf("section %q undefined", name)
		errs = append(errs, err)
	}

	return errs
}

func singleerror(format string, args ...interface{}) []error {
	return []error{fmt.Errorf(format, args...)}
}
