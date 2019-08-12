package fp

import (
	"strings"

	"github.com/mmcloughlin/ec3/addchain/acc"
	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/addchain/acc/pass"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

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

	a.Import("math/big")

	// Variables.
	a.NL()
	a.Commentf("modulus is the field prime modulus.")
	a.Linef("var modulus, _ = new(big.Int).SetString(\"%d\", 10)", a.Field.Prime())

	// Define element type.
	a.NL()
	a.Comment("Size of a field element in bytes.")
	a.Linef("const Size = %d", a.Field.ElementSize())

	a.NL()
	a.Commentf("%s is a field element.", a.ElementTypeName)
	a.Linef("type %s %s", a.Type(), a.Type().Underlying())

	// Conversion to/from integer types.
	a.SetInt64()
	a.SetInt()
	a.Int()

	// Encoding and decoding for montgomery fields.
	if a.Montgomery() {
		a.Decode()
		a.Encode()
	}

	// Implement field operations.
	a.Square()
	a.Inverse()

	return a.Formatted()
}

// SetInt64 generates a function to construct a field element from an int64.
func (a *api) SetInt64() {
	a.Comment("SetInt64 constructs a field element from an integer.")
	a.Printf("func (x %s) SetInt64(y int64) %s", a.PointerType(), a.PointerType())
	a.EnterBlock()
	a.Linef("x.SetInt(big.NewInt(y))")
	a.Linef("return x")
	a.LeaveBlock()
}

// SetInt generates a function to construct a field element from a big integer.
func (a *api) SetInt() {
	a.Comment("SetInt constructs a field element from a big integer.")
	a.Printf("func (x %s) SetInt(y *big.Int) %s", a.PointerType(), a.PointerType())
	a.EnterBlock()

	a.Comment("Reduce if outside range.")
	a.Linef("if y.Sign() < 0 || y.Cmp(modulus) >= 0 {")
	a.Linef("y = new(big.Int).Mod(y, modulus)")
	a.Linef("}")

	a.Comment("Copy bytes into field element.")
	a.Linef("b := y.Bytes()")
	a.Linef("i := 0")
	a.Linef("for ; i < len(b); i++ {")
	a.Linef("x[i] = b[len(b)-1-i]")
	a.Linef("}")
	a.Linef("for ; i < Size; i++ {")
	a.Linef("x[i] = 0")
	a.Linef("}")

	a.Linef("return x")
	a.LeaveBlock()
}

// Int generates a function to convert to a big integer.
func (a *api) Int() {
	a.Comment("Int converts to a big integer.")
	a.Printf("func (x %s) Int() *big.Int", a.PointerType())
	a.EnterBlock()

	a.Comment("Endianness swap.")
	a.Linef("var be %s", a.Type())
	a.Linef("for i := 0; i < Size; i++ {")
	a.Linef("be[Size-1-i] = x[i]")
	a.Linef("}")

	a.Comment("Build big.Int.")
	a.Linef("return new(big.Int).SetBytes(be[:])")

	a.LeaveBlock()
}

// Decode generates a decode function for Montgomery fields.
func (a *api) Decode() {
	a.Comment("one is the field element 1.")
	a.Linef("var one = new(%s).SetInt64(1)", a.Type())

	a.Comment("Decode from the Montgomery domain.")
	a.Function("Decode", a.Signature("z", "x"))
	a.Linef("Mul(z, x, one)")
	a.LeaveBlock()
}

// Encode generates an encode function for Montgomery fields.
func (a *api) Encode() {
	// Define the R^2 constant.
	a.Comment("r2 is the multiplier R^2 for encoding into the Montgomery domain.")
	n := a.Field.ElementBits()
	a.Linef("var r2 = new(%s).SetInt(new(big.Int).Lsh(big.NewInt(1), 2*%d))", a.Type(), n)

	a.Comment("Encode into the Montgomery domain.")
	a.Function("Encode", a.Signature("z", "x"))
	a.Linef("Mul(z, x, r2)")
	a.LeaveBlock()
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
	// Function header.
	a.Comment("Inv computes z = 1/x (mod p).")
	a.Function("Inv", a.Signature("z", "x"))

	// Comment describing the addition chain.
	p := a.InverseChain.Clone()
	script, err := acc.String(p)
	if err != nil {
		a.SetError(err)
		return
	}
	a.Comment("Inversion computation is derived from the addition chain:", "")
	a.Comment(strings.Split(script, "\n")...)

	if err := pass.Eval(p); err != nil {
		a.SetError(err)
		return
	}

	sqrs, muls := p.Program.Count()
	a.Commentf("Operations: %d squares %d multiplies", sqrs, muls)

	// Perform temporary variable allocation.
	alloc := pass.Allocator{
		Input:  "x",
		Output: "z",
		Format: "&t[%d]",
	}
	if err := alloc.Execute(p); err != nil {
		a.SetError(err)
		return
	}

	// Allocate required temporaries.
	n := len(p.Temporaries)
	a.NL()
	a.Commentf("Allocate %d temporaries.", n)
	a.Linef("var t [%d]%s", n, a.Type())

	for _, inst := range p.Instructions {
		a.NL()
		a.Commentf("Step %d: %s = x^%#x.", inst.Output.Index, inst.Output, p.Chain[inst.Output.Index])
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
