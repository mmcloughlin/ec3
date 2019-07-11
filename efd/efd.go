package efd

//go:generate go run make.go -archive efd.tar.gz -output zdb.go

var All = Formulae()

type Predicate func(*Formula) bool

type Collection []*Formula

func (c Collection) Filter(predicates ...Predicate) Collection {
	if len(predicates) == 0 {
		return c
	}
	result := make(Collection, 0)
	for _, f := range c {
		keep := true
		for _, predicate := range predicates {
			keep = predicate(f) && keep
		}
		if keep {
			result = append(result, f)
		}
	}
	return result
}

func Formulae(predicates ...Predicate) Collection {
	return Collection(formulae).Filter(predicates...)
}
