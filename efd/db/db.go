package db

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

// WalkFunc is the type of function called for every file visited by Walk. The
// filename argument is the name of the file within the tar file, and r is a
// reader over its contents.
type WalkFunc func(filename string, r io.Reader) error

// Walk walks the EFD archive file at filename, calling the given function for
// every file.
func Walk(filename string, fn WalkFunc) error {
	// Build tar reader for compressed archive.
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	zr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer zr.Close()

	tr := tar.NewReader(zr)

	// Iterate through files.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		if err := fn(hdr.Name, tr); err != nil {
			return err
		}
	}
}
