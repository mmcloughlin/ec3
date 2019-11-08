package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// Bool is a boolean value.
type Bool struct {
	value
}

//go:generate go run wrap.go -type Bool -input $GOFILE -output zbool.go

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
