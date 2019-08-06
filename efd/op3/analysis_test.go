package op3

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

func TestRenameVariables(t *testing.T) {
	a := ast.Variable("a")
	b := ast.Variable("b")
	c := ast.Variable("c")
	d := ast.Variable("d")
	p := &ast.Program{
		Assignments: []ast.Assignment{
			{LHS: b, RHS: ast.Add{X: a, Y: a}}, // b = a+a
			{LHS: c, RHS: ast.Add{X: b, Y: b}}, // c = b+b
		},
	}
	replacements := map[ast.Variable]ast.Variable{
		b: d,
	}
	expect := &ast.Program{
		Assignments: []ast.Assignment{
			{LHS: d, RHS: ast.Add{X: a, Y: a}}, // d = a+a
			{LHS: c, RHS: ast.Add{X: d, Y: d}}, // c = d+d
		},
	}
	r := RenameVariables(p, replacements)
	if !reflect.DeepEqual(r, expect) {
		t.Fail()
	}
}

func TestRenameVariablesAllFormulae(t *testing.T) {
	// Verify that RenameVariables with an empty replacement map is a no-op.
	noop := map[ast.Variable]ast.Variable{}
	for _, f := range efd.All {
		p := f.Program
		if p == nil {
			continue
		}
		r := RenameVariables(p, noop)
		if !reflect.DeepEqual(p, r) {
			t.Fatal("expected noop")
		}
	}
}

func TestPare(t *testing.T) {
	a := ast.Variable("a")
	b := ast.Variable("b")
	c := ast.Variable("c")
	cases := []struct {
		Assignments []ast.Assignment
		Outputs     []ast.Variable
		Expect      []ast.Assignment
	}{
		// No outputs.
		{
			Assignments: []ast.Assignment{
				{LHS: a, RHS: b},
			},
			Outputs: []ast.Variable{}, // none
			Expect:  []ast.Assignment{},
		},
		// Trivial case.
		{
			Assignments: []ast.Assignment{
				{LHS: a, RHS: b},
			},
			Outputs: []ast.Variable{a},
			Expect: []ast.Assignment{
				{LHS: a, RHS: b},
			},
		},
		// Dependency chain.
		{
			Assignments: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}}, // b = a+a
				{LHS: c, RHS: ast.Add{X: b, Y: b}}, // c = b+b
			},
			Outputs: []ast.Variable{c},
			Expect: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}},
				{LHS: c, RHS: ast.Add{X: b, Y: b}},
			},
		},
		// Unnecessary instruction.
		{
			Assignments: []ast.Assignment{
				{LHS: b, RHS: ast.Add{X: a, Y: a}}, // b = a+a
				{LHS: c, RHS: ast.Add{X: a, Y: a}}, // c = a+a
			},
			Outputs: []ast.Variable{c},
			Expect: []ast.Assignment{
				{LHS: c, RHS: ast.Add{X: a, Y: a}},
			},
		},
		// Multiple assignment to the output variable.
		{
			Assignments: []ast.Assignment{
				{LHS: c, RHS: ast.Add{X: a, Y: a}}, // c = a+a
				{LHS: c, RHS: ast.Add{X: a, Y: a}}, // c = a+a
				{LHS: c, RHS: ast.Add{X: a, Y: a}}, // c = a+a
				{LHS: c, RHS: ast.Add{X: b, Y: b}}, // c = b+b
				{LHS: c, RHS: ast.Add{X: a, Y: a}}, // c = a+a
			},
			Outputs: []ast.Variable{c},
			Expect: []ast.Assignment{
				{LHS: c, RHS: ast.Add{X: a, Y: a}},
			},
		},
	}
	for _, c := range cases {
		p := &ast.Program{Assignments: c.Assignments}
		expect := &ast.Program{Assignments: c.Expect}
		got, err := Pare(p, c.Outputs)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, expect) {
			t.Logf("got:\n%s", got)
			t.Logf("expect:\n%s", expect)
			t.Fatalf("mismatch")
		}
	}
}

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

func TestLowerAllFormulae(t *testing.T) {
	for _, f := range efd.All {
		p := f.Program
		if p == nil || IsPrimitive(p) {
			continue
		}

		low, err := Lower(p)
		if err != nil {
			t.Fatal(err)
		}

		if !IsPrimitive(low) {
			t.Errorf("%s: lowered program is not primitive", f.ID)
		}
	}
}
