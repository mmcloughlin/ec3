package mp

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/eval/m64"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/test"
)

func TestAddInto(t *testing.T) {
	k := 4
	n := uint(64 * k)

	// Build program.
	ctx := build.NewContext()
	X := ctx.Int("X", k)
	Y := ctx.Int("Y", k)
	Z := ctx.Int("Z", k)
	c := ctx.Register("c")
	AddInto(ctx, Z, X, Y, c)

	p, err := ctx.Program()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("program:\n%s", p)

	// got: use the evaluator
	f := func(x, y *big.Int) *big.Int {
		e := m64.NewEvaluator()
		e.SetInt(X, x)
		e.SetInt(Y, y)
		if err := e.Execute(p); err != nil {
			t.Fatal(err)
		}
		z, err := e.Int(Z)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}

	// expect: add function on two big integers
	g := func(x, y *big.Int) *big.Int {
		z := new(big.Int).Add(x, y)
		z.And(z, bigint.Ones(n))
		return z
	}

	// Random trials.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	test.Repeat(t, func(t *testing.T) bool {
		x := bigint.RandBits(r, n)
		y := bigint.RandBits(r, n)
		got := f(x, y)
		expect := g(x, y)
		if !bigint.Equal(expect, got) {
			t.Fail()
		}
		return true
	})
}

func TestSubInto(t *testing.T) {
	k := 4
	n := uint(64 * k)

	// Build program.
	ctx := build.NewContext()
	X := ctx.Int("X", k)
	Y := ctx.Int("Y", k)
	Z := ctx.Int("Z", k)
	b := ctx.Register("b")
	SubInto(ctx, Z, X, Y, b)

	p, err := ctx.Program()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("program:\n%s", p)

	// got: use the evaluator
	f := func(x, y *big.Int) *big.Int {
		e := m64.NewEvaluator()
		e.SetInt(X, x)
		e.SetInt(Y, y)
		if err := e.Execute(p); err != nil {
			t.Fatal(err)
		}
		z, err := e.Int(Z)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}

	// expect: unsigned subtraction function on two big integers
	g := func(x, y *big.Int) *big.Int {
		z := new(big.Int).Sub(x, y)
		z.Add(z, bigint.Pow2(n))
		z.And(z, bigint.Ones(n))
		return z
	}

	// Random trials.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	test.Repeat(t, func(t *testing.T) bool {
		x := bigint.RandBits(r, n)
		y := bigint.RandBits(r, n)
		got := f(x, y)
		expect := g(x, y)
		if !bigint.Equal(expect, got) {
			t.Fail()
		}
		return true
	})
}

func TestMulInto(t *testing.T) {
	k := 4
	n := uint(64 * k)

	// Build program.
	ctx := build.NewContext()
	X := ctx.Int("X", k)
	Y := ctx.Int("Y", k)
	Z := ctx.Int("Z", 2*k)
	MulInto(ctx, Z, X, Y)

	p, err := ctx.Program()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("program:\n%s", p)

	// got: use the evaluator
	f := func(x, y *big.Int) *big.Int {
		e := m64.NewEvaluator()
		e.SetInt(X, x)
		e.SetInt(Y, y)
		if err := e.Execute(p); err != nil {
			t.Fatal(err)
		}
		z, err := e.Int(Z)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}

	// expect: add function on two big integers
	g := func(x, y *big.Int) *big.Int {
		return new(big.Int).Mul(x, y)
	}

	// Random trials.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	test.Repeat(t, func(t *testing.T) bool {
		x := bigint.RandBits(r, n)
		y := bigint.RandBits(r, n)
		got := f(x, y)
		expect := g(x, y)
		if !bigint.Equal(expect, got) {
			t.Fail()
		}
		return true
	})
}
