package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

// Result contains one of: an "ok" value or an "err" value
type Result struct {
	ok  Object
	err *Error
}

func (rv *Result) Type() Type {
	return RESULT
}

func (rv *Result) Inspect() string {
	if rv.ok != nil {
		return fmt.Sprintf("ok(%s)", rv.ok.Inspect())
	}
	return fmt.Sprintf("err(%q)", rv.err.Value())
}

func (rv *Result) GetAttr(name string) (Object, bool) {
	switch name {
	case "is_err":
		return &Builtin{
			name: "result.is_err",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("is_err", 0, len(args))
				}
				return NewBool(rv.IsErr())
			},
		}, true
	case "is_ok":
		return &Builtin{
			name: "result.is_ok",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("is_ok", 0, len(args))
				}
				return NewBool(rv.IsOk())
			},
		}, true
	case "unwrap":
		return &Builtin{
			name: "result.unwrap",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("unwrap", 0, len(args))
				}
				return rv.Unwrap()
			},
		}, true
	case "unwrap_err":
		return &Builtin{
			name: "result.unwrap_err",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("unwrap_err", 0, len(args))
				}
				return rv.UnwrapErr()
			},
		}, true
	case "err_msg":
		return &Builtin{
			name: "result.err_msg",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("err_msg", 0, len(args))
				}
				return rv.ErrMsg()
			},
		}, true
	case "unwrap_or":
		return &Builtin{
			name: "result.unwrap_or",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("unwrap_or", 1, len(args))
				}
				return rv.UnwrapOr(args[0])
			},
		}, true
	case "expect":
		return &Builtin{
			name: "result.expect",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("expect", 1, len(args))
				}
				return rv.Expect(args[0])
			},
		}, true
	default:
		if rv.ok != nil {
			return rv.ok.GetAttr(name)
		} else if rv.err != nil {
			return Errorf("result error: attempted to unwrap an error: %s", rv.err), true
		}
	}
	return nil, false
}

func (rv *Result) Interface() interface{} {
	if rv.ok != nil {
		return rv.ok.Interface()
	}
	return rv.err.Interface()
}

func (rv *Result) String() string {
	if rv.ok != nil {
		return fmt.Sprintf("result(ok: %s)", rv.ok)
	} else if rv.err != nil {
		return fmt.Sprintf("result(err: %s)", rv.err)
	}
	return "result()"
}

func (rv *Result) IsOk() bool {
	return rv.ok != nil
}

func (rv *Result) IsErr() bool {
	return rv.err != nil
}

func (rv *Result) Unwrap() Object {
	if rv.ok != nil {
		return rv.ok
	}
	return Errorf("result error: unwrap() called on an error: %s", rv.err.Inspect())
}

func (rv *Result) UnwrapErr() *Error {
	if rv.err != nil {
		return rv.err
	}
	return Errorf("result error: unwrap_err() called on an ok: %s", rv.ok.Inspect())
}

func (rv *Result) ErrMsg() *String {
	if rv.err != nil {
		return rv.err.Message()
	}
	return NewString("")
}

func (rv *Result) UnwrapOr(other Object) Object {
	if rv.ok != nil {
		return rv.ok
	}
	return other
}

func (rv *Result) Expect(other Object) Object {
	if _, ok := other.(*String); !ok {
		return Errorf("type error: expect() argument should be a string (%s given)", other.Type())
	}
	if rv.ok != nil {
		return rv.ok
	}
	return other
}

func (rv *Result) Equals(other Object) Object {
	if other.Type() != RESULT {
		return False
	}
	otherResult := other.(*Result)
	if rv.ok != nil {
		return NewBool(otherResult.ok != nil && rv.ok.Equals(otherResult.ok).(*Bool).value)
	} else if rv.err != nil {
		return NewBool(otherResult.err != nil && rv.err.Equals(otherResult.err).(*Bool).value)
	}
	if otherResult.ok == nil && otherResult.err == nil {
		return True
	}
	return False
}

func (rv *Result) IsTruthy() bool {
	return rv.IsOk()
}

func (rv *Result) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for result: %v", opType))
}

func NewErrResult(err *Error) *Result {
	return &Result{err: err}
}

func NewOkResult(ok Object) *Result {
	return &Result{ok: ok}
}
