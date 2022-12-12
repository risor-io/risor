package object

import (
	"fmt"
)

// Int wraps int64 and implements Object and Hashable interfaces.
type Int struct {
	// Value holds the int64 wrapped by this object.
	Value int64
}

func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Int) Type() Type {
	return INT
}

func (i *Int) HashKey() Key {
	return Key{Type: i.Type(), IntValue: i.Value}
}

func (i *Int) InvokeMethod(method string, args ...Object) Object {
	if method == "chr" {
		return &String{Value: string(rune(i.Value))}
	}
	return NewError("type error: %s object has no method %s", i.Type(), method)
}

func (i *Int) ToInterface() interface{} {
	return i.Value
}

func (i *Int) String() string {
	return fmt.Sprintf("Int(%v)", i.Value)
}

func (i *Int) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		if float64(i.Value) == other.Value {
			return 0, nil
		}
		if float64(i.Value) > other.Value {
			return 1, nil
		}
		return -1, nil
	case *Int:
		if i.Value == other.Value {
			return 0, nil
		}
		if i.Value > other.Value {
			return 1, nil
		}
		return -1, nil
	default:
		return CompareTypes(i, other), nil
	}
}

func NewInt(value int64) *Int {
	return &Int{Value: value}
}
