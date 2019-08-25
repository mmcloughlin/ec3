package tmpl

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mmcloughlin/ec3/internal/assert"
)

func Corpus(t *testing.T) []string {
	t.Helper()
	filenames, err := filepath.Glob("*.go")
	assert.NoError(t, err)
	return filenames
}

func TestRoundTrip(t *testing.T) {
	for _, filename := range Corpus(t) {
		src, err := FileSystemLoader.Load(filename)
		assert.NoError(t, err)

		tpl, err := ParseFile(filename, src)
		assert.NoError(t, err)

		got, err := tpl.Bytes()
		assert.NoError(t, err)

		if !bytes.Equal(got, src) {
			t.Logf("src:\n%s", src)
			t.Logf("got:\n%s", got)
			t.Fatalf("%s: roundtrip error", filename)
		}
	}
}
