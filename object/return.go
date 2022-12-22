package object

import "fmt"

// ReturnValue wraps Object and implements Object interface.
type ReturnValue struct {
	// Value is the object that is to be returned
	value Object
}

func (rv *ReturnValue) Type() Type {
	return RETURN_VALUE
}

func (rv *ReturnValue) Value() Object {
	return rv.value
}

func (rv *ReturnValue) Inspect() string {
	return rv.value.Inspect()
}

func (rv *ReturnValue) String() string {
	return fmt.Sprintf("return(%s)", rv.value)
}

func (rv *ReturnValue) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (rv *ReturnValue) Interface() interface{} {
	return rv.value.Interface()
}

func (rv *ReturnValue) Equals(other Object) Object {
	if other.Type() == RETURN_VALUE && rv == other.(*ReturnValue) {
		return True
	}
	return False
}

func NewReturnValue(value Object) *ReturnValue {
	if value == nil {
		panic("return value cannot be nil")
	}
	return &ReturnValue{value: value}
}
