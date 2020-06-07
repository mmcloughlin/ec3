// Package name provides utilities for symbol naming.
package name

import (
	"go/token"
	"strings"
)

// Standard naming schemes.
var (
	Plain     = SchemeFunc(func(name string) string { return name })
	LowerCase = SchemeFunc(strings.ToLower)
)

// Scheme defines a naming scheme.
type Scheme interface {
	Name(name string) string
}

// SchemeFunc adapts a function to the Scheme interface.
type SchemeFunc func(string) string

// Name calls f.
func (f SchemeFunc) Name(name string) string {
	return f(name)
}

// CompositeScheme builds a Scheme that applies the given schemes in order.
func CompositeScheme(schemes ...Scheme) Scheme {
	return SchemeFunc(func(name string) string {
		for _, scheme := range schemes {
			name = scheme.Name(name)
		}
		return name
	})
}

// Prefixed returns a scheme that adds the given prefix to all names. The result
// will have the same visibility as the input name.
func Prefixed(prefix string) Scheme {
	return SchemeFunc(func(name string) string {
		return SetExported(prefix+name, IsExported(name))
	})
}

// SetExported returns the given name
func SetExported(name string, exported bool) string {
	if name == "" {
		return ""
	}
	if exported {
		return strings.ToUpper(name[:1]) + name[1:]
	}
	return strings.ToLower(name[:1]) + name[1:]
}

// IsExported reports whether name is an exported Go symbol.
func IsExported(name string) bool {
	return token.IsExported(name)
}
