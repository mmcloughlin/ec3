package db

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/internal/test"
)

func TestDirectory(t *testing.T) {
	// Setup a temporary directory.
	dir, clean := test.TempDir(t)
	defer clean()

	paths := []string{"one", "two", "three"}
	filename := "data"
	contents := "Hello, World!"

	root := dir
	for _, path := range paths {
		root = filepath.Join(root, path)

		err := os.Mkdir(root, 0755)
		assert.NoError(t, err)

		ioutil.WriteFile(filepath.Join(root, filename), []byte(contents), 0640)
		assert.NoError(t, err)
	}

	expect := map[string]string{
		"one/data":           contents,
		"one/two/data":       contents,
		"one/two/three/data": contents,
	}

	// Read files with Directory Store.
	d, err := Directory(dir)
	assert.NoError(t, err)

	got := map[string]string{}
	err = Walk(d, VisitorFunc(func(f File) error {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		got[f.Path()] = string(data)
		return nil
	}))
	assert.NoError(t, err)

	if !reflect.DeepEqual(expect, got) {
		t.Logf("   got = %v", got)
		t.Logf("expect = %v", expect)
		t.Fail()
	}
}
