package fmla

import (
	"sort"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/gotypes"
	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/asm"
	asmfp "github.com/mmcloughlin/ec3/asm/fp"
	"github.com/mmcloughlin/ec3/asm/mp"
	"github.com/mmcloughlin/ec3/efd/op3"
	"github.com/mmcloughlin/ec3/efd/op3/ast"
	"github.com/mmcloughlin/ec3/gen/fp"
	"github.com/mmcloughlin/ec3/internal/errutil"
)

type Asm struct {
	cfg   fp.Config
	field asmfp.Builder
	ctx   *build.Context
}

func NewAsm(cfg fp.Config) *Asm {
	ctx := build.NewContext()
	return &Asm{
		cfg:   cfg,
		field: cfg.Field.Build(ctx),
		ctx:   ctx,
	}
}

func (a *Asm) Context() *build.Context {
	return a.ctx
}

func (a *Asm) Function(name string, p *ast.Program, outputs []ast.Variable) error {
	field := a.cfg.Field

	// Declare the function.
	a.ctx.Function(name)
	a.ctx.Pragma("noescape")

	params := []string{}
	for _, output := range outputs {
		params = append(params, paramname(output))
	}
	for _, input := range op3.Inputs(p) {
		params = append(params, paramname(input))
	}
	sort.Strings(params)
	sig := a.cfg.Signature(params...)
	a.ctx.Signature(gotypes.NewSignature(nil, sig))

	// Allocate stack space.
	size := field.ElementSize()
	stack := map[ast.Variable]mp.Int{}
	for _, v := range op3.Variables(p) {
		addr := a.ctx.AllocLocal(size)
		stack[v] = mp.NewIntFromMem(addr, field.Limbs())
	}

	// Copy inputs to stack.
	t := mp.NewIntLimb64(a.ctx, field.Limbs())
	for _, input := range op3.Inputs(p) {
		x := mp.Param(a.ctx, paramname(input), field.Limbs())
		mp.Copy(a.ctx, t, x)
		mp.Copy(a.ctx, stack[input], t)
	}

	// Allocate space for temporary multiplication results.
	m := mp.AllocLocal(a.ctx, 2*field.Limbs())

	// Generate program.
	for step, asgn := range p.Assignments {
		a.ctx.Commentf("Step %d: %s", step+1, asgn.RHS)
		// TODO(mbm): refactor common code in case blocks
		switch e := asgn.RHS.(type) {
		case ast.Variable:
			ops, err := a.operands(stack, asgn.LHS, e)
			if err != nil {
				return err
			}
			mp.Copy(a.ctx, t, ops[1])
			mp.Copy(a.ctx, ops[0], t)
		case ast.Pow:
			if e.N != 2 {
				return xerrors.New("non-square powers are not supported")
			}
			ops, err := a.operands(stack, asgn.LHS, e.X)
			if err != nil {
				return err
			}
			x := mp.CopyIntoRegisters(a.ctx, ops[1])
			mp.Mul(a.ctx, m, x, x)
			a.field.ReduceDouble(ops[0], m)
		case ast.Mul:
			ops, err := a.operands(stack, asgn.LHS, e.X, e.Y)
			if err != nil {
				return err
			}
			x := mp.CopyIntoRegisters(a.ctx, ops[1])
			y := mp.CopyIntoRegisters(a.ctx, ops[2])
			mp.Mul(a.ctx, m, x, y)
			a.field.ReduceDouble(ops[0], m)
		case ast.Sub:
			ops, err := a.operands(stack, asgn.LHS, e.X, e.Y)
			if err != nil {
				return err
			}
			x := mp.CopyIntoRegisters(a.ctx, ops[1])
			y := mp.CopyIntoRegisters(a.ctx, ops[2])
			a.field.Sub(x, y)
			mp.Copy(a.ctx, ops[0], x)
		case ast.Add:
			ops, err := a.operands(stack, asgn.LHS, e.X, e.Y)
			if err != nil {
				return err
			}
			x := mp.CopyIntoRegisters(a.ctx, ops[1])
			y := mp.CopyIntoRegisters(a.ctx, ops[2])
			a.field.Add(x, y)
			mp.Copy(a.ctx, ops[0], x)
		case ast.Inv, ast.Neg, ast.Cond, ast.Constant:
			return xerrors.Errorf("operation %T is not supported in assembly", e)
		default:
			return errutil.UnexpectedType(e)
		}
	}

	// Store outputs.
	for _, output := range outputs {
		x := mp.Param(a.ctx, paramname(output), field.Limbs())
		mp.Copy(a.ctx, t, stack[output])
		mp.Copy(a.ctx, x, t)
	}

	// Return.
	a.ctx.RET()

	return nil
}

func (a *Asm) operands(vars map[ast.Variable]mp.Int, ops ...ast.Operand) ([]mp.Int, error) {
	xs := make([]mp.Int, 0, len(ops))
	for _, op := range ops {
		x, err := a.operand(vars, op)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, nil
}

func (a *Asm) operand(vars map[ast.Variable]mp.Int, o ast.Operand) (mp.Int, error) {
	switch op := o.(type) {
	case ast.Variable:
		x, ok := vars[op]
		if !ok {
			return nil, xerrors.Errorf("unknown variable %q", op)
		}
		return x, nil
	case ast.Constant:
		return mp.ImmUint(uint(op), a.field.Limbs())
	default:
		return nil, errutil.UnexpectedType(op)
	}
}

func paramname(v ast.Variable) string {
	name := string(v)
	for asm.IsRegisterName(name) {
		name += "_"
	}
	return name
}
