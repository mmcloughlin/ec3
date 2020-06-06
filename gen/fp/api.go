package fp

import (
	"math/big"
	"strings"

	"github.com/mmcloughlin/addchain/acc"
	"github.com/mmcloughlin/addchain/acc/ir"
	"github.com/mmcloughlin/addchain/acc/pass"
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
	a.DefineVar("prime", p)

	// Conversion to/from integer types.
	for _, raw := range []bool{false, true} {
		a.SetInt64(raw)
		a.SetInt(raw)
		a.SetBytes(raw)
		a.Int(raw)
	}

	// Encoding and decoding for montgomery fields.
	if a.Montgomery() {
		a.Decode()
		a.Encode()
	}

	// Implement field operations.
	a.Negate()
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

// DefineVar defines a field element variable set to the specific integer value.
func (a *api) DefineVar(name string, x *big.Int) {
	a.Printf("var %s = %s", a.Name(name), a.Type())
	a.ByteArrayValue(bigint.BytesLittleEndian(x))
}

func rawname(name string, raw bool) string {
	if raw {
		name += "Raw"
	}
	return name
}

func (a *api) rawcomment(raw bool) {
	if raw {
		a.Comment("This raw variant sets the value directly, bypassing any encoding/decoding steps.")
	}
}

// SetInt64 generates a function to construct a field element from an int64.
func (a *api) SetInt64(raw bool) {
	a.Commentf("%s constructs a field element from an integer.", rawname("SetInt64", raw))
	a.rawcomment(raw)
	a.Printf("func (x %s) %s(y int64) %s", a.PointerType(), rawname("SetInt64", raw), a.PointerType())
	a.EnterBlock()
	a.Linef("x.%s(big.NewInt(y))", rawname("SetInt", raw))
	a.Linef("return x")
	a.LeaveBlock()
}

// SetInt64 generates a function to construct a field element from an int64.
func (a *api) SetBytes(raw bool) {
	a.Commentf("%s constructs a field element from bytes in big-endian order.", rawname("SetBytes", raw))
	a.rawcomment(raw)
	a.Printf("func (x %s) %s(b []byte) %s", a.PointerType(), rawname("SetBytes", raw), a.PointerType())
	a.EnterBlock()
	a.Linef("x.%s(new(big.Int).SetBytes(b))", rawname("SetInt", raw))
	a.Linef("return x")
	a.LeaveBlock()
}

// SetInt generates a function to construct a field element from a big integer.
func (a *api) SetInt(raw bool) {
	a.Commentf("%s constructs a field element from a big integer.", rawname("SetInt", raw))
	a.rawcomment(raw)
	a.Printf("func (x %s) %s(y *big.Int) %s", a.PointerType(), rawname("SetInt", raw), a.PointerType())
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

	if !raw && a.Montgomery() {
		a.Comment("Encode into the Montgomery domain.")
		a.Call("Encode", "x", "x")
	}

	a.Linef("return x")
	a.LeaveBlock()
}

// Int generates a function to convert to a big integer.
func (a *api) Int(raw bool) {
	a.Commentf("%s converts to a big integer.", rawname("Int", raw))
	a.rawcomment(raw)
	a.Printf("func (x %s) %s() *big.Int", a.PointerType(), rawname("Int", raw))
	a.EnterBlock()

	if !raw && a.Montgomery() {
		a.Linef("var z %s", a.Type())
		a.Comment("Decode from the Montgomery domain.")
		a.Call("Decode", "&z", "x")
	} else {
		a.Linef("z := *x")
	}

	a.Comment("Endianness swap.")
	a.Linef("for l, r := 0, %s-1; l < r; l, r = l+1, r-1 {", a.Size())
	a.Linef("z[l], z[r] = z[r], z[l]")
	a.Linef("}")

	a.Comment("Build big.Int.")
	a.Linef("return new(big.Int).SetBytes(z[:])")

	a.LeaveBlock()
}

// Decode generates a decode function for Montgomery fields.
func (a *api) Decode() {
	one := a.Name("one")
	a.Commentf("%s is the field element 1.", one)
	a.DefineVar("one", bigint.One())

	a.Commentf("%s decodes from the Montgomery domain.", a.Name("Decode"))
	a.Function(a.Name("Decode"), a.Signature("z", "x"))
	a.Call("Mul", "z", "x", "&"+one)
	a.LeaveBlock()
}

// Encode generates an encode function for Montgomery fields.
func (a *api) Encode() {
	// Define the R² constant.
	r2 := a.Name("r2")
	a.Commentf("%s is the multiplier R^2 for encoding into the Montgomery domain.", "r2")
	R2 := bigint.Pow2(2 * uint(a.Field.ElementBits()))
	R2.Mod(R2, a.Field.Prime())
	a.DefineVar("r2", R2)

	a.Commentf("%s encodes into the Montgomery domain.", a.Name("Encode"))
	a.Function(a.Name("Encode"), a.Signature("z", "x"))
	a.Call("Mul", "z", "x", "&"+r2)
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
