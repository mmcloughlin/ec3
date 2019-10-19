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
	sort := ctx.BVSort(256)
	f, err := NewField(sort, prime.NISTP256.Int())
	if err != nil {
		t.Fatal(err)
	}

	// Variables.
	x, y := f.Const("x"), f.Const("y")

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
