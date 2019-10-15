package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

import "errors"

type Solver struct {
	ctx    C.Z3_context
	solver C.Z3_solver
}

func (c *Context) Solver() *Solver {
	solver := C.Z3_mk_solver(c.ctx)
	C.Z3_solver_inc_ref(c.ctx, solver)

	return &Solver{
		ctx:    c.ctx,
		solver: solver,
	}
}

func (s *Solver) Close() error {
	C.Z3_solver_dec_ref(s.ctx, s.solver)
	return nil
}

func (s *Solver) Push() {
	C.Z3_solver_push(s.ctx, s.solver)
}

func (s *Solver) Pop() {
	s.PopN(1)
}

func (s *Solver) PopN(n uint) {
	C.Z3_solver_pop(s.ctx, s.solver, C.unsigned(n))
}

func (s *Solver) Assert(c *Bool) {
	C.Z3_solver_assert(s.ctx, s.solver, c.ast)
}

func (s *Solver) Check() (bool, error) {
	res := C.Z3_solver_check(s.ctx, s.solver)
	if res == C.Z3_L_UNDEF {	
		reason := C.Z3_solver_get_reason_unknown(s.ctx, s.solver)
		return false, errors.New(C.GoString(reason))
	}
	return res == C.Z3_L_TRUE, nil
}

func (s *Solver) Prove(f *Bool) (bool, error) {
	s.Push()
	defer s.Pop()
	s.Assert(f.Not())
	sat, err := s.Check()
	if err != nil {
		return false, err
	}
	return !sat, nil
}
