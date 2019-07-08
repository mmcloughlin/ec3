package fp

import (
	"go/token"
	"go/types"

	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Config struct {
	PackageName     string
	Field           fp.Crandall
	ElementTypeName string
}

func (c Config) Type() *types.Named {
	array := types.NewArray(types.Typ[types.Byte], int64(c.Field.ElementSize()))
	name := types.NewTypeName(token.NoPos, nil, c.ElementTypeName, nil)
	return types.NewNamed(name, array, nil)
}

func (c Config) PointerType() *types.Pointer {
	return types.NewPointer(c.Type())
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
	a.Mul()

	if err := fs.CompileAsm(cfg.PackageName, "fp_amd64", a.Context()); err != nil {
		return nil, err
	}

	return fs, nil
}

type api struct {
	Config
	gocode.Generator
}

func API(cfg Config) ([]byte, error) {
	a := &api{
		Config:    cfg,
		Generator: gocode.NewGenerator(),
	}
	return a.Generate()
}

func (a *api) Generate() ([]byte, error) {
	a.CodeGenerationWarning(gen.GeneratedBy)
	a.Package(a.Config.PackageName)

	// Define element type.
	a.NL()
	a.Comment("Size of a field element in bytes.")
	a.Linef("const Size = %d", a.Config.Field.ElementSize())

	a.NL()
	a.Commentf("%s is a field element.", a.Config.ElementTypeName)
	a.Linef("type %s %s", a.Config.Type(), a.Config.Type().Underlying())

	return a.Result()
}
