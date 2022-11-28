package object

import (
	"fmt"
)

// Integer wraps int64 and implements Object and Hashable interfaces.
type Integer struct {
	// Value holds the int64 wrapped by this object.
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

func (i *Integer) HashKey() Key {
	return Key{Type: i.Type(), IntValue: i.Value}
}

func (i *Integer) InvokeMethod(method string, args ...Object) Object {
	if method == "chr" {
		return &String{Value: string(rune(i.Value))}
	}
	return NewError("type error: %s object has no method %s", i.Type(), method)
}

func (i *Integer) ToInterface() interface{} {
	return i.Value
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%v)", i.Value)
}

func (i *Integer) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		if float64(i.Value) == other.Value {
			return 0, nil
		}
		if float64(i.Value) > other.Value {
			return 1, nil
		}
		return -1, nil
	case *Integer:
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

func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}
