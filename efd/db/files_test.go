package db

import (
	"testing"

	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestRealArchive(t *testing.T) {
	a, err := Archive("../data/efd.tar.gz")
	assert.NoError(t, err)

	_, err = Read(a)
	assert.NoError(t, err)
}
