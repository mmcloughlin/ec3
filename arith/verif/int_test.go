package verif

import (
	"testing"

	"github.com/mmcloughlin/ec3/z3"
)

func TestLimbsRoundTrip(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	// Initialize solver.
	solver := ctx.Solver()
	defer solver.Close()

	// Break a variable x into limbs and then reform.
	sort := ctx.BVSort(256)
	x := sort.Const("x")
	limbs := Limbs(x, 64)
	xr := FromLimbs(limbs)

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
