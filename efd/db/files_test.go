package db

import "testing"

func TestRealArchive(t *testing.T) {
	_, err := Read("../efd.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
}
