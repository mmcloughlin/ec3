package build

import (
	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/name"
)

type Context struct {
	prog *ir.Program
	regs name.UniqueGenerator
	errs errutil.Errors
}

func NewContext() *Context {
	return &Context{
		prog: &ir.Program{},
		regs: name.NewUniqueGenerator(),
	}
}

func (ctx *Context) Program() (*ir.Program, error) {
	return ctx.prog, ctx.errs.Err()
}

// RegisterFromSequence returns a unique register from the sequence s.
func (ctx *Context) RegisterFromSequence(s name.Sequence) ir.Register {
	return ir.Register(ctx.regs.New(s))
}

// Register returns a unique register in the given namespace.
func (ctx *Context) Register(namespace string) ir.Register {
	return ctx.RegisterFromSequence(name.Indexed(namespace + "%d"))
}

func (ctx *Context) Int(namespace string, k int) ir.Registers {
	x := make(ir.Registers, k)
	for i := 0; i < k; i++ {
		x[i] = ctx.Register(namespace)
	}
	return x
}

func (ctx *Context) MOV(src ir.Operand, dst ir.Register) {
	ctx.instruction(ir.MOV{
		Source:      src,
		Destination: dst,
	})
}

func (ctx *Context) ADD(x, y, ci, s, co ir.Operand) {
	ctx.instruction(ir.ADD{
		X:        x,
		Y:        y,
		CarryIn:  ci,
		Sum:      ctx.reg(s),
		CarryOut: ctx.reg(co),
	})
}

func (ctx *Context) SUB(x, y, bi, d, bo ir.Operand) {
	ctx.instruction(ir.SUB{
		X:         x,
		Y:         y,
		BorrowIn:  bi,
		Diff:      ctx.reg(d),
		BorrowOut: ctx.reg(bo),
	})
}

func (ctx *Context) MUL(x, y, hi, lo ir.Operand) {
	ctx.instruction(ir.MUL{
		X:    x,
		Y:    y,
		High: ctx.reg(hi),
		Low:  ctx.reg(lo),
	})
}

func (ctx *Context) SHL(x ir.Operand, s ir.Constant, r ir.Operand) {
	ctx.instruction(ir.SHL{
		X:      x,
		Shift:  s,
		Result: ctx.reg(r),
	})
}

func (ctx *Context) SHR(x ir.Operand, s ir.Constant, r ir.Operand) {
	ctx.instruction(ir.SHR{
		X:      x,
		Shift:  s,
		Result: ctx.reg(r),
	})
}

func (ctx *Context) reg(op ir.Operand) ir.Register {
	if r, ok := op.(ir.Register); ok {
		return r
	}
	ctx.errs.Addf("operand %v is not a register", op)
	return ir.Register("")
}

func (ctx *Context) instruction(i ir.Instruction) {
	// Check register naming.
	for _, reg := range ir.SelectRegisters(i.Operands()) {
		if !name.IsExported(string(reg)) && !ctx.regs.Used(string(reg)) {
			ctx.errs.Addf("unexported names must be managed by the build context: %q is unknown", reg)
			return
		}
	}

	// Add instruction.
	ctx.prog.Instructions = append(ctx.prog.Instructions, i)
}
