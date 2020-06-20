package mont

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/eval/m64"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/test"
	"github.com/mmcloughlin/ec3/prime"
)

func TestAdd(t *testing.T) {
	p := prime.NISTP256
	field := New(p)
	k := field.Limbs()

	// Build program.
	ctx := build.NewContext()
	X := ctx.Int("X", k)
	Y := ctx.Int("Y", k)
	Z := ctx.Int("Z", k)

	field.Add(ctx, Z, X, Y)

	prog, err := ctx.Program()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("program:\n%s", prog)

	// got: use the evaluator
	f := func(x, y *big.Int) *big.Int {
		e := m64.NewEvaluator()
		e.SetInt(X, x)
		e.SetInt(Y, y)
		if err := e.Execute(prog); err != nil {
			t.Fatal(err)
		}
		z, err := e.Int(Z)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}

	// expect: modular addition function on two big integers
	g := func(x, y *big.Int) *big.Int {
		z := new(big.Int).Add(x, y)
		z.Mod(z, p.Int())
		return z
	}

	// Random trials.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	test.Repeat(t, func(t *testing.T) bool {
		x := new(big.Int).Rand(r, p.Int())
		y := new(big.Int).Rand(r, p.Int())
		got := f(x, y)
		expect := g(x, y)
		if !bigint.Equal(expect, got) {
			t.Fail()
		}
		return true
	})
}

func TestSub(t *testing.T) {
	p := prime.NISTP256
	field := New(p)
	k := field.Limbs()

	// Build program.
	ctx := build.NewContext()
	X := ctx.Int("X", k)
	Y := ctx.Int("Y", k)
	Z := ctx.Int("Z", k)

	field.Sub(ctx, Z, X, Y)

	prog, err := ctx.Program()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("program:\n%s", prog)

	// got: use the evaluator
	f := func(x, y *big.Int) *big.Int {
		e := m64.NewEvaluator()
		e.SetInt(X, x)
		e.SetInt(Y, y)
		if err := e.Execute(prog); err != nil {
			t.Fatal(err)
		}
		z, err := e.Int(Z)
		if err != nil {
			t.Fatal(err)
		}
		return z
	}

	// expect: modular subtraction function on two big integers
	g := func(x, y *big.Int) *big.Int {
		z := new(big.Int).Sub(x, y)
		z.Add(z, p.Int())
		z.Mod(z, p.Int())
		return z
	}

	// Random trials.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	test.Repeat(t, func(t *testing.T) bool {
		x := new(big.Int).Rand(r, p.Int())
		y := new(big.Int).Rand(r, p.Int())
		got := f(x, y)
		expect := g(x, y)
		if !bigint.Equal(expect, got) {
			t.Fail()
		}
		return true
	})
}
