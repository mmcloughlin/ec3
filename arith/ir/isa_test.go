package ir

import (
	"testing"
)

func TestProgramString(t *testing.T) {
	p := &Program{
		Instructions: []Instruction{
			MOV{Source: Register("X"), Destination: Register("Y")},
			ADD{
				X:        Constant(42),
				Y:        Register("Y"),
				CarryIn:  Flag(1),
				Sum:      Register("S"),
				CarryOut: Register("c"),
			},
		},
	}
	got := p.String()
	expect := "MOV\tX, Y\nADD\t$0x2a, Y, $1, S, c\n"
	if got != expect {
		t.Fatalf("got:\n%sexpect:\n%s", got, expect)
	}
}
