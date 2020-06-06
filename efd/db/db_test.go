package db

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/internal/test"
)

func AssertFiles(t *testing.T, s Store, expect map[string]string) {
	got := map[string]string{}
	err := Walk(s, VisitorFunc(func(f File) error {
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
		t.FailNow()
	}
}

func TestDirectory(t *testing.T) {
	// Setup a temporary directory.
	dir := test.TempDir(t)

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

	AssertFiles(t, d, expect)
}

// Single is a mock Store that returns a single file.
func Single(path, data string) Store {
	return &single{path: path, data: data}
}

type single struct {
	path   string
	data   string
	called bool
}

func (s *single) Next() (File, error) {
	if s.called {
		return nil, io.EOF
	}
	s.called = true
	return &file{
		path:       s.path,
		ReadCloser: ioutil.NopCloser(strings.NewReader(s.data)),
	}, nil
}

func (s *single) Close() error { return nil }

func TestMerge(t *testing.T) {
	m := Merge(
		Single("a", "Hello, A!"),
		Single("b", "Hello, B!"),
		Single("c", "Hello, C!"),
		Single("d", "Hello, D!"),
	)

	expect := map[string]string{
		"a": "Hello, A!",
		"b": "Hello, B!",
		"c": "Hello, C!",
		"d": "Hello, D!",
	}

	AssertFiles(t, m, expect)
}
