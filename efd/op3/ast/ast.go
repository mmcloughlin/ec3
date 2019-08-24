package ast

import (
	"fmt"
	"strconv"
)

type Program struct {
	Assignments []Assignment
}

func (p Program) String() string {
	var s string
	for _, a := range p.Assignments {
		s += a.String() + "\n"
	}
	return s
}

type Assignment struct {
	LHS Variable
	RHS Expression
}

func (a Assignment) String() string {
	return fmt.Sprintf("%s = %s", a.LHS, a.RHS)
}

type Expression interface {
	Inputs() []Operand
	fmt.Stringer
}

type Pow struct {
	X Variable
	N Constant
}

func (p Pow) Inputs() []Operand { return []Operand{p.X} }

func (p Pow) String() string { return fmt.Sprintf("%s^%s", p.X, p.N) }

type Inv struct{ X Operand }

func (i Inv) Inputs() []Operand { return []Operand{i.X} }

func (i Inv) String() string { return fmt.Sprintf("1/%s", i.X) }

type Mul struct{ X, Y Operand }

func (m Mul) Inputs() []Operand { return []Operand{m.X, m.Y} }

func (m Mul) String() string { return fmt.Sprintf("%s*%s", m.X, m.Y) }

type Neg struct{ X Operand }

func (n Neg) Inputs() []Operand { return []Operand{n.X} }

func (n Neg) String() string { return fmt.Sprintf("-%s", n.X) }

type Add struct{ X, Y Operand }

func (a Add) Inputs() []Operand { return []Operand{a.X, a.Y} }

func (a Add) String() string { return fmt.Sprintf("%s+%s", a.X, a.Y) }

type Sub struct{ X, Y Operand }

func (s Sub) Inputs() []Operand { return []Operand{s.X, s.Y} }

func (s Sub) String() string { return fmt.Sprintf("%s-%s", s.X, s.Y) }

type Cond struct{ X, C Variable }

func (c Cond) Inputs() []Operand { return []Operand{c.X, c.C} }

func (c Cond) String() string { return fmt.Sprintf("%s?%s", c.X, c.C) }

type Operand interface {
	fmt.Stringer
}

type Variable string

func Variables(operands []Operand) []Variable {
	vs := make([]Variable, 0, len(operands))
	for _, operand := range operands {
		if v, ok := operand.(Variable); ok {
			vs = append(vs, v)
		}
	}
	return vs
}

func (v Variable) Inputs() []Operand { return []Operand{v} }
func (v Variable) String() string    { return string(v) }
func (v Variable) GoString() string  { return fmt.Sprintf("ast.Variable(%q)", v) }

type Constant uint

func (c Constant) Inputs() []Operand { return []Operand{c} }
func (c Constant) String() string    { return strconv.FormatUint(uint64(c), 10) }
func (c Constant) GoString() string  { return fmt.Sprintf("ast.Constant(%d)", c) }
