package gen

import (
	"github.com/mmcloughlin/avo/ir"
	"github.com/mmcloughlin/avo/printer"
)

func Asm(pkg string, f *ir.File) Interface {
	return asm(pkg, f, printer.NewGoAsm)
}

func Stubs(pkg string, f *ir.File) Interface {
	return asm(pkg, f, printer.NewStubs)
}

func asm(pkg string, f *ir.File, b printer.Builder) Interface {
	p := b(printer.Config{
		Pkg:  pkg,
		Name: GeneratedBy,
	})
	return Func(func() ([]byte, error) {
		return p.Print(f)
	})
}
