package fp

import (
	"go/token"
	"go/types"

	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/fp/mont"
	"github.com/mmcloughlin/ec3/gen"
)

type Config struct {
	Field        fp.Field
	InverseChain *ir.Program

	PackageName     string
	ElementTypeName string
}

// Montgomery reports whether this is a montgomery field. Fields implemented this way require encoding and decoding before
func (c Config) Montgomery() bool {
	_, ok := c.Field.(mont.Field)
	return ok
}

func (c Config) Type() *types.Named {
	array := types.NewArray(types.Typ[types.Byte], int64(c.Field.ElementSize()))
	name := types.NewTypeName(token.NoPos, nil, c.ElementTypeName, nil)
	return types.NewNamed(name, array, nil)
}

func (c Config) PointerType() *types.Pointer {
	return types.NewPointer(c.Type())
}

func (c Config) Signature(params ...string) *types.Signature {
	ptr := c.PointerType()
	vars := []*types.Var{}
	for _, param := range params {
		vars = append(vars, types.NewParam(token.NoPos, nil, param, ptr))
	}
	return types.NewSignature(nil, types.NewTuple(vars...), nil, false)
}

func Package(cfg Config) (gen.Files, error) {
	fs := gen.Files{}

	// Exported functions.
	b, err := API(cfg)
	if err != nil {
		return nil, err
	}

	fs.Add("fp.go", b)

	// Assembly backend.
	a := NewAsm(cfg)
	a.Add()
	a.Sub()
	a.Mul()

	if err := fs.CompileAsm(cfg.PackageName, "fp_amd64", a.Context()); err != nil {
		return nil, err
	}

	return fs, nil
}
