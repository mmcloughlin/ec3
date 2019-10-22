package name

import "fmt"

// Sequence generates a sequence of names.
type Sequence interface {
	New() string
}

type SequenceFunc func() string

func (f SequenceFunc) New() string {
	return f()
}

// Indexed generates names using an increasing index and format string.
func Indexed(format string) Sequence {
	i := 0
	return SequenceFunc(func() string {
		v := fmt.Sprintf(format, i)
		i++
		return string(v)
	})
}

// Temporaries generates temporary variables in a standard form.
func Temporaries() Sequence {
	return Indexed("t%d")
}

type Uniquer interface {
	// MarkUsed marks names as used.
	MarkUsed(...string)
}

// UniqueGenerator is a method of generating new unique names.
type UniqueGenerator interface {
	Uniquer

	// New generates a new unique name.
	New(s Sequence) string
}

type uniquegenerator struct {
	used map[string]bool
}

// NewUniqueGenerator builds a UniqueGenerator based on a function that
// produces a sequence of possible variable names.
func NewUniqueGenerator() UniqueGenerator {
	return &uniquegenerator{
		used: make(map[string]bool),
	}
}

func (g *uniquegenerator) MarkUsed(vs ...string) {
	for _, v := range vs {
		g.used[v] = true
	}
}

func (g uniquegenerator) New(s Sequence) string {
	for {
		v := s.New()
		if !g.used[v] {
			g.MarkUsed(v)
			return v
		}
	}
}

// UniqueSequence generates unique names from a sequence.
type UniqueSequence interface {
	Uniquer
	Sequence
}

type uniquesequence struct {
	UniqueGenerator
	s Sequence
}

func Uniqued(s Sequence) UniqueSequence {
	return &uniquesequence{
		UniqueGenerator: NewUniqueGenerator(),
		s:               s,
	}
}

func (u *uniquesequence) New() string {
	return u.UniqueGenerator.New(u.s)
}
