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

//wrap:doc Not returns bitwise negation of the vector.
//wrap:unary Not Z3_mk_bvnot

//wrap:doc All returns the conjunction of bits in the vector, as a vector of length 1.
//wrap:unary All Z3_mk_bvredand

//wrap:doc Any returns the disjunction of bits in the vector, as a vector of length 1.
//wrap:unary Any Z3_mk_bvredor

//wrap:doc And returns the bitwise and of the input vectors.
//wrap:binary And Z3_mk_bvand

//wrap:doc Or returns the bitwise or of the input vectors.
//wrap:binary Or Z3_mk_bvor

//wrap:doc Xor returns the bitwise exclusive-or of the input vectors.
//wrap:binary Xor Z3_mk_bvxor

//wrap:doc Nand returns the bitwise nand of the input vectors.
//wrap:binary Nand Z3_mk_bvnand

//wrap:doc Nor returns the bitwise nor of the input vectors.
//wrap:binary Nor Z3_mk_bvnor

//wrap:doc Xnor returns the bitwise xnor of the input vectors.
//wrap:binary Xnor Z3_mk_bvxnor

//wrap:doc Neg returns twos complement unary minus.
//wrap:unary Neg Z3_mk_bvneg

//wrap:doc Add returns standard twos complement addition.
//wrap:binary Add Z3_mk_bvadd

//wrap:doc Sub returns standard twos complement subtraction.
//wrap:binary Sub Z3_mk_bvsub

//wrap:doc Mul returns standard twos complement multiplication.
//wrap:binary Mul Z3_mk_bvmul

//wrap:doc Udiv returns unsigned division.
//wrap:binary Udiv Z3_mk_bvudiv

//wrap:doc Sdiv returns twos complement signed division.
//wrap:binary Sdiv Z3_mk_bvsdiv

//wrap:doc Urem returns unsigned remainder.
//wrap:binary Urem Z3_mk_bvurem

//wrap:doc Srem returns twos complement signed remainder.
//wrap:binary Srem Z3_mk_bvsrem

//wrap:doc Smod returns twos complement signed remainder (sign follows divisor).
//wrap:binary Smod Z3_mk_bvsrem

//wrap:doc ULT is unsigned less than.
//wrap:binary ULT:Bool Z3_mk_bvult

//wrap:doc SLT is twos complement signed less than.
//wrap:binary SLT:Bool Z3_mk_bvslt

//wrap:doc ULE is unsigned less than or equal to.
//wrap:binary ULE:Bool Z3_mk_bvule

//wrap:doc SLE is twos complement signed less than or equal to.
//wrap:binary SLE:Bool Z3_mk_bvsle

//wrap:doc UGE is unsigned greater than or equal to.
//wrap:binary UGE:Bool Z3_mk_bvuge

//wrap:doc SGE is twos complement signed greater than or equal to.
//wrap:binary SGE:Bool Z3_mk_bvsge

//wrap:doc UGT is unsigned greater than.
//wrap:binary UGT:Bool Z3_mk_bvugt

//wrap:doc SGT is twos complement signed greater than.
//wrap:binary SGT:Bool Z3_mk_bvsgt

//wrap:doc Concat concatenates the given bit-vectors.
//wrap:binary Concat Z3_mk_concat

//wrap:doc Extract the bits high down to low from a bit-vector of size m to yield a new bit-vector of size n, where n = high - low + 1.
//wrap:go Extract x high:uint low:uint
//wrap:c Z3_mk_extract high:unsigned low:unsigned x

//wrap:doc SignExtend the given bit-vector to the (signed) equivalent bit-vector of size m+i, where m is the size of the given bit-vector.
//wrap:go SignExtend x i:uint
//wrap:c Z3_mk_sign_ext i:unsigned x

//wrap:doc ZeroExtend extends the given bit-vector with zeros to the (unsigned) equivalent bit-vector of size m+i, where m is the size of the given bit-vector.
//wrap:go ZeroExtend x i:uint
//wrap:c Z3_mk_zero_ext i:unsigned x

//wrap:doc Repeat the given bit-vector up length i.
//wrap:go Repeat x i:uint
//wrap:c Z3_mk_repeat i:unsigned x

//wrap:doc Lsh returns x << y.
//wrap:binary Lsh Z3_mk_bvshl

//wrap:doc Rsh returns x >> y.
//wrap:binary Rsh Z3_mk_bvlshr

//wrap:doc Srsh returns the arithmetic right shift of x by y.
//wrap:binary Srsh Z3_mk_bvashr

//wrap:doc RotateLeftConst rotates the bits of x to the left i times.
//wrap:go RotateLeftConst x i:uint
//wrap:c Z3_mk_rotate_left i:unsigned x

//wrap:doc RotateRightConst rotates the bits of x to the right i times.
//wrap:go RotateRightConst x i:uint
//wrap:c Z3_mk_rotate_right i:unsigned x

//wrap:doc RotateLeft rotates the bits of x to the left y times.
//wrap:binary RotateLeft Z3_mk_ext_rotate_left

//wrap:doc RotateRight rotates the bits of x to the right y times.
//wrap:binary RotateRight Z3_mk_ext_rotate_right

//
// Z3_ast Z3_API 	Z3_mk_int2bv (Z3_context c, unsigned n, Z3_ast t1)
//  	Create an n bit bit-vector from the integer argument t1. More...
//
// Z3_ast Z3_API 	Z3_mk_bv2int (Z3_context c, Z3_ast t1, bool is_signed)
//  	Create an integer from the bit-vector argument t1. If is_signed is false, then the bit-vector t1 is treated as unsigned. So the result is non-negative and in the range [0..2^N-1], where N are the number of bits in t1. If is_signed is true, t1 is treated as a signed bit-vector. More...
//
// Z3_ast Z3_API 	Z3_mk_bvadd_no_overflow (Z3_context c, Z3_ast t1, Z3_ast t2, bool is_signed)
//  	Create a predicate that checks that the bit-wise addition of t1 and t2 does not overflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvadd_no_underflow (Z3_context c, Z3_ast t1, Z3_ast t2)
//  	Create a predicate that checks that the bit-wise signed addition of t1 and t2 does not underflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvsub_no_overflow (Z3_context c, Z3_ast t1, Z3_ast t2)
//  	Create a predicate that checks that the bit-wise signed subtraction of t1 and t2 does not overflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvsub_no_underflow (Z3_context c, Z3_ast t1, Z3_ast t2, bool is_signed)
//  	Create a predicate that checks that the bit-wise subtraction of t1 and t2 does not underflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvsdiv_no_overflow (Z3_context c, Z3_ast t1, Z3_ast t2)
//  	Create a predicate that checks that the bit-wise signed division of t1 and t2 does not overflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvneg_no_overflow (Z3_context c, Z3_ast t1)
//  	Check that bit-wise negation does not overflow when t1 is interpreted as a signed bit-vector. More...
//
// Z3_ast Z3_API 	Z3_mk_bvmul_no_overflow (Z3_context c, Z3_ast t1, Z3_ast t2, bool is_signed)
//  	Create a predicate that checks that the bit-wise multiplication of t1 and t2 does not overflow. More...
//
// Z3_ast Z3_API 	Z3_mk_bvmul_no_underflow (Z3_context c, Z3_ast t1, Z3_ast t2)
//  	Create a predicate that checks that the bit-wise signed multiplication of t1 and t2 does not underflow. More...
