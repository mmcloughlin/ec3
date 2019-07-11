package db

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

// Visitor provides a Visit method that is called for every file visited by
// Walk. The visit method receives the filename within the archive, as well as a
// reader for the contents.
type Visitor interface {
	Visit(filename string, r io.Reader) error
}

// VisitorFunc adapts a plain function to the Visitor interface.
type VisitorFunc func(filename string, r io.Reader) error

// Visit calls f.
func (f VisitorFunc) Visit(filename string, r io.Reader) error {
	return f(filename, r)
}

// Walk walks the EFD archive file at filename, calling the given function for
// every file.
func Walk(filename string, v Visitor) error {
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

		if err := v.Visit(hdr.Name, tr); err != nil {
			return err
		}
	}
}
