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

func TestSUB(t *testing.T) {
	cases := []struct {
		X, Y, BorrowIn  uint64
		Diff, BorrowOut uint64
	}{
		{X: 5, Y: 2, BorrowIn: 0, Diff: 3, BorrowOut: 0},
		{X: 5, Y: 2, BorrowIn: 1, Diff: 2, BorrowOut: 0},
		{X: 2, Y: 5, BorrowIn: 0, Diff: (1 << 64) - 3, BorrowOut: 1},
		{X: 2, Y: 5, BorrowIn: 1, Diff: (1 << 64) - 4, BorrowOut: 1},
	}
	for _, c := range cases {
		p := &ir.Program{
			Instructions: []ir.Instruction{
				ir.MOV{Source: ir.Constant(c.X), Destination: ir.Register("X")},
				ir.MOV{Source: ir.Constant(c.BorrowIn), Destination: ir.Register("CF")},
				ir.SUB{
					X:         ir.Register("X"),
					Y:         ir.Constant(c.Y),
					BorrowIn:  ir.Register("CF"),
					Diff:      ir.Register("D"),
					BorrowOut: ir.Register("CF"),
				},
			},
		}
		Execute(t, p, []Expectation{
			{ir.Register("X"), c.X},
			{ir.Register("D"), c.Diff},
			{ir.Register("CF"), c.BorrowOut},
		})
	}
}

func TestMUL(t *testing.T) {
	cases := []struct {
		X, Y      uint64
		High, Low uint64
	}{
		{X: 5, Y: 2, High: 0, Low: 10},
		{X: 0xfeedbeefcafe, Y: 0x123456789, High: 0x1220d, Low: 0x5d391f6e1975d3ee},
	}
	for _, c := range cases {
		p := &ir.Program{
			Instructions: []ir.Instruction{
				ir.MOV{Source: ir.Constant(c.X), Destination: ir.Register("X")},
				ir.MUL{
					X:    ir.Register("X"),
					Y:    ir.Constant(c.Y),
					High: ir.Register("H"),
					Low:  ir.Register("L"),
				},
			},
		}
		Execute(t, p, []Expectation{
			{ir.Register("X"), c.X},
			{ir.Register("H"), c.High},
			{ir.Register("L"), c.Low},
		})
	}
}

func TestSHL(t *testing.T) {
	cases := []struct {
		X, Shift uint64
		Result   uint64
	}{
		{X: 5, Shift: 3, Result: 5 << 3},
		{X: 5, Shift: 0, Result: 5},
		{X: 5, Shift: 100, Result: 0},
	}
	for _, c := range cases {
		p := &ir.Program{
			Instructions: []ir.Instruction{
				ir.MOV{Source: ir.Constant(c.X), Destination: ir.Register("X")},
				ir.SHL{
					X:      ir.Register("X"),
					Shift:  ir.Constant(c.Shift),
					Result: ir.Register("S"),
				},
			},
		}
		Execute(t, p, []Expectation{
			{ir.Register("X"), c.X},
			{ir.Register("S"), c.Result},
		})
	}
}

func TestSHR(t *testing.T) {
	cases := []struct {
		X, Shift uint64
		Result   uint64
	}{
		{X: 500, Shift: 3, Result: 500 >> 3},
		{X: 500, Shift: 0, Result: 500},
		{X: 500, Shift: 100, Result: 0},
	}
	for _, c := range cases {
		p := &ir.Program{
			Instructions: []ir.Instruction{
				ir.MOV{Source: ir.Constant(c.X), Destination: ir.Register("X")},
				ir.SHR{
					X:      ir.Register("X"),
					Shift:  ir.Constant(c.Shift),
					Result: ir.Register("S"),
				},
			},
		}
		Execute(t, p, []Expectation{
			{ir.Register("X"), c.X},
			{ir.Register("S"), c.Result},
		})
	}
}
