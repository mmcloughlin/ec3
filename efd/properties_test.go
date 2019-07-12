package efd

import (
	"testing"

	"github.com/mmcloughlin/ec3/efd/op3"
)

func TestSSA(t *testing.T) {
	for _, f := range All {
		if f.Program == nil {
			continue
		}
		if err := op3.CheckSSA(f.Program); err != nil {
			t.Error(err)
		}
	}
}
