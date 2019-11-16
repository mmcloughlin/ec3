package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

import (
	"strconv"
	"time"
	"unsafe"
)

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

// SetParam sets a configuration parameter.
func (c *Config) SetParam(id, value string) {
	cid := C.CString(id)
	cvalue := C.CString(value)

	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cvalue))

	C.Z3_set_param_value(c.cfg, cid, cvalue)
}

// SetParamBool sets a boolean parameter.
func (c *Config) SetParamBool(id string, value bool) {
	c.SetParam(id, strconv.FormatBool(value))
}

// SetProof configures proof generation.
func (c *Config) SetProof(enabled bool) {
	c.SetParamBool("proof", enabled)
}

// SetTrace configures tracing support for VCC.
func (c *Config) SetTrace(enabled bool) {
	c.SetParamBool("trace", enabled)
}

// SetTraceFilename sets the trace output file for VCC traces.
func (c *Config) SetTraceFilename(filename string) {
	c.SetParam("trace_file_name", filename)
}

// SetTimeout configures the default timeout used for solvers.
func (c *Config) SetTimeout(d time.Duration) {
	c.SetParam("timeout", strconv.FormatInt(d.Milliseconds(), 10))
}

// SetAutoConfig configures whether to use heuristics to automatically select
// solver and configure it.
func (c *Config) SetAutoConfig(enabled bool) {
	c.SetParamBool("auto_config", enabled)
}

// SetModel configures model generation for solvers. This parameter can be
// overwritten when creating a solver.
func (c *Config) SetModel(enabled bool) {
	c.SetParamBool("model", enabled)
}

// SetModelValidate configures whether to validate models produced by solvers.
func (c *Config) SetModelValidate(enabled bool) {
	c.SetParamBool("model_validate", enabled)
}
