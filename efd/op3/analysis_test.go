package op3

import (
	"reflect"
	"sort"
	"testing"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

func TestVariables(t *testing.T) {
	a := ast.Variable("a")
	b := ast.Variable("b")
	c := ast.Variable("c")
	p := &ast.Program{
		Assignments: []ast.Assignment{
			{LHS: a, RHS: ast.Add{X: b, Y: c}},
			{LHS: b, RHS: ast.Add{X: a, Y: b}},
		},
	}

	expect := []ast.Variable{a, b, c}
	got := Variables(p)
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	if !reflect.DeepEqual(expect, got) {
		t.FailNow()
	}
}

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

func TestRenameVariablesCorpus(t *testing.T) {
	// Verify that RenameVariables with an empty replacement map is a no-op.
	noop := map[ast.Variable]ast.Variable{}
	for _, p := range Corpus() {
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
