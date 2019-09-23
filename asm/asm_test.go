package asm

import "testing"

func TestIsRegisterName(t *testing.T) {
	cases := []struct {
		Name   string
		Expect bool
	}{
		{"SP", true},
		{"BX", true},
		{"R13", true},
		{"X3", true},
		{"Y13", true},
		{"Z7", true},
		{"HELLO", false},
	}
	for _, c := range cases {
		if got := IsRegisterName(c.Name); got != c.Expect {
			t.Errorf("IsRegisterName(%q) = %v; expect %v", c.Name, got, c.Expect)
		}
	}
}
