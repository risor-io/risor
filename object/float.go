package object

import (
	"fmt"
	"strconv"
)

// Float wraps float64 and implements Object and Hashable interfaces.
type Float struct {
	// Value holds the float64 wrapped by this object.
	Value float64
}

func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (f *Float) Type() Type {
	return FLOAT
}

func (f *Float) HashKey() Key {
	return Key{Type: f.Type(), FltValue: f.Value}
}

func (f *Float) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", f.Type(), method)
}

func (f *Float) ToInterface() interface{} {
	return f.Value
}

func (f *Float) String() string {
	return fmt.Sprintf("Float(%v)", f.Value)
}

func (f *Float) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		if f.Value == other.Value {
			return 0, nil
		}
		if f.Value > other.Value {
			return 1, nil
		}
		return -1, nil
	case *Int:
		if f.Value == float64(other.Value) {
			return 0, nil
		}
		if f.Value > float64(other.Value) {
			return 1, nil
		}
		return -1, nil
	default:
		return CompareTypes(f, other), nil
	}
}

func (f *Float) Equals(other Object) Object {
	switch other.Type() {
	case INT:
		if f.Value == float64(other.(*Int).Value) {
			return True
		}
	case FLOAT:
		if f.Value == other.(*Float).Value {
			return True
		}
	}
	return False
}

func NewFloat(value float64) *Float {
	return &Float{Value: value}
}
