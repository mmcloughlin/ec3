package ir

import "testing"

func TestIntInterface(t *testing.T) {
	var _ Int = Operands{}
	var _ Int = Registers{}
	var _ Int = Constants{}
}
