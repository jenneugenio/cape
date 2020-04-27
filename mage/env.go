package mage

import (
	"os"
)

// Env represents a set of environment variables
type Env map[string]string

// Source checks the os environment for the current keys. If any are set (e.g.
// no-zero length string) then the value is overridden.
func (e Env) Source() {
	for k := range e {
		v := os.Getenv(k)
		if len(v) > 0 {
			e[k] = v
		}
	}
}
