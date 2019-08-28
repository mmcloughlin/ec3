// Package tmpl provides code generation from pure Go templates.
package tmpl

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"

	"golang.org/x/xerrors"
)

var (
	FileSystemLoader = LoaderFunc(ioutil.ReadFile)
	Default          = Environment{Loader: FileSystemLoader}
)

// Loader is a method of loading templates.
type Loader interface {
	Load(name string) ([]byte, error)
}

// LoaderFunc adapts a function to the Loader interface.
type LoaderFunc func(name string) ([]byte, error)

// Load calls l.
func (l LoaderFunc) Load(name string) ([]byte, error) {
	return l(name)
}

type Environment struct {
	Loader Loader
}

func (e Environment) Load(name string) (*Template, error) {
	src, err := e.Loader.Load(name)
	if err != nil {
		return nil, err
	}
	return ParseFile(name, src)
}

type Transform interface {
	Transform(ast.Node) (ast.Node, error)
}

type Template struct {
	fset *token.FileSet
	node ast.Node
}

func ParseFile(filename string, src []byte) (*Template, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse template: %w", err)
	}
	return &Template{
		fset: fset,
		node: f,
	}, nil
}

func (t *Template) Apply(transforms ...Transform) error {
	for _, transform := range transforms {
		n, err := transform.Transform(t.node)
		if err != nil {
			return err
		}
		t.node = n
	}
	return nil
}

func (t *Template) Node() ast.Node {
	return t.node
}

func (t *Template) Format(w io.Writer) error {
	return format.Node(w, t.fset, t.node)
}

func (t *Template) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Format(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
