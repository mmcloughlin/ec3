package fp

import (
	"go/token"
	"go/types"

	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/addchain/acc/pass"
	"github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Config struct {
	Field        fp.Crandall
	InverseChain *ir.Program

	PackageName     string
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
	a.Linef("const Size = %d", a.Field.ElementSize())

	a.NL()
	a.Commentf("%s is a field element.", a.ElementTypeName)
	a.Linef("type %s %s", a.Type(), a.Type().Underlying())

	// Implement field operations.
	a.Square()
	a.Inverse()

	return a.Formatted()
}

// Square generates a square function. This is currently implemented naively
// using multiply.
func (a *api) Square() {
	a.Comment("Sqr computes z = x^2 (mod p).")
	a.Function("Sqr", a.Signature("z", "x"))
	a.Linef("Mul(z, x, x)")
	a.LeaveBlock()
}

func (a *api) Inverse() {
	// Perform temporary variable allocation.
	p := a.InverseChain
	alloc := pass.Allocator{
		Input:  "x",
		Output: "z",
		Format: "&t[%d]",
	}
	if err := alloc.Execute(p); err != nil {
		a.SetError(err)
		return
	}
	n := len(p.Temporaries)

	// Output the function.
	a.Comment("Inv computes z = 1/x (mod p).")
	a.Function("Inv", a.Signature("z", "x"))

	// Allocate the temporaries on the stack.
	a.Linef("var t [%d]%s", n, a.Type())

	for _, inst := range p.Instructions {
		switch op := inst.Op.(type) {
		case ir.Add:
			a.Linef("Mul(%s, %s, %s)", inst.Output, op.X, op.Y)
		case ir.Double:
			a.Linef("Sqr(%s, %s)", inst.Output, op.X)
		case ir.Shift:
			first := 0
			if inst.Output.Identifier != op.X.Identifier {
				a.Linef("Sqr(%s, %s)", inst.Output, op.X)
				first++
			}
			a.Linef("for s := %d; s < %d; s++ {", first, op.S)
			a.Linef("Sqr(%s, %s)", inst.Output, inst.Output)
			a.Linef("}")
		default:
			a.SetError(errutil.UnexpectedType(op))
		}
	}

	a.LeaveBlock()
}
