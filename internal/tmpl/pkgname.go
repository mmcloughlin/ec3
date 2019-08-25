package tmpl

import (
	"go/ast"

	"golang.org/x/xerrors"
)

type pkgname struct {
	name  string
	found bool
}

func SetPackageName(name string) Transform {
	return &pkgname{name: name}
}

func (p *pkgname) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return p
	}
	switch n := node.(type) {
	case *ast.File:
		n.Name.Name = p.name
		p.found = true
		return nil
	}
	return p
}

func (p *pkgname) Transform(n ast.Node) (ast.Node, error) {
	ast.Walk(p, n)
	if !p.found {
		return nil, xerrors.New("package statement not found")
	}
	return n, nil
}
