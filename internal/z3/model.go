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
