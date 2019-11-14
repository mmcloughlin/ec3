package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// Params is a parameter set used to configure many components such as:
// simplifiers, tactics, solvers, etc.
type Params struct {
	ctx    *Context
	params C.Z3_params
}

// Params builds an empty parameter set tied to this context.
func (c *Context) Params() *Params {
	params := C.Z3_mk_params(c.ctx)
	C.Z3_params_inc_ref(c.ctx, params)
	return &Params{
		ctx:    c,
		params: params,
	}
}

// Close frees memory associated with the parameter set.
func (p *Params) Close() error {
	C.Z3_params_dec_ref(p.ctx.ctx, p.params)
	return nil
}

func (p *Params) String() string {
	return C.GoString(C.Z3_params_to_string(p.ctx.ctx, p.params))
}

// SetBool sets a boolean parameter k with value v.
func (p *Params) SetBool(k string, v bool) {
	C.Z3_params_set_bool(p.ctx.ctx, p.params, p.ctx.symbol(k), C.bool(v))
}

// SetUint sets an unsigned integer parameter k with value v.
func (p *Params) SetUint(k string, v uint) {
	C.Z3_params_set_uint(p.ctx.ctx, p.params, p.ctx.symbol(k), C.unsigned(v))
}

// SetFloat64 sets an floating point parameter k with value v.
func (p *Params) SetFloat64(k string, v float64) {
	C.Z3_params_set_double(p.ctx.ctx, p.params, p.ctx.symbol(k), C.double(v))
}

// SetString sets a string parameter k with value v.
func (p *Params) SetString(k, v string) {
	C.Z3_params_set_symbol(p.ctx.ctx, p.params, p.ctx.symbol(k), p.ctx.symbol(v))
}
