package op3

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
)

// Corpus returns a suite of test programs.
func Corpus() map[string]*ast.Program {
	corpus := map[string]*ast.Program{}

	// Include everything from the EFD.
	for _, f := range efd.All {
		p := f.Program
		if p == nil {
			continue
		}
		corpus[f.ID] = p
	}

	// Conditionals do not appear in the EFD. Include a simple program that uses them.
	corpus["cond"] = &ast.Program{
		Assignments: []ast.Assignment{
			{
				LHS: ast.Variable("a"),
				RHS: ast.Cond{X: ast.Variable("b"), C: ast.Variable("c")},
			},
		},
	}

	return corpus
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
