package op3

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"strconv"
	"testing"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/efd/op3/eval"
	"github.com/mmcloughlin/ec3/efd/op3/parse"
	"github.com/mmcloughlin/ec3/internal/assert"
	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/test"
)

func TestAliasCorrectCases(t *testing.T) {
	tmpfmt := "TMP%d"
	cases := []struct {
		Input   string
		Aliases [][]ast.Variable
		Outputs []ast.Variable
		Expect  string
	}{
		// Excerpt from the complete addition formula that revealed the aliasing
		// problem in the first place. Here X3 is assigned a temporary and we expect
		// an additional assignment to be added at the end.
		{
			Input: "X3 = t1 + t2\nt4 = t4 - X3\nX3 = X1 + Z1\n",
			Aliases: [][]ast.Variable{
				{"X1", "X3"},
			},
			Outputs: []ast.Variable{"X3"},
			Expect:  "TMP0 = t1 + t2\nt4 = t4 - TMP0\nTMP0 = X1 + Z1\nX3 = TMP0\n",
		},
		// As above, but we also make X3 an input variable. This should trigger the
		// case where an extra assignment is needed prior to the program as well.
		{
			Input: "t1 = X3 + t2\nX3 = t1 + t2\nt4 = t4 - X3\nX3 = X1 + Z1\n",
			Aliases: [][]ast.Variable{
				{"X1", "X3"},
			},
			Outputs: []ast.Variable{"X3"},
			Expect:  "TMP0 = X3\nt1 = TMP0 + t2\nTMP0 = t1 + t2\nt4 = t4 - TMP0\nTMP0 = X1 + Z1\nX3 = TMP0\n",
		},
	}
	for _, c := range cases {
		p, err := parse.String(c.Input)
		assert.NoError(t, err)

		expect, err := parse.String(c.Expect)
		assert.NoError(t, err)

		q := AliasCorrect(p, c.Aliases, c.Outputs, IndexedVariables(tmpfmt))
		if !reflect.DeepEqual(q, expect) {
			t.Logf("got:\n%s", q)
			t.Logf("expect:\n%s", expect)
			t.FailNow()
		}
	}
}

// TestAliasCorrectEFDBinaryXYZ is a fairly elaborate test that applies alias
// correction to a subset of formulae in the EFD database. The subset is
// specifically binary formulae taking arguments (X1:Y1:Z1) (X2:Y2:Z2) and
// returning (X3:Y3:Z3). This was chosen since it is a common format, and this
// subset is known to contain formulae which will fail aliased execution.
//
// For each formula we first check whether it already handles aliasing, just for
// informational purposes. Then we apply alias correction and confirm that the
// result is the same as the original, and that the result is robust to
// aliasing.
func TestAliasCorrectEFDBinaryXYZ(t *testing.T) {
	fs := efd.Select(
		WithInputSet(
			"X1", "Y1", "Z1",
			"X2", "Y2", "Z2",
		),
		WithOutputs(
			"X3", "Y3", "Z3",
		),
	)
	for _, f := range fs {
		f := f // scopelint
		t.Run(f.ID, func(t *testing.T) {
			p := f.Program

			// Check if the original formula works when aliased.
			originalok := true
			PN := NonAliasedEvaluator(t, p)
			PL := LeftAliasedEvaluator(t, p)
			if !CheckEqual(t, PN, PL) {
				originalok = false
				t.Logf("mismatch: left-aliased")
			}

			PR := RightAliasedEvaluator(t, p)
			if !CheckEqual(t, PN, PR) {
				originalok = false
				t.Logf("mismatch: right-aliased")
			}

			// Apply alias correction.
			aliases := [][]ast.Variable{
				{"X1", "X2", "X3"},
				{"Y1", "Y2", "Y3"},
				{"Z1", "Z2", "Z3"},
			}
			outputs := []ast.Variable{"X3", "Y3", "Z3"}
			q := AliasCorrect(p, aliases, outputs, Temporaries())

			// This should work when aliased.
			QL := LeftAliasedEvaluator(t, q)
			if !CheckEqual(t, PN, QL) {
				t.Errorf("left-aliased fails")
			}

			QR := RightAliasedEvaluator(t, q)
			if !CheckEqual(t, PN, QR) {
				t.Errorf("right-aliased fails")
			}

			// If the original was fine, we expect that no changes were made.
			if originalok && !reflect.DeepEqual(p, q) {
				t.Logf("original:\n%s", p)
				t.Logf("got:\n%s", q)
				t.Errorf("original was alias-safe; expect no changes")
			}
		})
	}
}

