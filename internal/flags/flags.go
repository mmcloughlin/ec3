// Package flags provides helpers for the standard flag package.
package flags

import (
	"flag"
	"fmt"
)

type strings struct {
	value []string
}

func (s strings) String() string {
	return fmt.Sprint(s.value)
}

func (s *strings) Set(val string) error {
	s.value = append(s.value, val)
	return nil
}

// Strings defines a string slice flag with specified name, default value, and
// usage string. The return value is the address of a string slice variable that
// stores the value of the flag.
func Strings(name string, value []string, usage string) *[]string {
	s := append([]string{}, value...)
	v := &strings{
		value: s,
	}
	flag.Var(v, name, usage)
	return &v.value
}
