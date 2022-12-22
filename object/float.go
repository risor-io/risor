package object

import (
	"fmt"
	"strconv"
)

// Float wraps float64 and implements Object and Hashable interfaces.
type Float struct {
	// value holds the float64 wrapped by this object.
	value float64
}

func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

func (f *Float) Type() Type {
	return FLOAT
}

func (f *Float) Value() float64 {
	return f.value
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), FltValue: f.value}
}

func (f *Float) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (f *Float) Interface() interface{} {
	return f.value
}

func (f *Float) String() string {
	return fmt.Sprintf("float(%v)", f.value)
}

func (f *Float) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		if f.value == other.value {
			return 0, nil
		}
		if f.value > other.value {
			return 1, nil
		}
		return -1, nil
	case *Int:
		if f.value == float64(other.value) {
			return 0, nil
		}
		if f.value > float64(other.value) {
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
		if f.value == float64(other.(*Int).value) {
			return True
		}
	case FLOAT:
		if f.value == other.(*Float).value {
			return True
		}
	}
	return False
}

func NewFloat(value float64) *Float {
	return &Float{value: value}
}
