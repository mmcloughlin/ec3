package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Linter interface {
	Lint(*Database) []error
}

type LinterFunc func(*Database) []error

func (f LinterFunc) Lint(db *Database) []error {
	return f(db)
}

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

type ReferenceLinterFunc func(*Reference) []error

func (f ReferenceLinterFunc) Lint(db *Database) []error {
	var errs []error
	for _, r := range db.References {
		if e := f(r); len(e) > 0 {
			errs = append(errs, e...)
		}
	}
	return errs
}

func RequireURL(r *Reference) []error {
	if r.URL != "" {
		return nil
	}
	return singleerror("missing url")
}

func RequireTitle(r *Reference) []error {
	if r.Title != "" {
		return nil
	}
	return singleerror("missing title")
}

func ValidURL(r *Reference) []error {
	_, err := url.Parse(r.URL)
	if err != nil {
		return []error{err}
	}
	return nil
}

func AuthorPeriod(r *Reference) []error {
	if strings.HasSuffix(r.Author, ".") {
		return singleerror("author %q ends with period", r.Author)
	}
	return nil
}

func EprintCanonical(r *Reference) []error {
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

// Done:
// * all fields recognized
// * title and url are present, all other fields optional
// * valid url
// * author field does not end in period "."
// * eprint links do not include ".pdf"
//
// Remaining:
// * section/tags map to defined section
// * URLs are valid (link checker)
// * no fields have newlines in them (title/author/note)
// * no private links (drive.google.com)
// * repeat URLs

func singleerror(format string, args ...interface{}) []error {
	return []error{fmt.Errorf(format, args...)}
}
