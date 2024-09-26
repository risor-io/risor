package object

import (
	"fmt"
	"math"

	"github.com/risor-io/risor/errz"
)

// Deprecated: retained for backwards compatibility only. Prefer errz.ArgsError.
type ArgumentsError = errz.ArgsError

// Deprecated: retained for backwards compatibility only. Prefer errz.ArgsErrorf.
func NewArgumentsError(message string, args ...interface{}) error {
	return &ArgumentsError{Err: fmt.Errorf(message, args...)}
}

func NewArgsError(fn string, takes, given int) *Error {
	return NewError(errz.ArgsErrorf("args error: %s() takes exactly %d arguments (%d given)",
		fn, takes, given))
}

func NewArgsRangeError(fn string, takesMin, takesMax, given int) *Error {
	if math.Abs(float64(takesMax-takesMin)) <= 0.0001 {
		return NewError(errz.ArgsErrorf("args error: %s() takes %d or %d arguments (%d given)",
			fn, takesMin, takesMax, given))
	}
	return NewError(errz.ArgsErrorf("args error: %s() takes between %d and %d arguments (%d given)",
		fn, takesMin, takesMax, given))
}
