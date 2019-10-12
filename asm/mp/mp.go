package mp

import (
	"math/big"

	"github.com/mmcloughlin/ec3/internal/ints"

	"github.com/mmcloughlin/avo/attr"
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/internal/bigint"
)

// Int represents a multi-precision integer.
type Int []operand.Op

// Limb returns the ith limb of x, or the immediate value 0 if i ⩾ len(x).
func (x Int) Limb(i int) operand.Op {
	if i < len(x) {
		return x[i]
	}
	return operand.U32(0)
}

// Extend returns a new integer with limbs appended to the end.
func (x Int) Extend(limbs ...operand.Op) Int {
	ext := Int{}
	ext = append(ext, x...)
	ext = append(ext, limbs...)
	return ext
}

// NewInt builds an empty integer with k limbs.
func NewInt(k int) Int {
	return make(Int, k)
}

// NewIntZero builds an integer with all k limbs set to the immediate value 0.
func NewIntZero(k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = operand.U32(0)
	}
	return x
}

// NewIntLimb64 builds multi-precision integer with k 64-bit limbs.
func NewIntLimb64(ctx *build.Context, k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = ctx.GP64()
	}
	return x
}

// Imm returns an integer representing the big integer constant c with k 64-bit limbs.
func Imm(c *big.Int, k int) (Int, error) {
	limbs := bigint.Uint64s(c)
	if len(limbs) > k {
		return nil, xerrors.Errorf("constant %d cannot be represented in %d 64-bit limbs", c, k)
	}
	x := NewIntZero(k)
	for i, limb := range limbs {
		x[i] = operand.Imm(limb)
	}
	return x, nil
}

// ImmUint returns an integer representing the unsigned integer constant c with k 64-bit limbs.
func ImmUint(c uint, k int) (Int, error) {
	return Imm(new(big.Int).SetUint64(uint64(c)), k)
}

// NewIntFromMem builds a multi-precision integer referencing the k 64-bit limbs
// at memory address m.
func NewIntFromMem(m operand.Mem, k int) Int {
	x := NewInt(k)
	for i := 0; i < k; i++ {
		x[i] = m.Offset(8 * i)
	}
	return x
}

// AllocLocal allocates an integer with k 64-bit limbs on the stack of the
// currently active function.
func AllocLocal(ctx *build.Context, k int) Int {
	addr := ctx.AllocLocal(8 * k)
	return NewIntFromMem(addr, k)
}

// Param builds a multi-precision integer from a function parameter. The
// parameter is expected to be a pointer to the start of the integer.
func Param(ctx *build.Context, name string, k int) Int {
	base := ctx.Load(ctx.Param(name), ctx.GP64())
	addr := operand.Mem{Base: base}
	return NewIntFromMem(addr, k)
}

// Copy copies x to y with 64-bit move instructions. If x and y are different
// sizes it will copy to the smaller size.
func Copy(ctx *build.Context, y, x Int) {
	for i := 0; i < len(x) && i < len(y); i++ {
		ctx.MOVQ(x[i], y[i])
	}
}

// CopyIntoRegisters will copy x into registers.
func CopyIntoRegisters(ctx *build.Context, x Int) Int {
	r := NewIntLimb64(ctx, len(x))
	Copy(ctx, r, x)
	return r
}

// StaticGlobal returns a multi-precision integer stored in a static global data
// section.
func StaticGlobal(ctx *build.Context, name string, limbs []uint64) Int {
	addr := ctx.StaticGlobal(name)
	ctx.DataAttributes(attr.RODATA | attr.NOPTR)
	for _, limb := range limbs {
		ctx.AppendDatum(operand.U64(limb))
	}
	return NewIntFromMem(addr, len(limbs))
}

// ConditionalMove copies x into y if c ≡ 1.
func ConditionalMove(ctx *build.Context, y, x Int, c operand.Op) {
	ctx.TESTQ(c, c)
	for i := 0; i < len(x) && i < len(y); i++ {
		ctx.CMOVQNE(x[i], y[i])
	}
}

