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
		n := fmt.Sprintf(format, i)
		i++
		return n
	})
}

// Temporaries generates temporary names in a standard form.
func Temporaries() Sequence {
	return Indexed("t%d")
}

type Uniquer interface {
	// Used reports whether a name is used.
	Used(string) bool

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
// produces a sequence of possible names.
func NewUniqueGenerator() UniqueGenerator {
	return &uniquegenerator{
		used: make(map[string]bool),
	}
}

func (g *uniquegenerator) Used(name string) bool {
	return g.used[name]
}

func (g *uniquegenerator) MarkUsed(names ...string) {
	for _, name := range names {
		g.used[name] = true
	}
}

func (g uniquegenerator) New(s Sequence) string {
	for {
		n := s.New()
		if !g.Used(n) {
			g.MarkUsed(n)
			return n
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
