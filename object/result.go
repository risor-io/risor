package object

import (
	"fmt"
)

// Result contains one of: an "ok" value and an "err" value
type Result struct {
	Ok  Object
	Err Object
}

func (rv *Result) Type() Type {
	return RESULT
}

func (rv *Result) Inspect() string {
	if rv.Ok != nil {
		return rv.Ok.Inspect()
	}
	return rv.Err.Inspect()
}

func (rv *Result) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "is_err":
		if len(args) != 0 {
			return NewError(fmt.Sprintf("type error: is_err() takes exactly 0 arguments (%d given)", len(args)))
		}
		if rv.Err != nil {
			return True
		}
		return False
	case "is_ok":
		if len(args) != 0 {
			return NewError(fmt.Sprintf("type error: is_ok() takes exactly 0 arguments (%d given)", len(args)))
		}
		if rv.Ok != nil {
			return True
		}
		return False
	case "unwrap":
		if len(args) != 0 {
			return NewError(fmt.Sprintf("type error: unwrap() takes exactly 0 arguments (%d given)", len(args)))
		}
		if rv.Ok != nil {
			return rv.Ok
		}
		return &Error{Message: fmt.Sprintf("result error: attempted to unwrap an error: %s", rv.Err)}
	case "unwrap_or":
		if len(args) != 1 {
			return NewError(fmt.Sprintf("type error: unwrap_or() takes exactly 1 argument (%d given)", len(args)))
		}
		if rv.Ok != nil {
			return rv.Ok
		}
		return args[0]
	case "expect":
		if len(args) != 1 {
			return NewError(fmt.Sprintf("type error: expect() takes exactly 1 argument (%d given)", len(args)))
		}
		arg := args[0]
		if _, ok := arg.(*String); !ok {
			return NewError(fmt.Sprintf("type error: expect() argument should be a string (%s given)", arg.Type()))
		}
		if rv.Ok != nil {
			return rv.Ok
		}
		return arg
	default:
		if rv.Ok != nil {
			return rv.Ok.InvokeMethod(method, args...)
		}
		return NewError(fmt.Sprintf("result error: %v", rv.Err))
	}
}

func (rv *Result) ToInterface() interface{} {
	if rv.Ok != nil {
		return rv.Ok.ToInterface()
	}
	return rv.Err.ToInterface()
}

func (rv *Result) String() string {
	if rv.Ok != nil {
		return fmt.Sprintf("Result(Ok: %s)", rv.Ok)
	} else if rv.Err != nil {
		return fmt.Sprintf("Result(Err: %s)", rv.Err)
	}
	return "Result()"
}

func (rv *Result) IsOk() bool {
	return rv.Ok != nil
}

func (rv *Result) IsErr() bool {
	return rv.Err != nil
}

func (rv *Result) Equals(other Object) Object {
	if other.Type() != RESULT {
		return False
	}
	otherResult := other.(*Result)
	if rv.Ok != nil {
		return NewBoolean(otherResult.Ok != nil && rv.Ok.Equals(otherResult.Ok).(*Bool).Value)
	} else if rv.Err != nil {
		return NewBoolean(otherResult.Err != nil && rv.Err.Equals(otherResult.Err).(*Bool).Value)
	}
	if otherResult.Ok == nil && otherResult.Err == nil {
		return True
	}
	return False
}
