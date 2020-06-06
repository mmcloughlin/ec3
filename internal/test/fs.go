package test

import (
	"io/ioutil"
	"os"
	"testing"
)

// TempDir creates a temp directory. Returns the path to the directory.
func TempDir(t *testing.T) string {
	t.Helper()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})

	return dir
}
