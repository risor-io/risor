package object

import (
	"errors"
	"fmt"
)

// Error wraps string and implements Object interface.
type Error struct {
	// Message contains the error-message we're wrapping
	Message string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("Error(%s)", e.Message)
}

func (e *Error) Interface() interface{} {
	return errors.New(e.Message)
}

func (e *Error) Compare(other Object) (int, error) {
	typeComp := CompareTypes(e, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Error)
	if e.Message == otherStr.Message {
		return 0, nil
	}
	if e.Message > otherStr.Message {
		return 1, nil
	}
	return -1, nil
}

func (e *Error) Equals(other Object) Object {
	if other.Type() != ERROR {
		return False
	}
	if e.Message == other.(*Error).Message {
		return True
	}
	return False
}

func (e *Error) GetAttr(name string) (Object, bool) {
	return nil, false
}
