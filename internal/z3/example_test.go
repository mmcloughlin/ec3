package z3_test

import (
	"fmt"

	"github.com/mmcloughlin/ec3/internal/z3"
)

// This example disproves that x - 10 ⩽ 0 IFF x ⩽ 10 for (32-bit) machine integers.
func Example_bitVectorInequality() {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	// Initialize solver.
	solver := ctx.Solver()
	defer solver.Close()

	// Bit-vectors of size 32.
	bv := ctx.BVSort(32)

	// Inputs.
	x := bv.Const("x")
	zero := bv.Uint64(0)
	ten := bv.Uint64(10)

	// Build theorem.
	xminus10 := x.Sub(ten)

	c1 := x.SLE(ten)
	c2 := xminus10.SLE(zero)

	thm := c1.Iff(c2)

	// Prove.
	result, err := solver.Prove(thm)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("result:", result)
	// Output:
	// result: false
}
