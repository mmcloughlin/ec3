package parse

import (
	"reflect"
	"testing"

	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/assert"
)

func TestParseErrors(t *testing.T) {
	src := `A = X*Y
error: not sure how to handle (R*(-t14)-P^3*t15)/2
t5 = 9*C`
	_, err := String(src)
	assert.Error(t, err)
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Source string
		Expect ast.Assignment
	}{
		{
			Source: "w4 = K*t11",
			Expect: ast.Assignment{
				LHS: "w4",
				RHS: ast.Mul{
					X: ast.Variable("K"),
					Y: ast.Variable("t11"),
				},
			},
		},
		{
			Source: "J = A^2",
			Expect: ast.Assignment{
				LHS: "J",
				RHS: ast.Pow{
					X: ast.Variable("A"),
					N: ast.Constant(2),
				},
			},
		},
		{
			Source: "t10 = t7-t6",
			Expect: ast.Assignment{
				LHS: "t10",
				RHS: ast.Sub{
					X: ast.Variable("t7"),
					Y: ast.Variable("t6"),
				},
			},
		},
		{
			Source: "t10 = t9+t8",
			Expect: ast.Assignment{
				LHS: "t10",
				RHS: ast.Add{
					X: ast.Variable("t9"),
					Y: ast.Variable("t8"),
				},
			},
		},
		{
			Source: "Q = 2*t12",
			Expect: ast.Assignment{
				LHS: "Q",
				RHS: ast.Mul{
					X: ast.Constant(2),
					Y: ast.Variable("t12"),
				},
			},
		},
		{
			Source: "w3 = w1",
			Expect: ast.Assignment{
				LHS: ast.Variable("w3"),
				RHS: ast.Variable("w1"),
			},
		},
		{
			Source: "t51 = Z1^3",
			Expect: ast.Assignment{
				LHS: "t51",
				RHS: ast.Pow{
					X: ast.Variable("Z1"),
					N: ast.Constant(3),
				},
			},
		},
		{
			Source: "t17 = 4*t16",
			Expect: ast.Assignment{
				LHS: "t17",
				RHS: ast.Mul{
					X: ast.Constant(4),
					Y: ast.Variable("t16"),
				},
			},
		},
		{
			Source: "t11 = 1/t10",
			Expect: ast.Assignment{
				LHS: "t11",
				RHS: ast.Inv{
					X: ast.Variable("t10"),
				},
			},
		},
		{
			Source: "Y3 = 3*t39",
			Expect: ast.Assignment{
				LHS: "Y3",
				RHS: ast.Mul{
					X: ast.Constant(3),
					Y: ast.Variable("t39"),
				},
			},
		},
		{
			Source: "t10 = 8*Q1",
			Expect: ast.Assignment{
				LHS: "t10",
				RHS: ast.Mul{
					X: ast.Constant(8),
					Y: ast.Variable("Q1"),
				},
			},
		},
		{
			Source: "t31 = X1^4",
			Expect: ast.Assignment{
				LHS: "t31",
				RHS: ast.Pow{
					X: ast.Variable("X1"),
					N: ast.Constant(4),
				},
			},
		},
		{
			Source: "Z3 = 1",
			Expect: ast.Assignment{
				LHS: "Z3",
				RHS: ast.Constant(1),
			},
		},
		{
			Source: "t0 = 1+w2",
			Expect: ast.Assignment{
				LHS: "t0",
				RHS: ast.Add{
					X: ast.Constant(1),
					Y: ast.Variable("w2"),
				},
			},
		},
		{
			Source: "Z3 = 1-t3",
			Expect: ast.Assignment{
				LHS: "Z3",
				RHS: ast.Sub{
					X: ast.Constant(1),
					Y: ast.Variable("t3"),
				},
			},
		},
		{
			Source: "t7 = X3+1",
			Expect: ast.Assignment{
				LHS: "t7",
				RHS: ast.Add{
					X: ast.Variable("X3"),
					Y: ast.Constant(1),
				},
			},
		},
		{
			Source: "t5 = 6*t4",
			Expect: ast.Assignment{
				LHS: "t5",
				RHS: ast.Mul{
					X: ast.Constant(6),
					Y: ast.Variable("t4"),
				},
			},
		},
		{
			Source: "Y3 = -t5",
			Expect: ast.Assignment{
				LHS: "Y3",
				RHS: ast.Neg{
					X: ast.Variable("t5"),
				},
			},
		},
		{
			Source: "t14 = 16*t12",
			Expect: ast.Assignment{
				LHS: "t14",
				RHS: ast.Mul{
					X: ast.Constant(16),
					Y: ast.Variable("t12"),
				},
			},
		},
		{
			Source: "Z3 = A-1",
			Expect: ast.Assignment{
				LHS: "Z3",
				RHS: ast.Sub{
					X: ast.Variable("A"),
					Y: ast.Constant(1),
				},
			},
		},
		{
			Source: "t4 = 12*t3",
			Expect: ast.Assignment{
				LHS: "t4",
				RHS: ast.Mul{
					X: ast.Constant(12),
					Y: ast.Variable("t3"),
				},
			},
		},
		{
			Source: "t7 = 2d*t6",
			Expect: ast.Assignment{
				LHS: "t7",
				RHS: ast.Mul{
					X: ast.Variable("2d"),
					Y: ast.Variable("t6"),
				},
			},
		},
	}

	// Build one source from all the cases.
	src := ""
	expect := &ast.Program{}
	for _, c := range cases {
		src += c.Source + "\n"
		expect.Assignments = append(expect.Assignments, c.Expect)
	}

	t.Logf("src:\n%s", src)

	// Verify.
	got, err := String(src)
	assert.NoError(t, err)
	if !reflect.DeepEqual(got, expect) {
		t.Logf("got:\n%#v", got)
		t.Logf("expect:\n%#v", expect)
		t.Fatal("mismatch")
	}
}
