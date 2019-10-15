package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// Context is an environment for interacting with Z3.
type Context struct {
	ctx C.Z3_context
}

// NewContext builds a new Z3 environment with the given configuration.
func NewContext(cfg *Config) *Context {
	return &Context{
		ctx: C.Z3_mk_context(cfg.cfg),
	}
}

// Close frees memory associated with this context.
func (c *Context) Close() error {
	C.Z3_del_context(c.ctx)
	return nil
}

func (c *Context) symbol(name string) C.Z3_symbol {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	return C.Z3_mk_string_symbol(c.ctx, n)
}
