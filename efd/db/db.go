package db

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/mmcloughlin/ec3/internal/errutil"
)

// File in a Store.
type File interface {
	Name() string
	io.ReadCloser
}

type file struct {
	name string
	io.ReadCloser
}

func (f file) Name() string { return f.name }

// Store is a method of accessing database files.
type Store interface {
	Next() (File, error)
	io.Closer
}

// archive provides access to a gzipped tarball.
type archive struct {
	f *os.File
	z *gzip.Reader
	r *tar.Reader
}

// Archive opens a gzipped tarball database.
func Archive(filename string) (Store, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	z, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}

	return archive{
		f: f,
		z: z,
		r: tar.NewReader(z),
	}, nil
}

// Next returns the next file in the archive.
func (a archive) Next() (File, error) {
	for {
		hdr, err := a.r.Next()
		if err == io.EOF {
			return nil, io.EOF
		}
		if err != nil {
			return nil, err
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		return file{
			name:       hdr.Name,
			ReadCloser: ioutil.NopCloser(a.r),
		}, nil
	}
}

// Close the archive.
func (a archive) Close() error {
	var errs errutil.Errors
	for _, c := range []io.Closer{a.z, a.f} {
		if err := c.Close(); err != nil {
			errs.Add(err)
		}
	}
	return errs.Err()
}

// Visitor provides a Visit method that is called for every file visited by
// Walk. The visit method receives the filename within the Store, as well as a
// reader for the contents.
type Visitor interface {
	Visit(f File) error
}

// VisitorFunc adapts a plain function to the Visitor interface.
type VisitorFunc func(f File) error

// Visit calls v.
func (v VisitorFunc) Visit(f File) error {
	return v(f)
}

// Walk walks the Store, calling the Visitor for every file.
func Walk(s Store, v Visitor) error {
	defer s.Close()
	for {
		f, err := s.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := v.Visit(f); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}
}
