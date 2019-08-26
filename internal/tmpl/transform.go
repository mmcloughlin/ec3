package tmpl

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/xerrors"
)

type visitor struct {
	Pre   astutil.ApplyFunc
	Post  astutil.ApplyFunc
	Error func() error
}

func (v visitor) Transform(n ast.Node) (ast.Node, error) {
	node := astutil.Apply(n, v.Pre, v.Post)
	if err := v.Error(); err != nil {
		return nil, err
	}
	return node, nil
}

func SetPackageName(name string) Transform {
	found := false
	return visitor{
		Post: func(c *astutil.Cursor) bool {
			if _, ok := c.Parent().(*ast.File); ok && c.Name() == "Name" {
				c.Replace(&ast.Ident{
					Name: name,
				})
				found = true
				return false
			}
			return true
		},
		Error: func() error {
			if !found {
				return xerrors.New("package declaration not found")
			}
			return nil
		},
	}
}

func replace(name string, sub func(*ast.Ident) ast.Node) Transform {
	n := 0
	return visitor{
		Post: func(c *astutil.Cursor) bool {
			if i, ok := c.Node().(*ast.Ident); ok && i.Name == name {
				c.Replace(sub(i))
				n++
			}
			return true
		},
		Error: func() error {
			if n == 0 {
				return xerrors.New("no replacements")
			}
			return nil
		},
	}
}

func Rename(from, to string) Transform {
	return replace(from, func(i *ast.Ident) ast.Node {
		r := &ast.Ident{}
		*r = *i
		r.Name = to
		return r
	})
}

func DefineLiteral(name string, kind token.Token, value string) Transform {
	return replace(name, func(*ast.Ident) ast.Node {
		return &ast.BasicLit{
			Kind:  kind,
			Value: value,
		}
	})
}

func DefineLiteralf(name string, kind token.Token, format string, args ...interface{}) Transform {
	return DefineLiteral(name, kind, fmt.Sprintf(format, args...))
}

func DefineString(name, value string) Transform {
	return DefineLiteralf(name, token.STRING, "%q", value)
}

func DefineIntDecimal(name string, value int) Transform {
	return DefineLiteralf(name, token.INT, "%d", value)
}

func DefineIntHex(name string, value int) Transform {
	return DefineLiteralf(name, token.INT, "%#x", value)
}
