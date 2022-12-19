package object

import (
	"context"
	"fmt"
)

// Result contains one of: an "ok" value or an "err" value
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

func (rv *Result) GetAttr(name string) (Object, bool) {
	switch name {
	case "is_err":
		return &Builtin{
			Name: "result.is_err",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("is_err", 0, len(args))
				}
				return NewBool(rv.IsErr())
			},
		}, true
	case "is_ok":
		return &Builtin{
			Name: "result.is_ok",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("is_ok", 0, len(args))
				}
				return NewBool(rv.IsOk())
			},
		}, true
	case "unwrap":
		return &Builtin{
			Name: "result.unwrap",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("unwrap", 0, len(args))
				}
				return rv.Unwrap()
			},
		}, true
	case "unwrap_or":
		return &Builtin{
			Name: "result.unwrap_or",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("unwrap_or", 1, len(args))
				}
				return rv.UnwrapOr(args[0])
			},
		}, true
	case "expect":
		return &Builtin{
			Name: "result.expect",
			Fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("expect", 1, len(args))
				}
				return rv.Expect(args[0])
			},
		}, true
	default:
		if rv.Ok != nil {
			return rv.Ok.GetAttr(name)
		}
	}
	return nil, false
}

func (rv *Result) Interface() interface{} {
	if rv.Ok != nil {
		return rv.Ok.Interface()
	}
	return rv.Err.Interface()
}

func (rv *Result) String() string {
	if rv.Ok != nil {
		return fmt.Sprintf("result(Ok: %s)", rv.Ok)
	} else if rv.Err != nil {
		return fmt.Sprintf("result(Err: %s)", rv.Err)
	}
	return "result()"
}

func (rv *Result) IsOk() bool {
	return rv.Ok != nil
}

func (rv *Result) IsErr() bool {
	return rv.Err != nil
}

func (rv *Result) Unwrap() Object {
	if rv.Ok != nil {
		return rv.Ok
	}
	return &Error{Message: fmt.Sprintf("result error: attempted to unwrap an error: %s", rv.Err)}
}

func (rv *Result) UnwrapOr(other Object) Object {
	if rv.Ok != nil {
		return rv.Ok
	}
	return other
}

func (rv *Result) Expect(other Object) Object {
	if _, ok := other.(*String); !ok {
		return NewError(fmt.Sprintf("type error: expect() argument should be a string (%s given)", other.Type()))
	}
	if rv.Ok != nil {
		return rv.Ok
	}
	return other
}

func (rv *Result) Equals(other Object) Object {
	if other.Type() != RESULT {
		return False
	}
	otherResult := other.(*Result)
	if rv.Ok != nil {
		return NewBool(otherResult.Ok != nil && rv.Ok.Equals(otherResult.Ok).(*Bool).Value)
	} else if rv.Err != nil {
		return NewBool(otherResult.Err != nil && rv.Err.Equals(otherResult.Err).(*Bool).Value)
	}
	if otherResult.Ok == nil && otherResult.Err == nil {
		return True
	}
	return False
}
