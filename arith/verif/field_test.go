package verif

import (
	"testing"

	"github.com/mmcloughlin/ec3/internal/z3"
	"github.com/mmcloughlin/ec3/prime"
)

func TestFieldAddCommutative(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	// Initialize solver.
	solver := ctx.Solver()
	defer solver.Close()

	// Construct field.
	f, err := NewField(ctx, 256, prime.NISTP256.Int())
	if err != nil {
		t.Fatal(err)
	}

	// Variables.
	x, y := f.Var("x"), f.Var("y")

	sum := f.Add(x, y)
	alt := f.Add(y, x)

	thm := sum.Eq(alt)

	// Prove.
	result, err := solver.Prove(thm)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("expected true")
	}
}

func TestFieldLimbsRoundTrip(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	// Initialize solver.
	solver := ctx.Solver()
	defer solver.Close()

	// Construct field.
	f, err := NewField(ctx, 256, prime.NISTP256.Int())
	if err != nil {
		t.Fatal(err)
	}

	// Break a variable x into limbs and then reform.
	x := f.Var("x")
	limbs := f.Limbs(x, 64)
	xr := f.FromLimbs(limbs)

	thm := x.Eq(xr)

	// Prove.
	result, err := solver.Prove(thm)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("expected true")
	}
}
