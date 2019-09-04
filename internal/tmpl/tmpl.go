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
	"path/filepath"

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

// NewBasePath restricts the given loader to a given path.
func NewBasePath(loader Loader, path string) Loader {
	return LoaderFunc(func(name string) ([]byte, error) {
		return loader.Load(filepath.Join(path, name))
	})
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

func (e Environment) Package(names ...string) (*Package, error) {
	pkg := NewPackage()
	for _, name := range names {
		src, err := e.Loader.Load(name)
		if err != nil {
			return nil, err
		}

		if err := pkg.AddFile(name, src); err != nil {
			return nil, err
		}
	}
	return pkg, nil
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
	f, err := parse(fset, filename, src)
	if err != nil {
		return nil, err
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

type Package struct {
	fset  *token.FileSet
	files map[string]*ast.File
}

func NewPackage() *Package {
	return &Package{
		fset:  token.NewFileSet(),
		files: make(map[string]*ast.File),
	}
}

func (p *Package) AddFile(name string, src []byte) error {
	if _, ok := p.files[name]; ok {
		return xerrors.Errorf("file named %q already exists", name)
	}

	f, err := parse(p.fset, name, src)
	if err != nil {
		return err
	}

	p.files[name] = f
	return nil
}

// Template returns a template for the named file.
func (p *Package) Template(name string) (*Template, bool) {
	f, ok := p.files[name]
	if !ok {
		return nil, false
	}
	return &Template{
		fset: p.fset,
		node: f,
	}, true
}

func (p *Package) Templates() map[string]*Template {
	tpls := map[string]*Template{}
	for name, f := range p.files {
		tpls[name] = &Template{
			fset: p.fset,
			node: f,
		}
	}
	return tpls
}

// Apply transforms to an entire package.
func (p *Package) Apply(transforms ...Transform) error {
	tpl := &Template{
		fset: p.fset,
		node: &ast.Package{Files: p.files},
	}
	return tpl.Apply(transforms...)
}

// parse a template.
func parse(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse template: %w", err)
	}
	return f, nil
}
