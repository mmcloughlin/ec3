package fp

import (
	"strings"

	"github.com/mmcloughlin/ec3/addchain/acc"
	"github.com/mmcloughlin/ec3/addchain/acc/ir"
	"github.com/mmcloughlin/ec3/addchain/acc/pass"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/bigint"
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

	// Define element type.
	a.NL()
	a.Commentf("%s is the size of a field element in bytes.", a.Size())
	a.Linef("const %s = %d", a.Size(), a.Field.ElementSize())

	a.NL()
	a.Commentf("%s is a field element.", a.ElementTypeName)
	a.Linef("type %s %s", a.Type(), a.Type().Underlying())

	// Variables.
	a.NL()

	p := a.Field.Prime()
	a.Commentf("%s is the field prime modulus as a big integer.", a.Name("p"))
	a.Linef("var %s, _ = new(big.Int).SetString(\"%d\", 10)", a.Name("p"), p)

	a.Commentf("%s is the prime field modulus as a field element.", a.Name("prime"))
	a.Printf("var %s = %s", a.Name("prime"), a.Type())
	a.ByteArrayValue(bigint.BytesLittleEndian(p))

	// Conversion to/from integer types.
	a.SetInt64()
	a.SetInt()
	a.SetBytes()
	a.Int()

	// Encoding and decoding for montgomery fields.
	if a.Montgomery() {
		a.Decode()
		a.Encode()
	}

	// Implement field operations.
	a.Negate()
	a.Square()
	a.Inverse()

	return a.Formatted()
}

// Name of the field element size constant.
func (a *api) Size() string { return a.Name("Size") }

// Call generates a function call.
func (a *api) Call(name string, args ...interface{}) {
	format := a.Name(name) + "(%s" + strings.Repeat(", %s", len(args)-1) + ")"
	a.Linef(format, args...)
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

// SetInt64 generates a function to construct a field element from an int64.
func (a *api) SetBytes() {
	a.Comment("SetBytes constructs a field element from bytes in big-endian order.")
	a.Printf("func (x %s) SetBytes(b []byte) %s", a.PointerType(), a.PointerType())
	a.EnterBlock()
	a.Linef("x.SetInt(new(big.Int).SetBytes(b))")
	a.Linef("return x")
	a.LeaveBlock()
}

// SetInt generates a function to construct a field element from a big integer.
func (a *api) SetInt() {
	a.Comment("SetInt constructs a field element from a big integer.")
	a.Printf("func (x %s) SetInt(y *big.Int) %s", a.PointerType(), a.PointerType())
	a.EnterBlock()

	a.Comment("Reduce if outside range.")
	a.Linef("if y.Sign() < 0 || y.Cmp(%s) >= 0 {", a.Name("p"))
	a.Linef("y = new(big.Int).Mod(y, %s)", a.Name("p"))
	a.Linef("}")

	a.Comment("Copy bytes into field element.")
	a.Linef("b := y.Bytes()")
	a.Linef("i := 0")
	a.Linef("for ; i < len(b); i++ {")
	a.Linef("x[i] = b[len(b)-1-i]")
	a.Linef("}")
	a.Linef("for ; i < %s; i++ {", a.Size())
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
	a.Linef("for i := 0; i < %s; i++ {", a.Size())
	a.Linef("be[%s-1-i] = x[i]", a.Size())
	a.Linef("}")

	a.Comment("Build big.Int.")
	a.Linef("return new(big.Int).SetBytes(be[:])")

	a.LeaveBlock()
}

// Decode generates a decode function for Montgomery fields.
func (a *api) Decode() {
	one := a.Name("one")
	a.Commentf("%s is the field element 1.", one)
	a.Linef("var %s = new(%s).SetInt64(1)", one, a.Type())

	a.Commentf("%s decodes from the Montgomery domain.", a.Name("Decode"))
	a.Function(a.Name("Decode"), a.Signature("z", "x"))
	a.Call("Mul", "z", "x", one)
	a.LeaveBlock()
}

// Encode generates an encode function for Montgomery fields.
func (a *api) Encode() {
	// Define the R^2 constant.
	r2 := a.Name("r2")
	a.Commentf("%s is the multiplier R^2 for encoding into the Montgomery domain.", r2)
	n := a.Field.ElementBits()
	a.Linef("var %s = new(%s).SetInt(new(big.Int).Lsh(big.NewInt(1), 2*%d))", r2, a.Type(), n)

	a.Commentf("%s encodes into the Montgomery domain.", a.Name("Encode"))
	a.Function(a.Name("Encode"), a.Signature("z", "x"))
	a.Call("Mul", "z", "x", r2)
	a.LeaveBlock()
}

// Negate generates a negation function. This is currently implemented using
// subtraction.
func (a *api) Negate() {
	a.Commentf("%s computes z = -x (mod p).", a.Name("Neg"))
	a.Function(a.Name("Neg"), a.Signature("z", "x"))
	a.Call("Sub", "z", "&"+a.Name("prime"), "x")
	a.LeaveBlock()
}

// Square generates a square function. This is currently implemented naively
// using multiply.
func (a *api) Square() {
	a.Commentf("%s computes z = x^2 (mod p).", a.Name("Sqr"))
	a.Function(a.Name("Sqr"), a.Signature("z", "x"))
	a.Call("Mul", "z", "x", "x")
	a.LeaveBlock()
}

func (a *api) Inverse() {
	// Function header.
	a.Commentf("%s computes z = 1/x (mod p).", a.Name("Inv"))
	a.Function(a.Name("Inv"), a.Signature("z", "x"))

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
			a.Call("Mul", inst.Output, op.X, op.Y)
		case ir.Double:
			a.Call("Sqr", inst.Output, op.X)
		case ir.Shift:
			first := 0
			if inst.Output.Identifier != op.X.Identifier {
				a.Call("Sqr", inst.Output, op.X)
				first++
			}
			a.Linef("for s := %d; s < %d; s++ {", first, op.S)
			a.Call("Sqr", inst.Output, inst.Output)
			a.Linef("}")
		default:
			a.SetError(errutil.UnexpectedType(op))
		}
	}

	a.LeaveBlock()
}
