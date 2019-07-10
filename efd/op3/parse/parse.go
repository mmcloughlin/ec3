package parse

import (
	"io"
	"strings"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/efd/op3/parse/internal/parser"
)

//go:generate pigeon -o internal/parser/zparser.go op3.peg

// File parses filename.
func File(filename string) (*ast.Program, error) {
	return cast(parser.ParseFile(filename))
}

// Reader parses the data from r using filename as information in
// error messages.
func Reader(filename string, r io.Reader) (*ast.Program, error) {
	return cast(parser.ParseReader(filename, r))
}

// String parses s.
func String(s string) (*ast.Program, error) {
	return Reader("string", strings.NewReader(s))
}

func cast(i interface{}, err error) (*ast.Program, error) {
	if err != nil {
		return nil, err
	}
	return i.(*ast.Program), nil
}