// Mul does a full multiply z = x*y.
func Mul(ctx *build.Context, z, x, y Int) {
	// TODO(mbm): multi-precision multiply is ugly

	acc := make([]operand.Op, len(z))
	zero := ctx.GP64()

	for j := 0; j < len(y); j++ {
		ctx.Commentf("y[%d]", j)
		ctx.MOVQ(y[j], reg.RDX)
		ctx.XORQ(zero, zero) // clears flags
		carryinto := [2]int{-1, -1}
		for i := 0; i < len(x); i++ {
			k := i + j
			ctx.Commentf("x[%d] * y[%d] -> z[%d]", i, j, k)

			// Determine where the results should go.
			var product [2]operand.Op
			var add [2]bool
			for b := 0; b < 2; b++ {
				if acc[k+b] == nil {
					acc[k+b] = ctx.GP64()
					product[b] = acc[k+b]
				} else {
					product[b] = ctx.GP64()
					add[b] = true
				}
			}

			// Do the multiply.
			ctx.MULXQ(x[i], product[0], product[1])

			// Do the adds.
			if add[0] {
				ctx.ADCXQ(product[0], acc[k])
				carryinto[0] = k + 1
			}
			if add[1] {
				ctx.ADOXQ(product[1], acc[k+1])
				carryinto[1] = k + 2
			}
		}

		if carryinto[0] > 0 {
			ctx.ADCXQ(zero, acc[carryinto[0]])
		}
		if carryinto[1] > 0 {
			ctx.ADOXQ(zero, acc[carryinto[1]])
		}

		//
		ctx.MOVQ(acc[j], z[j])
	}

	for j := len(y); j < len(z); j++ {
		ctx.MOVQ(acc[j], z[j])
	}
}

type Accumulator struct {
	x   Int
	max *big.Int
	ctx *build.Context
}

func NewAccumulator(ctx *build.Context) *Accumulator {
	return &Accumulator{
		ctx: ctx,
		max: new(big.Int),
	}
}

// Int returns the value in the accumulator.
func (a *Accumulator) Int() Int { return a.x }

// AddAtMax adds a into the accumulator starting at limb i, assuming that the
// registers have maximum value max.
func (a *Accumulator) AddAtMax(y Int, max *big.Int, i int) {
	zero := operand.U32(0)

	// Update maximum value.
	Y := new(big.Int).Lsh(max, 64*uint(i))
	a.max.Add(a.max, Y)

	// Add y.
	carry := false
	for len(y) > 0 {
		x := a.get(i)
		switch {
		case !carry && x == nil:
			l := a.ctx.GP64()
			a.ctx.MOVQ(y[0], l)
			a.set(i, l)
		case carry && x == nil:
			l := a.ctx.GP64()
			a.ctx.MOVQ(zero, l)
			a.ctx.ADCQ(y[0], l)
			a.set(i, l)
			carry = true
		case !carry && x != nil:
			a.ctx.ADDQ(y[0], x)
			carry = true
		case carry && x != nil:
			a.ctx.ADCQ(y[0], x)
			carry = true
		}
		i++
		y = y[1:]
	}

	// Consider a possible carry off the top.
	if !carry || i >= a.limbs() {
		return
	}

	if x := a.get(i); x == nil {
		l := a.ctx.GP64()
		a.ctx.MOVQ(zero, l)
		a.set(i, l)
	}
	a.ctx.ADCQ(zero, a.get(i))
}

// AddAt adds a into the accumulator starting at limb i.
func (a *Accumulator) AddAt(y Int, i int) {
	bits := uint(64 * len(y))
	a.AddAtMax(y, bigint.Ones(bits), i)
}

// limbs returns the required number of limbs to represent the max possible value.
func (a *Accumulator) limbs() int {
	return ints.NextMultiple(a.max.BitLen(), 64) / 64
}

func (a *Accumulator) get(i int) operand.Op {
	if i < len(a.x) {
		return a.x[i]
	}
	return nil
}

func (a *Accumulator) set(i int, v operand.Op) {
	for len(a.x) <= i {
		a.x = append(a.x, nil)
	}
	a.x[i] = v
}

// Sqr does a full square z = x^2.
func Sqr(ctx *build.Context, z, x Int) {
	k := len(x)
	acc := NewAccumulator(ctx)

	// Max value of product.
	max64 := bigint.Ones(64)
	maxprod := new(big.Int).Mul(max64, max64)

	// Compute all products x_i * x_j for i < j. Due to symmetry these are are the
	// terms of the product that will appear twice.
	for n := 0; n < 2*k; n++ {
		for i := ints.Max(0, n-k+1); 2*i < n && i < k; i++ {
			j := n - i
			ctx.Commentf("x[%d] * x[%d]", i, j)

			// Multiply.
			ctx.MOVQ(x[i], reg.RAX)
			ctx.MULQ(x[j])
			hi, lo := reg.RDX, reg.RAX
			acc.AddAtMax(Int{lo, hi}, maxprod, n)
		}
	}

	// Double the result so far.
	ctx.Comment("*2")
	v := acc.Int()
	acc.AddAt(v[1:], 1)

	// Add in squared terms.
	for i := 0; i < k; i++ {
		ctx.Commentf("x[%d] * x[%d]", i, i)
		ctx.MOVQ(x[i], reg.RAX)
		ctx.MULQ(x[i])
		hi, lo := reg.RDX, reg.RAX
		acc.AddAtMax(Int{lo, hi}, maxprod, 2*i)
	}

	// TODO(mbm): is copy necessary?
	ctx.Comment("Copy to result.")
	Copy(ctx, z, acc.Int())
}
