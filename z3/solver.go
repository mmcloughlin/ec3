package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

import "errors"

// Solver checks a collection of predicates for satisfiability.
type Solver struct {
	ctx    C.Z3_context
	solver C.Z3_solver
}

// Solver builds an empty Solver tied to this context.
func (c *Context) Solver() *Solver {
	solver := C.Z3_mk_solver(c.ctx)
	C.Z3_solver_inc_ref(c.ctx, solver)
	return &Solver{
		ctx:    c.ctx,
		solver: solver,
	}
}

// SolverFor builds a solver customized for the given logic.
func (c *Context) SolverForLogic(logic string) *Solver {
	solver := C.Z3_mk_solver_for_logic(c.ctx, c.symbol(logic))
	C.Z3_solver_inc_ref(c.ctx, solver)
	return &Solver{
		ctx:    c.ctx,
		solver: solver,
	}
}

// Close frees memory associated with this solver.
func (s *Solver) Close() error {
	C.Z3_solver_dec_ref(s.ctx, s.solver)
	return nil
}

// Push creates a backtracking point.
func (s *Solver) Push() {
	C.Z3_solver_push(s.ctx, s.solver)
}

// Pop backtracks one layer in the stack.
func (s *Solver) Pop() {
	s.PopN(1)
}

// PopN backtracks n layers in the stack.
func (s *Solver) PopN(n uint) {
	C.Z3_solver_pop(s.ctx, s.solver, C.unsigned(n))
}

// Assert a constraint into the solver.
func (s *Solver) Assert(c *Bool) {
	C.Z3_solver_assert(s.ctx, s.solver, c.ast)
}

// Check whether the assertions in a given solver are consistent or not.
func (s *Solver) Check() (bool, error) {
	res := C.Z3_solver_check(s.ctx, s.solver)
	if res == C.Z3_L_UNDEF {
		reason := C.Z3_solver_get_reason_unknown(s.ctx, s.solver)
		return false, errors.New(C.GoString(reason))
	}
	return res == C.Z3_L_TRUE, nil
}

// Prove attempts to prove that f is true. This is a convenience for checking that not(f) is unsatisfiable.
func (s *Solver) Prove(f *Bool) (bool, error) {
	s.Assert(f.Not())
	sat, err := s.Check()
	if err != nil {
		return false, err
	}
	return !sat, nil
}
