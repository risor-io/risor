package arg

import (
	"github.com/risor-io/risor/object"
)

func Require(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		if count == 1 {
			return object.ArgsErrorf(
				"args error: %s() takes exactly 1 argument (%d given)",
				funcName, nArgs)
		}
		return object.ArgsErrorf(
			"args error: %s() takes exactly %d arguments (%d given)",
			funcName, count, nArgs)
	}
	return nil
}

func RequireRange(funcName string, min, max int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs < min {
		return object.ArgsErrorf(
			"args error: %s() takes at least %d %s (%d given)",
			funcName, min, pluralize("argument", nArgs > 1), nArgs)
	} else if nArgs > max {
		return object.ArgsErrorf(
			"args error: %s() takes at most %d %s (%d given)",
			funcName, max, pluralize("argument", nArgs > 1), nArgs)
	}
	return nil
}

func pluralize(s string, do bool) string {
	if do {
		return s + "s"
	}
	return s
}
