package ast

type Program struct {
	Assignments []Assignment
}

type Assignment struct {
	LHS Variable
	RHS Expression
}

type Expression interface{}

type Pow struct {
	X Variable
	N Constant
}

type Inv struct{ X Operand }

type Mul struct{ X, Y Operand }

type Neg struct{ X Operand }

type Add struct{ X, Y Operand }

type Sub struct{ X, Y Operand }

type Operand interface{}

type Variable string

type Constant uint
