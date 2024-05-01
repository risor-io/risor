package object

import (
	"context"
	"errors"
	"fmt"

	"github.com/risor-io/risor/op"
)

// Error wraps a Go error interface and implements Object.
type Error struct {
	*base
	err error
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("error(%q)", e.err.Error())
}

func (e *Error) String() string {
	return fmt.Sprintf("error(%s)", e.err.Error())
}

func (e *Error) Value() error {
	return e.err
}

func (e *Error) Interface() interface{} {
	return e.err
}

func (e *Error) Compare(other Object) (int, error) {
	otherErr, ok := other.(*Error)
	if !ok {
		return 0, fmt.Errorf("type error: unable to compare error and %s", other.Type())
	}
	thisMsg := e.Message().Value()
	otherMsg := otherErr.Message().Value()
	if thisMsg == otherMsg {
		return 0, nil
	}
	if thisMsg > otherMsg {
		return 1, nil
	}
	return -1, nil
}

func (e *Error) Equals(other Object) Object {
	switch other := other.(type) {
	case *Error:
		if e.Message() == other.Message() {
			return True
		}
		return False
	default:
		return False
	}
}

func (e *Error) GetAttr(name string) (Object, bool) {
	switch name {
	case "error":
		return NewBuiltin("error", func(ctx context.Context, args ...Object) Object {
			return e.Message()
		}), true
	default:
		return nil, false
	}
}

func (e *Error) Message() *String {
	return NewString(e.err.Error())
}

func (e *Error) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for error: %v", opType))
}

func Errorf(format string, a ...interface{}) *Error {
	var args []interface{}
	for _, arg := range a {
		if obj, ok := arg.(Object); ok {
			args = append(args, obj.Interface())
		} else {
			args = append(args, arg)
		}
	}
	return &Error{err: fmt.Errorf(format, args...)}
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return nil, errors.New("type error: unable to marshal error")
}

func NewError(err error) *Error {
	return &Error{err: err}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}
