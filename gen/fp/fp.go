package fp

import (
	"go/token"
	"go/types"

	"github.com/mmcloughlin/ec3/asm"
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

func Package(cfg Config) (*gen.Bundle, error) {
	b := gen.NewBundle()

	// Exported functions.
	b.Add("fp.go", API(cfg))

	// Assembly backend.
	a := NewAsm(cfg)
	a.Add()
	a.Mul()

	f, err := asm.Compile(a.Context())
	if err != nil {
		return nil, err
	}
	b.Add("fp_amd64.go", gen.Stubs(cfg.PackageName, f))
	b.Add("fp_amd64.s", gen.Asm(cfg.PackageName, f))

	return b, nil
}

type api struct {
	Config
	gocode.Generator
}

func API(cfg Config) gen.Interface {
	return gen.Formatted(&api{
		Config:    cfg,
		Generator: gocode.NewGenerator(),
	})
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
