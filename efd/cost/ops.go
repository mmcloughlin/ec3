package cost

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

type Operation interface {
	Code() string
}

type Inv struct{}

func (Inv) Code() string { return "I" }

type Mul struct{}

func (Mul) Code() string { return "M" }

type ConstMul uint

func (c ConstMul) Code() string {
	return fmt.Sprintf("*%d", c)
}

type ParamMul string

func (p ParamMul) Code() string {
	return fmt.Sprintf("*%s", p)
}

type Pow uint

func (p Pow) IsSquare() bool {
	return p == 2
}

func (p Pow) Code() string {
	if p.IsSquare() {
		return "S"
	}
	return fmt.Sprintf("^%d", p)
}

type Add struct{}

func (Add) Code() string { return "add" }

type Count struct {
	N  int
	Op Operation
}

func (c Count) Weight(m Model) float64 {
	return float64(c.N) * m.Weight(c.Op)
}

func (c Count) String() string {
	return fmt.Sprintf("%d%s", c.N, c.Op.Code())
}

type Counts []Count

// Sort by weight of the operation according to the model.
func (c Counts) Sort(m Model) {
	sort.Slice(c, func(i, j int) bool {
		wi := m.Weight(c[i].Op)
		wj := m.Weight(c[j].Op)
		if wi != wj {
			return wi > wj
		}
		return c[i].Op.Code() < c[j].Op.Code()
	})
}

func (c Counts) Weight(m Model) float64 {
	var w float64
	for _, count := range c {
		w += count.Weight(m)
	}
	return w
}

// Summary returns a concise string representation of the operation counts,
// discarding negligible operations.
func (c Counts) Summary() string {
	return c.string(Weights{I: 5, M: 4, S: 3, Pow: 2, ParamM: 1})
}

func (c Counts) String() string {
	return c.string(Weights{I: 7, M: 6, S: 5, Pow: 4, ParamM: 3, Add: 2, ConstM: 1})
}

func (c Counts) string(m Model) string {
	// Clone and sort according to the model.
	components := append(Counts{}, c...)
	components.Sort(m)

	// Empty is considered "0M".
	if len(components) == 0 {
		components = Counts{{Op: Mul{}}}
	}

	// Collect counts until the weight is zero.
	parts := []string{}
	for _, count := range components {
		if m.Weight(count.Op) == 0 {
			break
		}
		parts = append(parts, count.String())
	}

	return strings.Join(parts, " + ")
}

// Operations counts operation types used in the given formula.
func Operations(f *efd.Formula) (Counts, error) {
	if f.Program == nil {
		return nil, xerrors.New("missing program")
	}

	// Operations involving parameters are considered special. Build a set of all
	// such variables so we can identify them.
	isparam := map[ast.Variable]bool{}
	for _, p := range f.AllParameters() {
		isparam[ast.Variable(p)] = true
	}

	counts := map[Operation]int{}
	for _, a := range f.Program.Assignments {
		var op Operation
		var err error
		switch e := a.RHS.(type) {
		case ast.Inv:
			op = Inv{}
		case ast.Mul:
			op, err = mul(e, isparam)
		case ast.Pow:
			op = Pow(e.N)
		case ast.Add, ast.Sub, ast.Neg:
			op = Add{}
		case ast.Variable, ast.Constant:
			continue
		default:
			return nil, errutil.UnexpectedType(e)
		}

		if err != nil {
			return nil, err
		}

		counts[op]++
	}

	result := Counts{}
	for op, count := range counts {
		result = append(result, Count{
			N:  count,
			Op: op,
		})
	}

	return result, nil
}

// mul distinguishes between different types of multiplies.
func mul(m ast.Mul, isparam map[ast.Variable]bool) (Operation, error) {
	// Expect the second operand to always be a variable.
	if _, ok := m.Y.(ast.Variable); !ok {
		return nil, errutil.AssertionFailure("expect second multiply operand to be variable")
	}

	// Check for a const multiply.
	if c, ok := m.X.(ast.Constant); ok {
		return ConstMul(c), nil
	}

	// Check for parameter multiply.
	if v, ok := m.X.(ast.Variable); ok && isparam[v] {
		return ParamMul(v), nil
	}

	if v, ok := m.Y.(ast.Variable); ok && isparam[v] {
		return ParamMul(v), nil
	}

	// Fallback to a generic multiply.
	return Mul{}, nil
}
