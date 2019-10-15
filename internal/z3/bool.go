package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

type BoolSort struct {
	ctx  *Context
	sort C.Z3_sort
}

func (c *Context) BoolSort() *BoolSort {
	return &BoolSort{
		ctx:  c,
		sort: C.Z3_mk_bool_sort(c.ctx),
	}
}

type Bool struct {
	ctx C.Z3_context
	ast C.Z3_ast
}

func (l *Bool) Not() *Bool {
	return &Bool{
		ctx: l.ctx,
		ast: C.Z3_mk_not(l.ctx, l.ast),
	}
}

func (l *Bool) Iff(r *Bool) *Bool {
	return &Bool{
		ctx: l.ctx,
		ast: C.Z3_mk_iff(l.ctx, l.ast, r.ast),
	}
}
