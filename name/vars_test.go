package name

import (
	"testing"
)

func TestIndexedVariables(t *testing.T) {
	g := IndexedVariables("x%d")

	// Mark some already used.
	g.MarkUsed("x2", "x4", "x7", "x8", "x9")

	// Generate 10 variables.
	expect := []string{
		"x0",
		"x1",
		// x2
		"x3",
		// x4
		"x5",
		"x6",
		// x7
		// x8
		// x9
		"x10",
		"x11",
	}
	for i := range expect {
		got := g.New()
		if got != expect[i] {
			t.Fail()
		}
	}
}
