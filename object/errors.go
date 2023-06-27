package object

import (
	"fmt"
	"math"
)

type ArgumentsError struct {
	message string
}

func (e *ArgumentsError) Error() string {
	return e.message
}

func NewArgumentsError(message string, args ...interface{}) error {
	return &ArgumentsError{message: fmt.Sprintf(message, args...)}
}

func NewArgsError(fn string, takes, given int) *Error {
	return NewError(NewArgumentsError("type error: %s() takes exactly %d arguments (%d given)",
		fn, takes, given))
}

func NewArgsRangeError(fn string, takesMin, takesMax, given int) *Error {
	if math.Abs(float64(takesMax-takesMin)) <= 0.0001 {
		return NewError(NewArgumentsError("type error: %s() takes %d or %d arguments (%d given)",
			fn, takesMin, takesMax, given))
	}
	return NewError(NewArgumentsError("type error: %s() takes between %d and %d arguments (%d given)",
		fn, takesMin, takesMax, given))
}
