package db

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/internal/errutil"
)

// File in a Store.
type File interface {
	// Path relative to Store root.
	Path() string

	io.ReadCloser
}

type file struct {
	path string
	io.ReadCloser
}

func (f file) Path() string { return f.path }

// Store is a method of accessing database files.
type Store interface {
	Next() (File, error)
	io.Closer
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

// Open is a convenience for opening a database. It will make an informed guess
// about the storage format.
func Open(filename string) (Store, error) {
	s, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	switch {
	case s.IsDir():
		return Directory(filename)
	case strings.HasSuffix(filename, ".tar.gz"):
		return Archive(filename)
	default:
		return nil, xerrors.Errorf("unknown database type for file %q", filename)
	}
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
			path:       hdr.Name,
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

// directory provides access to files under a directory.
type directory struct {
	root  string
	paths []string
}

// Directory opens a database stored at root.
func Directory(root string) (Store, error) {
	d := &directory{root: root}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			d.paths = append(d.paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Next returns the next file in the directory.
func (d *directory) Next() (File, error) {
	if len(d.paths) == 0 {
		return nil, io.EOF
	}
	path := d.paths[0]
	d.paths = d.paths[1:]

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(d.root, path)
	if err != nil {
		return nil, errutil.AssertionFailure("failed to construct relative path: %w", err)
	}

	return file{
		path:       rel,
		ReadCloser: f,
	}, nil
}

// Close is a no-op.
func (directory) Close() error { return nil }
