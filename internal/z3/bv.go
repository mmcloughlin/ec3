package z3

/*
#cgo LDFLAGS: -lz3
#include <stdint.h>
#include <z3.h>
*/
import "C"

type BVSort struct {
	ctx  *Context
	sort C.Z3_sort
}

func (c *Context) BVSort(bits uint) *BVSort {
	return &BVSort{
		ctx:  c,
		sort: C.Z3_mk_bv_sort(c.ctx, C.unsigned(bits)),
	}
}

func (s *BVSort) Uint64(x uint64) *BV {
	return s.wrap(C.Z3_mk_unsigned_int64(s.ctx.ctx, C.uint64_t(x), s.sort))
}

func (s *BVSort) Const(name string) *BV {
	return s.wrap(C.Z3_mk_const(s.ctx.ctx, s.ctx.symbol(name), s.sort))
}

func (s *BVSort) wrap(ast C.Z3_ast) *BV {
	return &BV{
		ctx: s.ctx.ctx,
		ast: ast,
	}
}

type BV struct {
	ctx C.Z3_context
	ast C.Z3_ast
}

//go:generate go run wrap.go -type BV -input $GOFILE -output zbv.go

//wrap:doc Bitwise negation.
//wrap:unary Not Z3_mk_bvnot

//wrap:doc Add returns standard twos complement addition.
//wrap:binary Add Z3_mk_bvadd

//wrap:doc Sub returns standard twos complement subtraction.
//wrap:binary Sub Z3_mk_bvsub

//wrap:binary SLE:Bool Z3_mk_bvsle
