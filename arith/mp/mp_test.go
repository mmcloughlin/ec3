package mp

import (
	"testing"
	"time"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/verif"
	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/z3"
)

// TODO(mbm): reduce verification testing boilerplate

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

func TestMul(t *testing.T) {
	// Z3 configuration.
	cfg := z3.NewConfig()
	defer cfg.Close()

	cfg.SetAutoConfig(true)
	cfg.SetModel(true)
	cfg.SetModelValidate(true)
	cfg.SetTrace(true)
	cfg.SetTraceFilename("/tmp/mul4.log")
	cfg.SetTimeout(time.Minute)

	// Initialize context.
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	ctx.SetASTPrintMode(z3.PrintModeSMTLIB2Compliant)

	spec := verif.MulFull(ctx, 4, 2)

	x, err := spec.Registers("x")
	assert.NoError(t, err)
	y, err := spec.Registers("y")
	assert.NoError(t, err)
	z, err := spec.Registers("z")
	assert.NoError(t, err)

	bld := build.NewContext()
	MulInto(bld, z, x, y)

	p, err := bld.Program()
	assert.NoError(t, err)
	t.Logf("program:\n%s", p)

	result, err := spec.Prove(p)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		assignments, err := spec.Counterexample()
		assert.NoError(t, err)
		for name, x := range assignments {
			t.Logf("%s = %#x", name, x)
		}
		t.Fatal("expected true")
	}
}
