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

//go:generate go run wrap.go -type Bool -input $GOFILE -output zbool.go

// -------------------------------------------------------------
// Common expressions.
// -------------------------------------------------------------
// Z3_ast Z3_API 	Z3_mk_eq (Z3_context c, Z3_ast l, Z3_ast r)
//  	Create an AST node representing l = r. More...
//
// Z3_ast Z3_API 	Z3_mk_distinct (Z3_context c, unsigned num_args, Z3_ast const args[])
//  	Create an AST node representing distinct(args[0], ..., args[num_args-1]). More...
//
// Z3_ast Z3_API 	Z3_mk_ite (Z3_context c, Z3_ast t1, Z3_ast t2, Z3_ast t3)
//  	Create an AST node representing an if-then-else: ite(t1, t2, t3). More...
//

//wrap:doc Not returns not(x).
//wrap:unary Not Z3_mk_not

//wrap:doc Iff returns x iff y.
//wrap:binary Iff Z3_mk_iff

//wrap:doc Implies returns x implies y.
//wrap:binary Implies Z3_mk_implies

//wrap:doc Xor returns x xor y.
//wrap:binary Xor Z3_mk_xor

//wrap:doc And returns the boolean and of all parameters.
//wrap:go And x y...
//wrap:c Z3_mk_and y...

//wrap:doc Or returns the boolean or of all parameters.
//wrap:go Or x y...
//wrap:c Z3_mk_or y...
