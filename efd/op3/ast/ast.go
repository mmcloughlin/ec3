package ast

import (
	"fmt"
	"strconv"
)

type Program struct {
	Assignments []Assignment
}

type Assignment struct {
	LHS Variable
	RHS Expression
}

func (a Assignment) String() string {
	return fmt.Sprintf("%s = %s", a.LHS, a.RHS)
}

type Expression interface {
	fmt.Stringer
}

type Pow struct {
	X Variable
	N Constant
}

func (p Pow) String() string { return fmt.Sprintf("%s^%s", p.X, p.N) }

type Inv struct{ X Operand }

func (i Inv) String() string { return fmt.Sprintf("1/%s", i.X) }

type Mul struct{ X, Y Operand }

func (m Mul) String() string { return fmt.Sprintf("%s*%s", m.X, m.Y) }

type Neg struct{ X Operand }

func (n Neg) String() string { return fmt.Sprintf("-%s", n.X) }

type Add struct{ X, Y Operand }

func (a Add) String() string { return fmt.Sprintf("%s+%s", a.X, a.Y) }

type Sub struct{ X, Y Operand }

func (s Sub) String() string { return fmt.Sprintf("%s-%s", s.X, s.Y) }

type Operand interface {
	fmt.Stringer
}

type Variable string

func (v Variable) String() string   { return string(v) }
func (v Variable) GoString() string { return fmt.Sprintf("ast.Variable(%q)", v) }

type Constant uint

func (c Constant) String() string   { return strconv.FormatUint(uint64(c), 10) }
func (c Constant) GoString() string { return fmt.Sprintf("ast.Constant(%d)", c) }
