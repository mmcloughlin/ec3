package parse_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/mmcloughlin/ec3/efd/db"
	"github.com/mmcloughlin/ec3/efd/op3/parse"
	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestParseAllEFDFiles(t *testing.T) {
	archive, err := db.Archive("../../data/efd.tar.gz")
	assert.NoError(t, err)

	err = db.Walk(archive, db.VisitorFunc(func(f db.File) error {
		filename := f.Path()

		if filepath.Ext(filename) != ".op3" {
			return nil
		}

		t.Logf("parsing %s", filename)

		// Read into a byte array.
		b, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		// Some files have error messages in them, in which case we shouldn't expect
		// to parse.
		expecterr := bytes.Contains(b, []byte("error:"))

		// Parse.
		_, err = parse.Reader(filename, bytes.NewReader(b))

		if expecterr && err == nil {
			t.Fatal("expected an error")
		}

		if !expecterr && err != nil {
			t.Fatal(err)
		}

		return nil
	}))

	assert.NoError(t, err)
}
