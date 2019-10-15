package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

type Sort interface {
	sort() C.Z3_sort
}

type sort struct {
	srt C.Z3_sort
}

func (s sort) sort() C.Z3_sort {
	return s.srt
}
