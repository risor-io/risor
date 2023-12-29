package http

import "github.com/risor-io/risor/object"

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"fetch": object.NewBuiltin("fetch", Fetch),
	}
}
