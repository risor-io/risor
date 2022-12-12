package object

import "fmt"

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
