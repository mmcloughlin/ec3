package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

import (
	"math/big"

	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

// Value is a symbolic value of some sort.
type Value interface {
	val() // sealed
}

type value struct {
	ctx C.Z3_context
	ast C.Z3_ast
}

func (value) val() {}

// newvalue constructs a Value of unknown sort kind.
func newvalue(ctx C.Z3_context, ast C.Z3_ast) (Value, error) {
	v := value{ctx: ctx, ast: ast}
	switch v.sortkind() {
	case C.Z3_BOOL_SORT:
		return &Bool{v}, nil
	case C.Z3_BV_SORT:
		return &BV{v}, nil
	default:
		return nil, errutil.AssertionFailure("unknown sort kind %d", kind)
	}
}

func (v value) sort() C.Z3_sort {
	return C.Z3_get_sort(v.ctx, v.ast)
}

func (v value) sortkind() C.Z3_sort_kind {
	return C.Z3_get_sort_kind(v.ctx, v.sort())
}

func (v value) kind() C.Z3_ast_kind {
	return C.Z3_get_ast_kind(v.ctx, v.ast)
}

func (v value) asInt() (*big.Int, bool) {
	if kind := v.kind(); kind != C.Z3_NUMERAL_AST {
		return nil, false
	}
	s := C.GoString(C.Z3_get_numeral_string(v.ctx, v.ast))
	return bigint.Decimal(s)
}
