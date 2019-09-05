package efd_test

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/db"
	"github.com/mmcloughlin/ec3/internal/assert"
)

// The purpose of these tests is to verify that code generation preserved the
// data as present in the database.

func TestDB(t *testing.T) {
	a, err := db.Archive("efd.tar.gz")
	assert.NoError(t, err)

	add, err := db.Directory("addenda")
	assert.NoError(t, err)

	m := db.Merge(a, add)

	d, err := db.Read(m)
	assert.NoError(t, err)

	for _, f := range efd.All {
		expect := d.Formulae[f.ID]
		if !reflect.DeepEqual(expect, f) {
			t.Fatalf("mismatch formula %q", f.ID)
		}
	}
}
