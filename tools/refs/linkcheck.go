package main

import (
	"net/http"

	"golang.org/x/xerrors"
)

// CheckLink checks whether the given URL exists.
func CheckLink(u string) error {
	r, err := http.Get(u)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return xerrors.Errorf("http status %d", r.StatusCode)
	}

	return nil
}
