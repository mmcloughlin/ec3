package eval

import (
	"math"
	"testing"

	"github.com/mmcloughlin/ec3/arith/ir"
)

type Expectation struct {
	Register ir.Register
	Value    uint64
}

func Execute(t *testing.T, p *ir.Program, expectations []Expectation) {
	e := NewEvaluator()
	if err := e.Execute(p); err != nil {
		t.Fatal(err)
	}
	for _, expect := range expectations {
		got, err := e.Register(expect.Register)
		if err != nil {
			t.Fatal(err)
		}
		if got != expect.Value {
			t.Errorf("%s = %#x; expect %#x", expect.Register, got, expect.Value)
		}
	}
}

func TestADD(t *testing.T) {
	cases := []struct {
		X, Y, CarryIn uint64
		Sum, CarryOut uint64
	}{
		{X: 1, Y: 2, CarryIn: 0, Sum: 3, CarryOut: 0},
		{X: 1, Y: 2, CarryIn: 1, Sum: 4, CarryOut: 0},
		{X: 10, Y: math.MaxUint64, CarryIn: 0, Sum: 9, CarryOut: 1},
		{X: math.MaxUint64, Y: 10, CarryIn: 1, Sum: 10, CarryOut: 1},
	}
	for _, c := range cases {
		p := &ir.Program{
			Instructions: []ir.Instruction{
				ir.MOV{Source: ir.Constant(c.X), Destination: ir.Register("X")},
				ir.MOV{Source: ir.Constant(c.CarryIn), Destination: ir.Register("CF")},
				ir.ADD{
					X:        ir.Register("X"),
					Y:        ir.Constant(c.Y),
					CarryIn:  ir.Register("CF"),
					Sum:      ir.Register("S"),
					CarryOut: ir.Register("CF"),
				},
			},
		}
		Execute(t, p, []Expectation{
			{ir.Register("X"), c.X},
			{ir.Register("S"), c.Sum},
			{ir.Register("CF"), c.CarryOut},
		})
	}
}
