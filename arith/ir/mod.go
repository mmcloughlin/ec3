package ir

type Section interface {
	section() // sealed
}

type Module struct {
	Sections []Section
}

type Function struct {
	Program

	Name      string
	Signature *Signature
}

func (Function) section() {}

type Signature struct {
	Params  []*Var
	Results []*Var
}

// Var looks up a variable among parameters and results, returning nil if not found.
func (s *Signature) Var(name string) *Var {
	vars := []*Var{}
	vars = append(vars, s.Params...)
	vars = append(vars, s.Results...)
	return lookupvar(name, vars)
}

// Param looks up a parameter by name, returning nil if not found.
func (s *Signature) Param(name string) *Var {
	return lookupvar(name, s.Params)
}

// Result looks up a result by name, returning nil if not found.
func (s *Signature) Result(name string) *Var {
	return lookupvar(name, s.Results)
}

func lookupvar(name string, vs []*Var) *Var {
	for _, v := range vs {
		if v.Name == name {
			return v
		}
	}
	return nil
}

type Var struct {
	Name string
	Type Type
}

func NewVar(name string, t Type) *Var {
	return &Var{Name: name, Type: t}
}

func NewVars(t Type, names ...string) []*Var {
	vs := []*Var{}
	for _, name := range names {
		v := NewVar(name, t)
		vs = append(vs, v)
	}
	return vs
}

type Type interface {
	typ() // sealed
}

// Integer is a multi-precision integer type.
type Integer struct {
	K uint // number of limbs
}

func (Integer) typ() {}
