package ec

import (
	"go/types"
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

func TestAliasSets(t *testing.T) {
	jacobian := Representation{
		Name:        "Jacobian",
		ElementType: types.Typ[types.Uint64],
		Coordinates: []string{"X", "Y", "Z"},
	}
	cases := []struct {
		X, Y   Parameter
		Expect [][]ast.Variable
	}{
		// Different types do not alias.
		{
			X:      Value("x", W, types.Typ[types.Uint16]),
			Y:      Value("y", W, types.Typ[types.Int32]),
			Expect: nil,
		},
		// Non-pointer types do not alias.
		{
			X:      Condition("c", W),
			Y:      Condition("d", W),
			Expect: nil,
		},
		// Pointer types can alias.
		{
			X: Pointer("x", W, types.Typ[types.Uint8]),
			Y: Pointer("y", W, types.Typ[types.Uint8]),
			Expect: [][]ast.Variable{
				{"x", "y"},
			},
		},
		// Two point types.
		{
			X: Point("p", W, jacobian, 1),
			Y: Point("q", W, jacobian, 2),
			Expect: [][]ast.Variable{
				{"X1", "X2"},
				{"Y1", "Y2"},
				{"Z1", "Z2"},
			},
		},
	}
	for _, c := range cases {
		got := c.X.AliasSets(c.Y)
		if !reflect.DeepEqual(c.Expect, got) {
			t.Errorf("X=%v Y=%v: got %v; expect %v", c.X, c.Y, got, c.Expect)
		}
	}
}
