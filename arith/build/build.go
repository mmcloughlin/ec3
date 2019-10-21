package build

import (
	"github.com/mmcloughlin/ec3/arith/ir"
)

type Context struct {
	prog *ir.Program
}

func NewContext() *Context {
	return &Context{
		prog: &ir.Program{},
	}
}

func (ctx *Context) Program() *ir.Program {
	return ctx.prog
}

func (ctx *Context) MOV(src ir.Operand, dst ir.Register) {
	ctx.instruction(ir.MOV{
		Source:      src,
		Destination: dst,
	})
}

func (ctx *Context) ADD(x, y, ci ir.Operand, s, co ir.Register) {
	ctx.instruction(ir.ADD{
		X:        x,
		Y:        y,
		CarryIn:  ci,
		Sum:      s,
		CarryOut: co,
	})
}

func (ctx *Context) SUB(x, y, bi ir.Operand, d, bo ir.Register) {
	ctx.instruction(ir.SUB{
		X:         x,
		Y:         y,
		BorrowIn:  bi,
		Diff:      d,
		BorrowOut: bo,
	})
}

func (ctx *Context) MUL(x, y ir.Operand, hi, lo ir.Register) {
	ctx.instruction(ir.MUL{
		X:    x,
		Y:    y,
		High: hi,
		Low:  lo,
	})
}

func (ctx *Context) SHL(x ir.Operand, s ir.Constant, r ir.Register) {
	ctx.instruction(ir.SHL{
		X:      x,
		Shift:  s,
		Result: r,
	})
}

func (ctx *Context) SHR(x ir.Operand, s ir.Constant, r ir.Register) {
	ctx.instruction(ir.SHR{
		X:      x,
		Shift:  s,
		Result: r,
	})
}

func (ctx *Context) instruction(i ir.Instruction) {
	ctx.prog.Instructions = append(ctx.prog.Instructions, i)
}
