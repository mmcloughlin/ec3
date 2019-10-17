package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>

extern void errorhandler(Z3_context c, Z3_error_code e);
*/
import "C"

import "unsafe"

// errorhandler is a callback for Z3 errors. Currently panics, simply to ensure errors do not go unnoticed.
// TODO(mbm): handle z3 errors more gracefully
//export errorhandler
func errorhandler(ctx C.Z3_context, code C.Z3_error_code) {
	msg := C.Z3_get_error_msg(ctx, code)
	panic(C.GoString(msg))
}

// Context is an environment for interacting with Z3.
type Context struct {
	ctx C.Z3_context
}

// NewContext builds a new Z3 environment with the given configuration.
func NewContext(cfg *Config) *Context {
	c := &Context{
		ctx: C.Z3_mk_context(cfg.cfg),
	}
	C.Z3_set_error_handler(c.ctx, (*C.Z3_error_handler)(C.errorhandler))
	return c
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
