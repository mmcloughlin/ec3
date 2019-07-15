package op3

import (
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
	}
	for _, c := range cases {
		if got := IsPrimitiveExpression(c.Expr); got != c.Expect {
			t.Errorf("IsPrimitiveExpression(%s) = %v; exepct %v", c.Expr, got, c.Expect)
		}
	}
}
