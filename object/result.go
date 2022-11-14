package object

import "fmt"

// Result contains one of: an "ok" value and an "err" value
type Result struct {
	Ok  Object
	Err Object
}

func (rv *Result) Type() Type {
	return RESULT_OBJ
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
		if rv.Err != nil {
			return TRUE
		}
		return FALSE
	case "is_ok":
		if rv.Ok != nil {
			return TRUE
		}
		return FALSE
	case "unwrap":
		if rv.Ok != nil {
			return rv.Ok
		}
		return &Error{Message: fmt.Sprintf("result error: %v", rv.Err)}
	case "unwrap_or":
		if len(args) != 1 {
			return &Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
		}
		if rv.Ok != nil {
			return rv.Ok
		}
		return args[0]
	case "expect":
		if len(args) != 1 {
			return &Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
		}
		arg := args[0]
		if _, ok := arg.(*String); !ok {
			return &Error{Message: fmt.Sprintf("expected a string; got %s", arg.Type())}
		}
		if rv.Ok != nil {
			return rv.Ok
		}
		return arg
	default:
		if rv.Ok != nil {
			return rv.Ok.InvokeMethod(method, args...)
		}
		return &Error{Message: fmt.Sprintf("result error: %v", rv.Err)}
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