// mod is the modulus to use for testing evaluations.
var mod = big.NewInt(7919)

// Point is a helper point type.
type Point struct {
	X, Y, Z *big.Int
}

// NewPoint allocates a new zero point.
func NewPoint() *Point {
	return &Point{
		X: new(big.Int),
		Y: new(big.Int),
		Z: new(big.Int),
	}
}

// RandPoint generates a random point with values up to the modulus mod.
func RandPoint(t *testing.T) *Point {
	t.Helper()
	p := NewPoint()
	for _, x := range []*big.Int{p.X, p.Y, p.Z} {
		r, err := rand.Int(rand.Reader, mod)
		assert.NoError(t, err)
		x.Set(r)
	}
	return p
}

// Clone the point.
func (p *Point) Clone() *Point {
	return &Point{
		X: bigint.Clone(p.X),
		Y: bigint.Clone(p.Y),
		Z: bigint.Clone(p.Z),
	}
}

// Store the point in the evaluator with the given index. So for index 2 the
// point would be stored in variables X2, Y2, Z3.
func (p *Point) Store(e *eval.Evaluator, i int) {
	e.Store(ast.Variable("X"+strconv.Itoa(i)), p.X)
	e.Store(ast.Variable("Y"+strconv.Itoa(i)), p.Y)
	e.Store(ast.Variable("Z"+strconv.Itoa(i)), p.Z)
}

// Equal reports whether p is equal to q.
func (p *Point) Equal(q *Point) bool {
	return bigint.Equal(p.X, q.X) && bigint.Equal(p.Y, q.Y) && bigint.Equal(p.Z, q.Z)
}

// Binary operation on points.
type Binary func(a, b *Point) *Point

// NonAliasedEvaluator returns a function that evaluates p with non aliased inputs.
func NonAliasedEvaluator(t *testing.T, p *ast.Program) Binary {
	return ProgramEvaluator(t, p, func(a, b *Point) *Point { return NewPoint() })
}

// LeftAliasedEvaluator returns a function that evaluates p with the return value aliased to the left operand.
func LeftAliasedEvaluator(t *testing.T, p *ast.Program) Binary {
	return ProgramEvaluator(t, p, func(a, b *Point) *Point { return a })
}

// RightAliasedEvaluator returns a function that evaluates p with the return value aliased to the right operand.
func RightAliasedEvaluator(t *testing.T, p *ast.Program) Binary {
	return ProgramEvaluator(t, p, func(a, b *Point) *Point { return b })
}

// ProgramEvaluator is a convenience for building the above evaluators. The
// output function returns the point object to write the output to.
func ProgramEvaluator(t *testing.T, p *ast.Program, output Binary) Binary {
	t.Helper()
	return func(a, b *Point) *Point {
		e := eval.NewEvaluator(mod)
		a.Store(e, 1)
		b.Store(e, 2)
		c := output(a, b)
		c.Store(e, 3)
		err := e.Execute(p)
		assert.NoError(t, err)
		return c
	}
}

// CheckEqual runs f and g against random points and checks for equality.
func CheckEqual(t *testing.T, f, g Binary) bool {
	equal := true
	test.Repeat(t, func(t *testing.T) bool {
		a, b := RandPoint(t), RandPoint(t)
		cf := f(a.Clone(), b.Clone())
		cg := g(a.Clone(), b.Clone())
		equal = equal && cf.Equal(cg)
		return equal
	})
	return equal
}
