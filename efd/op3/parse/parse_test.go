package parse

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestParseBasic(t *testing.T) {
	s := "X = 3*P\nY = 1+X\n"
	p, err := String(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", p)
}

func TestParseAllEFDFiles(t *testing.T) {
	archive := "../../efd.tar.gz"

	// Build tar reader for compressed archive.
	f, err := os.Open(archive)
	assert.NoError(t, err)
	defer f.Close()

	zr, err := gzip.NewReader(f)
	assert.NoError(t, err)
	defer zr.Close()

	tr := tar.NewReader(zr)

	// Iterate through files.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return
		}
		assert.NoError(t, err)

		if filepath.Ext(hdr.Name) != ".op3" {
			continue
		}

		AssertParses(t, hdr.Name, tr)
	}
}

func AssertParses(t *testing.T, filename string, r io.Reader) {
	t.Helper()

	t.Logf("parsing %s", filename)

	// Read into a byte array.
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	// Some files have error messages in them, in which case we shouldn't expect
	// to parse.
	expecterr := bytes.Contains(b, []byte("error:"))

	// Parse.
	_, err = Reader(filename, bytes.NewReader(b))

	if expecterr && err == nil {
		t.Fatal("expected an error")
	}

	if !expecterr && err != nil {
		t.Fatal(err)
	}
}
