package fp

import (
	"go/token"
	"go/types"

	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/gocode"
	"github.com/mmcloughlin/ec3/internal/ints"
	"github.com/mmcloughlin/ec3/prime"
)

type Config struct {
	PackageName     string
	Prime           prime.Crandall
	ElementTypeName string
}

func (c Config) ElementSize() int {
	return c.ElementBits() / 8
}

func (c Config) ElementBits() int {
	return ints.NextMultiple(c.Prime.Bits(), 8)
}

func (c Config) ElementType() *types.Named {
	array := types.NewArray(types.Typ[types.Byte], int64(c.ElementSize()))
	ptr := types.NewPointer(array)
	name := types.NewTypeName(token.NoPos, nil, c.ElementTypeName, nil)
	return types.NewNamed(name, ptr, nil)
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
	a.Linef("const Size = %d", a.Config.ElementSize())

	a.NL()
	a.Commentf("%s is a field element.", a.Config.ElementTypeName)
	a.Linef("type %s %s", a.Config.ElementTypeName, a.Config.ElementType().Underlying())

	return a.Result()
}
