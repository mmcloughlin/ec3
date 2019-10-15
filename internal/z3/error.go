package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// Error returns the error for the last API call, if any.
func (c *Context) Error() error {
	if code := C.Z3_get_error_code(c.ctx); code != C.Z3_OK {
		return &Error{
			Code:    int(code),
			Message: C.GoString(C.Z3_get_error_msg(c.ctx, code)),
		}
	}
	return nil
}
