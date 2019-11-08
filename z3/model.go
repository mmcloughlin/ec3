package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// Model is a model from a solver.
type Model struct {
	ctx   C.Z3_context
	model C.Z3_model
}

// Model retrieves the model from the last successful call to check.
func (s *Solver) Model() *Model {
	m := &Model{
		ctx:   s.ctx,
		model: C.Z3_solver_get_model(s.ctx, s.solver),
	}
	C.Z3_model_inc_ref(m.ctx, m.model)
	return m
}

// Close frees memory associated with this model.
func (m *Model) Close() error {
	C.Z3_model_dec_ref(m.ctx, m.model)
	return nil
}

func (m *Model) String() string {
	return C.GoString(C.Z3_model_to_string(m.ctx, m.model))
}

// Assignments returns assignments for all constants in the model.
func (m *Model) Assignments() (map[string]Value, error) {
	values := map[string]Value{}

	n := uint(C.Z3_model_get_num_consts(m.ctx, m.model))
	for i := uint(0); i < n; i++ {
		decl := C.Z3_model_get_const_decl(m.ctx, m.model, C.unsigned(i))

		// Determine declaration name.
		sym := C.Z3_get_decl_name(m.ctx, decl)
		name := C.GoString(C.Z3_get_symbol_string(m.ctx, sym))

		// Interpret the constant in the model.
		ast := C.Z3_model_get_const_interp(m.ctx, m.model, decl)

		// Returns nil if the constant does not have an interpretation, meaning the value does not matter.
		if ast == nil {
			values[name] = nil
			continue
		}

		// Convert to specific type.
		v, err := newvalue(m.ctx, ast)
		if err != nil {
			return nil, err
		}

		values[name] = v
	}

	return values, nil
}
