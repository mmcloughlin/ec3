package z3

import "testing"

func TestBitVectorExample1(t *testing.T) {
	// Reference: https://github.com/Z3Prover/z3/blob/77c3f1fb822035ef4de04685cabcaa4da3f0a456/examples/c/test_capi.c#L1407-L1436
	//
	//	/**
	//	   \brief Simple bit-vector example. This example disproves that x - 10 ⩽ 0 IFF x ⩽ 10 for (32-bit) machine integers
	//	*/
	//	void bitvector_example1()
	//	{
	//	    Z3_context ctx = mk_context();
	//	    Z3_solver  s = mk_solver(ctx);
	//	    Z3_sort        bv_sort;
	//	    Z3_ast             x, zero, ten, x_minus_ten, c1, c2, thm;
	//
	//	    printf("\nbitvector_example1\n");
	//	    LOG_MSG("bitvector_example1");
	//
	//
	//	    bv_sort   = Z3_mk_bv_sort(ctx, 32);
	//
	//	    x           = mk_var(ctx, "x", bv_sort);
	//	    zero        = Z3_mk_numeral(ctx, "0", bv_sort);
	//	    ten         = Z3_mk_numeral(ctx, "10", bv_sort);
	//	    x_minus_ten = Z3_mk_bvsub(ctx, x, ten);
	//	    /* bvsle is signed less than or equal to */
	//	    c1          = Z3_mk_bvsle(ctx, x, ten);
	//	    c2          = Z3_mk_bvsle(ctx, x_minus_ten, zero);
	//	    thm         = Z3_mk_iff(ctx, c1, c2);
	//	    printf("disprove: x - 10 ⩽ 0 IFF x ⩽ 10 for (32-bit) machine integers\n");
	//	    prove(ctx, s, thm, false);
	//
	//	    del_solver(ctx, s);
	//	    Z3_del_context(ctx);
	//	}
	//

	cfg := NewConfig()
	defer cfg.Close()

	ctx := NewContext(cfg)
	defer ctx.Close()

	solver := ctx.Solver()
	defer solver.Close()

	bv := ctx.BVSort(32)

	x := bv.Const("x")
	zero := bv.Uint64(0)
	ten := bv.Uint64(10)

	xminus10 := x.Sub(ten)

	c1 := x.SLE(ten)
	c2 := xminus10.SLE(zero)

	thm := c1.Iff(c2)

	result, err := solver.Prove(thm)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("expected theorem to be false")
	}
}
