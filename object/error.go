package object

// Error wraps string and implements Object interface.
type Error struct {
	// Message contains the error-message we're wrapping
	Message string
}

// Type returns the type of this object.
func (e *Error) Type() Type {
	return ERROR_OBJ
}

// Inspect returns a string-representation of the given object.
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", e.Type(), method)
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (e *Error) ToInterface() interface{} {
	return "<ERROR>"
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
