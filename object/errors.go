package object

import (
	"fmt"
	"math"
)

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}

func NewErrorResult(format string, a ...interface{}) *Result {
	return &Result{Err: NewError(format, a...)}
}

func NewOkResult(value Object) *Result {
	return &Result{Ok: value}
}

func NewArgsError(fn string, takes, given int) *Error {
	return NewError(fmt.Sprintf("type error: %s() takes exactly %d arguments (%d given)",
		fn, takes, given))
}

func NewArgsRangeError(fn string, takesMin, takesMax, given int) *Error {
	if math.Abs(float64(takesMax-takesMin)) <= 1.0 {
		return NewError(fmt.Sprintf("type error: %s() takes %d or %d arguments (%d given)",
			fn, takesMin, takesMax, given))
	}
	return NewError(fmt.Sprintf("type error: %s() takes between %d and %d arguments (%d given)",
		fn, takesMin, takesMax, given))
}
