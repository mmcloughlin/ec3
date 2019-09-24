package fmla

import (
	"fmt"
	"go/types"
	"reflect"
	"strconv"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

// Action specifies the read/write operation of a function parameter.
type Action uint8

// Possible Action types.
const (
	R  Action = 0x1
	W  Action = 0x2
	RW Action = R | W
)

// Contains reports whether a supports all actions in s.
func (a Action) Contains(s Action) bool {
	return (a & s) == s
}

type Variable interface {
	Pointer() string
	Value() string
}

type value string

func (v value) Pointer() string { return "&" + string(v) }
func (v value) Value() string   { return string(v) }

type pointer string

func (p pointer) Pointer() string { return string(p) }
func (p pointer) Value() string   { return "*" + string(p) }

type Parameter interface {
	Name() string
	Action() Action
	Type() types.Type
	Variables() map[ast.Variable]Variable
	AliasSets(Parameter) [][]ast.Variable
}

// ParametersVariables gathers variables for the given list of parameters.
func ParametersVariables(params ...Parameter) map[ast.Variable]Variable {
	variables := map[ast.Variable]Variable{}
	for _, param := range params {
		for name, v := range param.Variables() {
			variables[name] = v
		}
	}
	return variables
}

// ParametersVariableNames gathers all variable names for the given list of parameters.
func ParametersVariableNames(params ...Parameter) []ast.Variable {
	names := []ast.Variable{}
	for _, param := range params {
		for name := range param.Variables() {
			names = append(names, name)
		}
	}
	return names
}

// basic is a parameter representing a single variable.
type basic struct {
	name string
	a    Action
	typ  types.Type
	v    Variable
}

// Value returns a basic parameter passed by value.
func Value(name string, a Action, t types.Type) Parameter {
	return basic{
		name: name,
		a:    a,
		typ:  t,
		v:    value(name),
	}
}

// Pointer represents a parameter passed as a pointer.
func Pointer(name string, a Action, t types.Type) Parameter {
	return basic{
		name: name,
		a:    a,
		typ:  types.NewPointer(t),
		v:    pointer(name),
	}
}

// Condition returns a condition variable parameter.
func Condition(name string, a Action) Parameter {
	return Value(name, a, types.Typ[types.Uint])
}

func (b basic) Name() string { return b.name }

func (b basic) Action() Action { return b.a }

func (b basic) Type() types.Type { return b.typ }

func (b basic) Variables() map[ast.Variable]Variable {
	return map[ast.Variable]Variable{
		ast.Variable(b.name): b.v,
	}
}

func (b basic) AliasSets(p Parameter) [][]ast.Variable {
	// Only the same type can alias.
	other, ok := p.(basic)
	if !ok || !reflect.DeepEqual(b.typ, other.typ) {
		return nil
	}

	// Non-pointers cannot alias.
	if _, ok := b.typ.Underlying().(*types.Pointer); !ok {
		return nil
	}

	return [][]ast.Variable{{ast.Variable(b.name), ast.Variable(other.name)}}
}

type point struct {
	name     string
	a        Action
	repr     Representation
	indicies []int
}

func Point(name string, a Action, r Representation, indicies ...int) Parameter {
	return point{
		name:     name,
		a:        a,
		repr:     r,
		indicies: indicies,
	}
}

func (p point) Name() string { return p.name }

func (p point) Action() Action { return p.a }

func (p point) Type() types.Type {
	return types.NewPointer(p.repr.Type())
}

func (p point) Variables() map[ast.Variable]Variable {
	vars := map[ast.Variable]Variable{}
	for _, idx := range p.indicies {
		for _, coord := range p.repr.Coordinates {
			name := ast.Variable(fmt.Sprintf("%s%d", coord, idx))
			code := fmt.Sprintf("%s.%s", p.Name(), coord)
			vars[name] = value(code)
		}
	}
	return vars
}

func (p point) AliasSets(param Parameter) [][]ast.Variable {
	other, ok := param.(point)
	if !ok || !p.repr.Equals(other.repr) {
		return nil
	}

	aliases := [][]ast.Variable{}
	for _, coord := range p.repr.Coordinates {
		set := []ast.Variable{}
		idxs := append(p.indicies, other.indicies...)
		for _, idx := range idxs {
			v := coord + strconv.Itoa(idx)
			set = append(set, ast.Variable(v))
		}
		aliases = append(aliases, set)
	}

	return aliases
}
