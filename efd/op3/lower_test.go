package op3

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

func TestIsPrimitive(t *testing.T) {
	p := &ast.Program{
		Assignments: []ast.Assignment{
			{
				LHS: "a",
				RHS: ast.Inv{X: ast.Variable("b")},
			},
			{
				LHS: "c",
				RHS: ast.Add{X: ast.Variable("a"), Y: ast.Variable("b")},
			},
			{
				LHS: "d",
				RHS: ast.Cond{X: ast.Variable("a"), C: ast.Variable("C")},
			},
		},
	}

	if !IsPrimitive(p) {
		t.Fail()
	}

	p.Assignments = append(p.Assignments, ast.Assignment{
		LHS: "d",
		RHS: ast.Pow{X: ast.Variable("c"), N: 3},
	})

	if IsPrimitive(p) {
		t.Fail()
	}
}

func TestIsPrimitiveExpression(t *testing.T) {
	cases := []struct {
		Expr   ast.Expression
		Expect bool
	}{
		{
			Expr:   ast.Inv{X: ast.Variable("x")},
			Expect: true,
		},
		{
			Expr:   ast.Neg{X: ast.Variable("x")},
			Expect: true,
		},
		{
			Expr:   ast.Add{X: ast.Variable("x"), Y: ast.Variable("y")},
			Expect: true,
		},
		{
			Expr:   ast.Sub{X: ast.Variable("x"), Y: ast.Variable("y")},
			Expect: true,
		},
		{
			Expr:   ast.Pow{X: ast.Variable("x"), N: 2},
			Expect: true,
		},
		{
			Expr:   ast.Pow{X: ast.Variable("x"), N: 3},
			Expect: false,
		},
		{
			Expr:   ast.Mul{X: ast.Variable("x"), Y: ast.Variable("y")},
			Expect: true,
		},
		{
			Expr:   ast.Mul{X: ast.Constant(3), Y: ast.Variable("y")},
			Expect: false,
		},
		{
			Expr:   ast.Variable("x"),
			Expect: true,
		},
		{
			Expr:   ast.Constant(3),
			Expect: true,
		},
	}
	for _, c := range cases {
		if got := IsPrimitiveExpression(c.Expr); got != c.Expect {
			t.Errorf("IsPrimitiveExpression(%s) = %v; exepct %v", c.Expr, got, c.Expect)
		}
	}
}

func TestLowerCases(t *testing.T) {
	a := ast.Variable("a")
	b := ast.Variable("b")
	cases := []struct {
		Name       string
		Assignment ast.Assignment
		Expect     []ast.Assignment
	}{
		{
			Name: "mul2",
			Assignment: ast.Assignment{
				LHS: b,
				RHS: ast.Mul{X: ast.Constant(2), Y: a},
			},
			Expect: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}},
			},
		},
		{
			Name: "mul12",
			Assignment: ast.Assignment{
				LHS: b,
				RHS: ast.Mul{X: ast.Constant(12), Y: a},
			},
			Expect: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}}, // 2*a
				{LHS: b, RHS: ast.Add{X: b, Y: a}}, // 3*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 6*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 12*a
			},
		},
		{
			Name: "mul64",
			Assignment: ast.Assignment{
				LHS: b,
				RHS: ast.Mul{X: ast.Constant(64), Y: a},
			},
			Expect: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}}, // 2*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 4*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 8*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 16*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 32*a
				{LHS: b, RHS: ast.Add{X: b, Y: b}}, // 64*a
			},
		},
		{
			Name: "cube",
			Assignment: ast.Assignment{
				LHS: b,
				RHS: ast.Pow{X: a, N: 3},
			},
			Expect: []ast.Assignment{
				{LHS: b, RHS: ast.Pow{X: a, N: 2}}, // a^2
				{LHS: b, RHS: ast.Mul{X: b, Y: a}}, // a^3
			},
		},
	}
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			p := &ast.Program{
				Assignments: []ast.Assignment{c.Assignment},
			}

			got, err := Lower(p)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("got:\n%s", got)

			expect := &ast.Program{Assignments: c.Expect}
			if !reflect.DeepEqual(got, expect) {
				t.Logf("expect:\n%s", expect)
				t.Fail()
			}
		})
	}
}

func TestLowerCorpus(t *testing.T) {
	for id, p := range Corpus() {
		low, err := Lower(p)
		if err != nil {
			t.Fatal(err)
		}

		if !IsPrimitive(low) {
			t.Errorf("%s: lowered program is not primitive", id)
		}
	}
}
