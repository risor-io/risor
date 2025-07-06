package all

import (
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
)

// Builtins returns a map of standard builtins and globals for Risor scripts.
// This includes only the builtins and modules that are always available, without
// pulling in additional Go modules.
//
// Deprecated: Use risor.DefaultGlobals instead.
func Builtins() map[string]object.Object {
	globals := risor.DefaultGlobals(risor.DefaultGlobalsOpts{
		ListenersAllowed: true,
	})
	result := map[string]object.Object{}
	for k, v := range globals {
		result[k] = v.(object.Object)
	}
	return result
}
