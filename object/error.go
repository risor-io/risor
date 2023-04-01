package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

// Error wraps a Go error interface and implements Object.
type Error struct {
	// err is the Go error being wrapped.
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
	typeComp := CompareTypes(e, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Error)
	thisMsg := e.Message().Value()
	otherMsg := otherStr.Message().Value()
	if thisMsg == otherMsg {
		return 0, nil
	}
	if thisMsg > otherMsg {
		return 1, nil
	}
	return -1, nil
}

func (e *Error) Equals(other Object) Object {
	if other.Type() != ERROR {
		return False
	}
	if e.Message() == other.(*Error).Message() {
		return True
	}
	return False
}

func (e *Error) IsTruthy() bool {
	return true
}

func (e *Error) Message() *String {
	return NewString(e.err.Error())
}

func (e *Error) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (e *Error) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for error: %v", opType))
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

func NewError(err error) *Error {
	return &Error{err: err}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}
