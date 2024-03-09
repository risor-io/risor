package vm

import (
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
		switch paramsCount {
		case 0:
			return fmt.Errorf("type error: function takes no arguments (%d given)", argc)
		case 1:
			return fmt.Errorf("type error: function takes 1 argument (%d given)", argc)
		default:
			return fmt.Errorf("type error: function takes %d arguments (%d given)", paramsCount, argc)
		}
	}
	return nil
}
