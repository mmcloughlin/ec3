package eval

import (
	"crypto/rand"
	"math/big"
	"strings"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/efd/op3/parse"
	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/internal/bigint"
)

func TestExecute(t *testing.T) {
	vars := map[string]int64{
		"x": 2,
	}
	m := 17

	// Define a program.
	lines := []string{
		`a = x^3`,   // a = 8
		`b = 1/a`,   // b = 15
		`c = a*b`,   // c = 1
		`d = -b`,    // d = 2
		`e = a + c`, // e = 9
		`f = e-b`,   // f = 11
		`g = b`,     // g = 15
		`h = f ? c`, // h = 11
		`false = 0`,
		`i = 7`,
		`i = h ? false`, // i = 7
	}
	expect := map[string]int64{
		"a": 8,
		"b": 15,
		"c": 1,
		"d": 2,
		"e": 9,
		"f": 11,
		"g": 15,
		"h": 11,
		"i": 7,
	}

	src := strings.Join(lines, "\n") + "\n"
	p, err := parse.String(src)
	assert.NoError(t, err)
	t.Logf("\n%s", p)

	// Build an evaluator.
	e := NewEvaluator(big.NewInt(int64(m)))
	for v, n := range vars {
		if err := e.Initialize(ast.Variable(v), big.NewInt(n)); err != nil {
			t.Fatal(err)
		}
	}

	// Execute the program.
	if err := e.Execute(p); err != nil {
		t.Fatal(err)
	}

	for v, n := range expect {
		got, ok := e.Load(ast.Variable(v))
		if !ok {
			t.Fatalf("expected value for variable %q", v)
		}
		if !bigint.EqualInt64(got, n) {
			t.Errorf("%s = %d; expect %d", v, got, n)
		}
	}
}

func TestEFDPrograms(t *testing.T) {
	m := big.NewInt(5507)
	for _, f := range efd.Select(efd.WithProgram) {
		p := f.Program
		t.Run(f.ID, func(t *testing.T) {
			e := NewEvaluator(m)

			// Assign random values to inputs.
			inputs := op3.Inputs(p)
			for _, v := range inputs {
				r, err := rand.Int(rand.Reader, m)
				assert.NoError(t, err)

				err = e.Initialize(v, r)
				assert.NoError(t, err)
			}

			// Execute the program.
			err := e.Execute(p)
			assert.NoError(t, err)

			// Verify every variable has a value in the expected range.
			for _, v := range op3.Variables(p) {
				x, ok := e.Load(v)
				if !ok {
					t.Fatalf("missing value for variable %q", v)
				}

				if x.Sign() < 0 || x.Cmp(m) >= 0 {
					t.Errorf("variable %q outside expected range", v)
				}
			}
		})
	}
}
