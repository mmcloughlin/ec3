package mont

import (
	"testing"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/arith/verif"
	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/prime"
	"github.com/mmcloughlin/ec3/z3"
)

func TestAdd(t *testing.T) {
	// Initialize context.
	cfg := z3.NewConfig()
	defer cfg.Close()
	ctx := z3.NewContext(cfg)
	defer ctx.Close()

	p := prime.NISTP256
	spec, err := verif.AddMod(ctx, 64, 4, p.Int())
	assert.NoError(t, err)

	f := New(p)

	x, err := spec.Registers("x")
	assert.NoError(t, err)
	y, err := spec.Registers("y")
	assert.NoError(t, err)
	z, err := spec.Registers("z")
	assert.NoError(t, err)

	bld := build.NewContext()
	f.Add(bld, z, x, y)

	prog, err := bld.Program()
	assert.NoError(t, err)
	t.Logf("program:\n%s", prog)

	result, err := spec.Prove(prog)
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
