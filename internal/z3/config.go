package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// Config contains configuration options.
type Config struct {
	cfg C.Z3_config
}

// NewConfig creates a configuration.
func NewConfig() *Config {
	return &Config{
		cfg: C.Z3_mk_config(),
	}
}

// Close frees memory associated with this configuration.
func (c *Config) Close() error {
	C.Z3_del_config(c.cfg)
	return nil
}
