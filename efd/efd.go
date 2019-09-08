package efd

//go:generate go run make.go -input data/efd.tar.gz -input data -output zdb.go

var All = Select()

type Predicate func(*Formula) bool

type Formulae []*Formula

func (f Formulae) Filter(predicates ...Predicate) Formulae {
	if len(predicates) == 0 {
		return f
	}
	result := make(Formulae, 0)
	for _, formula := range f {
		keep := true
		for _, predicate := range predicates {
			keep = predicate(formula) && keep
		}
		if keep {
			result = append(result, formula)
		}
	}
	return result
}

func Select(predicates ...Predicate) Formulae {
	return Formulae(formulae).Filter(predicates...)
}

func WithClass(class string) Predicate {
	return func(f *Formula) bool { return f.Class == class }
}

func WithShape(shape string) Predicate {
	return func(f *Formula) bool { return f.Shape.Tag == shape }
}

func WithRepresentation(repr string) Predicate {
	return func(f *Formula) bool { return f.Representation.Tag == repr }
}

func WithOperation(op string) Predicate {
	return func(f *Formula) bool { return f.Operation == op }
}

func WithProgram(f *Formula) bool {
	return f.Program != nil
}

func LookupRepresentation(id string) *Representation {
	for _, r := range representations {
		if r.ID == id {
			return r
		}
	}
	return nil
}

func LookupFormula(id string) *Formula {
	for _, f := range formulae {
		if f.ID == id {
			return f
		}
	}
	return nil
}
