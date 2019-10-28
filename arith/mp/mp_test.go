package mp

import (
	"testing"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/verif"
	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/z3"
)

func TestAddInto(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	spec := verif.Add(ctx, 64, 4)

	x, err := spec.Registers("x")
	assert.NoError(t, err)
	y, err := spec.Registers("y")
	assert.NoError(t, err)
	z, err := spec.Registers("z")
	assert.NoError(t, err)

	bld := build.NewContext()
	c := bld.Register("c")
	AddInto(bld, z, x, y, c)

	p, err := bld.Program()
	assert.NoError(t, err)

	result, err := spec.Prove(p)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("expected true")
	}
}

func TestSubInto(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	spec := verif.Sub(ctx, 64, 4)

	x, err := spec.Registers("x")
	assert.NoError(t, err)
	y, err := spec.Registers("y")
	assert.NoError(t, err)
	z, err := spec.Registers("z")
	assert.NoError(t, err)

	bld := build.NewContext()
	b := bld.Register("b")
	SubInto(bld, z, x, y, b)

	p, err := bld.Program()
	assert.NoError(t, err)

	result, err := spec.Prove(p)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("expected true")
	}
}
