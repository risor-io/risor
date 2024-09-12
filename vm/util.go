package vm

import (
	"errors"
	"fmt"

	"github.com/risor-io/risor/object"
)

func checkCallArgs(fn *object.Function, argc int) error {
	// Number of parameters in the function signature
	paramsCount := len(fn.Parameters())

	// Number of required args when the function is called (those without defaults)
	requiredArgsCount := fn.RequiredArgsCount()

	// Check if too many or too few arguments were passed
	if argc > paramsCount || argc < requiredArgsCount {
		msg := "type error: function"
		if name := fn.Name(); name != "" {
			msg = fmt.Sprintf("%s %q", msg, name)
		}
		switch paramsCount {
		case 0:
			msg = fmt.Sprintf("%s takes 0 arguments (%d given)", msg, argc)
		case 1:
			msg = fmt.Sprintf("%s takes 1 argument (%d given)", msg, argc)
		default:
			msg = fmt.Sprintf("%s takes %d arguments (%d given)", msg, paramsCount, argc)
		}
		return errors.New(msg)
	}
	return nil
}
