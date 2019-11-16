package mp

import (
	"testing"

	"github.com/mmcloughlin/ec3/arith/build"
	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestAdhoc(t *testing.T) {
	ctx := build.NewContext()
	x := ctx.Int("x", 4)
	y := ctx.Int("y", 4)
	z := ctx.Int("z", 8)
	MulInto(ctx, z, x, y)
	p, err := ctx.Program()
	assert.NoError(t, err)
	t.Logf("prog:\n%s", p)
}
