// Package generic generates pure Go code from arithmetic programs.
package generic

import (
	"fmt"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/mmcloughlin/ec3/arith/ir"
	"github.com/mmcloughlin/ec3/gen"
	"github.com/mmcloughlin/ec3/internal/errutil"
	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Config struct {
	PackageName     string
	ElementSize     int
	ElementTypeName string
}

func (c Config) Type() *types.Named {
	array := types.NewArray(types.Typ[types.Byte], int64(c.ElementSize))
	name := types.NewTypeName(token.NoPos, nil, c.ElementTypeName, nil)
	return types.NewNamed(name, array, nil)
}

func (c Config) PointerType() *types.Pointer {
	return types.NewPointer(c.Type())
}

//func (c Config) Param(name string) *types.Var {
//	return types.NewParam(token.NoPos, nil, name, c.PointerType())
//}
//
//func (c Config) Params(params ...string) []*types.Var {
//	vars := []*types.Var{}
//	for _, param := range params {
//		vars = append(vars, c.Param(param))
//	}
//	return vars
//}
//
//func (c Config) Signature(params ...string) *types.Signature {
//	tuple := types.NewTuple(c.Params(params...)...)
//	return types.NewSignature(nil, tuple, nil, false)
//}

func (c Config) Signature(s *ir.Signature) (*types.Signature, error) {
	var vs []*types.Var
	for _, param := range s.Vars() {
		switch t := param.Type.(type) {
		case ir.Integer:
			v := types.NewVar(token.NoPos, nil, param.Name, c.PointerType())
			vs = append(vs, v)
		default:
			return nil, errutil.UnexpectedType(t)
		}
	}
	return types.NewSignature(nil, types.NewTuple(vs...), nil, false), nil
}

// TODO(mbm): build tags

// Print generic Go code for the given arithmetic module.
func Print(cfg Config, mod *ir.Module) ([]byte, error) {
	p := newprinter(cfg)
	return p.Generate(mod)
}

type printer struct {
	Config
	gocode.Generator

	declared map[ir.Register]bool
}

func newprinter(cfg Config) *printer {
	return &printer{
		Config:    cfg,
		Generator: gocode.NewGenerator(),
	}
}

func (p *printer) Generate(mod *ir.Module) ([]byte, error) {
	p.CodeGenerationWarning(gen.GeneratedBy)
	p.Package(p.PackageName)
	p.Import("math/bits")

	for _, s := range mod.Sections {
		p.section(s)
	}

	return p.Formatted()
}

func (p *printer) section(sec ir.Section) {
	switch s := sec.(type) {
	case *ir.Function:
		p.function(s)
	default:
		p.SetError(errutil.UnexpectedType(s))
	}
}

func (p *printer) function(f *ir.Function) {
	s, err := p.Signature(f.Signature)
	if err != nil {
		p.SetError(err)
		return
	}

	p.Function(f.Name, s)
	p.enterscope()
	for _, i := range f.Instructions {
		p.instruction(i)
	}
	p.LeaveBlock()
}

func (p *printer) instruction(inst ir.Instruction) {
	switch i := inst.(type) {
	case ir.MOV:
		p.assign(i.Destination)
		p.Linef("%s", p.operand(i.Source))
	case ir.CMOV:
		mask := fmt.Sprintf("-uint64(%s & 1)", p.operand(i.Flag))
		a, b := p.operand(i.Source), p.operand(i.Destination)
		if i.Equals == 0 {
			a, b = b, a
		}
		p.assign(i.Destination)
		p.Linef("(%s & %s) | (%s &^ %s)", a, mask, b, mask)
	case ir.ADD:
		p.assign(i.Sum, i.CarryOut)
		p.Linef("bits.Add64(%s, %s, %s)", p.operand(i.X), p.operand(i.Y), p.operand(i.CarryIn))
	case ir.SUB:
		p.assign(i.Diff, i.BorrowOut)
		p.Linef("bits.Sub64(%s, %s, %s)", p.operand(i.X), p.operand(i.Y), p.operand(i.BorrowIn))
	case ir.MUL:
		p.assign(i.High, i.Low)
		p.Linef("bits.Mul64(%s, %s)", p.operand(i.X), p.operand(i.Y))
	case ir.SHL:
		p.assign(i.Result)
		p.Linef("%s << %d", i.X, i.Shift)
	case ir.SHR:
		p.assign(i.Result)
		p.Linef("%s >> %d", i.X, i.Shift)
	default:
		p.SetError(errutil.UnexpectedType(i))
	}
}

func (p *printer) operand(op ir.Operand) string {
	switch o := op.(type) {
	case ir.Register:
		return string(o)
	case ir.Constant:
		return fmt.Sprintf("%#x", o)
	case ir.Flag:
		return strconv.FormatUint(uint64(o), 10)
	default:
		p.SetError(errutil.UnexpectedType(o))
		return "<unknown>"
	}
}

func (p *printer) enterscope() {
	p.declared = map[ir.Register]bool{}
}

func (p *printer) assign(regs ...ir.Register) {
	// Determine whether all registers have been declared.
	declared := true
	for _, reg := range regs {
		declared = declared && p.declared[reg]
	}

	// Choose correct assignment operator.
	operator := ":="
	if declared {
		operator = "="
	}

	// Print assignment.
	names := []string{}
	for _, reg := range regs {
		names = append(names, string(reg))
	}

	p.Print(strings.Join(names, ", "), operator)

	// Mark declared.
	for _, reg := range regs {
		p.declared[reg] = true
	}
}
