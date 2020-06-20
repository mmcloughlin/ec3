package main

import (
	"net/http"

	"golang.org/x/xerrors"
)

// DatabaseLinks extracts all links from a database.
func DatabaseLinks(db *Database) []string {
	var links []string
	for _, reference := range db.References {
		links = append(links, reference.URL)
		for _, supp := range reference.Supplements {
			links = append(links, supp.URL)
		}
	}
	return links
}

// CheckLink checks whether the given URL exists.
func CheckLink(u string) (err error) {
	r, err := http.Get(u)
	if err != nil {
		return err
	}
	defer func() {
		if errc := r.Body.Close(); errc != nil && err == nil {
			err = errc
		}
	}()

	if r.StatusCode != http.StatusOK {
		return xerrors.Errorf("http status %d", r.StatusCode)
	}

	return nil
}
