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

// Find x and y such that: x ^ y - 103 == x * y.
func Example_bitVectorEquation() {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	// Initialize solver.
	solver := ctx.Solver()
	defer solver.Close()

	// Construct x ^ y - 103 == x * y.
	bv := ctx.BVSort(32)
	x, y := bv.Const("x"), bv.Const("y")
	xor := x.Xor(y)
	c103 := bv.Uint64(103)
	lhs := xor.Sub(c103)
	rhs := x.Mul(y)
	ctr := lhs.Eq(rhs)

	// Add constraint to solver and check.
	solver.Assert(ctr)
	sat, err := solver.Check()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("satisfiable:", sat)

	// Inspect the model.
	m := solver.Model()
	defer m.Close()
	fmt.Printf("model:\n%s", m)
	// Output:
	// satisfiable: true
	// model:
	// y -> #xf4fce4c5
	// x -> #x2a3854c5
}
