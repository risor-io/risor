package object

import (
	"context"
	"fmt"
	"strings"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

// StackFrame represents a single frame in a traceback
type StackFrame struct {
	FunctionName string
	FileName     string
	LineNumber   int
}

// Error wraps a Go error interface and implements Object.
type Error struct {
	*base
	err       error
	raised    bool
	traceback []StackFrame
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("error(%q)", e.err.Error())
}

func (e *Error) String() string {
	return e.err.Error()
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
		return 0, errz.TypeErrorf("type error: unable to compare error and %s", other.Type())
	}
	thisMsg := e.Message().Value()
	otherMsg := otherErr.Message().Value()
	if thisMsg == otherMsg && e.raised == otherErr.raised {
		return 0, nil
	}
	if thisMsg > otherMsg {
		return 1, nil
	}
	if thisMsg < otherMsg {
		return -1, nil
	}
	if e.raised && !otherErr.raised {
		return 1, nil
	}
	if !e.raised && otherErr.raised {
		return -1, nil
	}
	return 0, nil
}

func (e *Error) Equals(other Object) Object {
	switch other := other.(type) {
	case *Error:
		if e.Message().Value() == other.Message().Value() && e.raised == other.raised {
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
	case "message":
		return NewBuiltin("message", func(ctx context.Context, args ...Object) Object {
			return e.Message()
		}), true
	case "traceback":
		return NewBuiltin("traceback", func(ctx context.Context, args ...Object) Object {
			return e.Traceback()
		}), true
	default:
		return nil, false
	}
}

func (e *Error) Message() *String {
	return NewString(e.err.Error())
}

func (e *Error) WithRaised(value bool) *Error {
	fmt.Printf("DEBUG: WithRaised called - current traceback len: %d\n", len(e.traceback))
	e.raised = value
	return e
}

func (e *Error) IsRaised() bool {
	return e.raised
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) RunOperation(opType op.BinaryOpType, right Object) Object {
	return TypeErrorf("type error: unsupported operation for error: %v", opType)
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
	return &Error{base: &base{}, err: fmt.Errorf(format, args...), raised: true}
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal error")
}

func NewError(err error) *Error {
	fmt.Printf("DEBUG: NewError called with error type: %T\n", err)
	switch err := err.(type) {
	case *Error: // unwrap to get the inner error, to avoid unhelpful nesting
		fmt.Printf("DEBUG: NewError - input is *Error with %d traceback frames\n", len(err.traceback))
		// Preserve the traceback from the original error
		newErr := &Error{base: &base{}, err: err.Unwrap(), raised: true, traceback: err.traceback}
		fmt.Printf("DEBUG: NewError - created new Error with %d traceback frames\n", len(newErr.traceback))
		return newErr
	default:
		fmt.Printf("DEBUG: NewError - input is not *Error, creating new error without traceback\n")
		return &Error{base: &base{}, err: err, raised: true}
	}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}

func (e *Error) WithTraceback(traceback []StackFrame) *Error {
	e.traceback = traceback
	return e
}

func (e *Error) Traceback() *String {
	fmt.Printf("DEBUG: Error.Traceback() called - traceback len: %d\n", len(e.traceback))
	for i, frame := range e.traceback {
		fmt.Printf("DEBUG: Frame %d: %s\n", i, frame.FunctionName)
	}
	
	if len(e.traceback) == 0 {
		return NewString("No traceback available")
	}
	
	var builder strings.Builder
	builder.WriteString("Traceback (most recent call last):\n")
	for _, frame := range e.traceback {
		if frame.FileName != "" {
			builder.WriteString(fmt.Sprintf("  File \"%s\", line %d, in %s\n", 
				frame.FileName, frame.LineNumber, frame.FunctionName))
		} else {
			builder.WriteString(fmt.Sprintf("  in %s\n", frame.FunctionName))
		}
	}
	builder.WriteString(fmt.Sprintf("Error: %s", e.err.Error()))
	return NewString(builder.String())
}

func NewErrorWithTraceback(err error, traceback []StackFrame) *Error {
	switch err := err.(type) {
	case *Error: // unwrap to get the inner error, to avoid unhelpful nesting
		return &Error{base: &base{}, err: err.Unwrap(), raised: true, traceback: traceback}
	default:
		return &Error{base: &base{}, err: err, raised: true, traceback: traceback}
	}
}

func ErrorfWithTraceback(traceback []StackFrame, format string, a ...interface{}) *Error {
	var args []interface{}
	for _, arg := range a {
		if obj, ok := arg.(Object); ok {
			args = append(args, obj.Interface())
		} else {
			args = append(args, arg)
		}
	}
	fmt.Printf("DEBUG: ErrorfWithTraceback creating error with %d frames\n", len(traceback))
	return &Error{base: &base{}, err: fmt.Errorf(format, args...), raised: true, traceback: traceback}
}
