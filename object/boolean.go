package object

import (
	"fmt"
)

// Boolean wraps bool and implements Object and Hashable interface.
type Boolean struct {
	// Value holds the boolean value we wrap.
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() Key {
	var value int64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return Key{Type: b.Type(), IntValue: value}
}

func (b *Boolean) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", b.Type(), method)
}

func (b *Boolean) ToInterface() interface{} {
	return b.Value
}

func (b *Boolean) String() string {
	return fmt.Sprintf("Boolean(%v)", b.Value)
}

func (b *Boolean) Compare(other Object) (int, error) {
	typeComp := CompareTypes(b, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherBool := other.(*Boolean)
	if b.Value == otherBool.Value {
		return 0, nil
	}
	if b.Value {
		return 1, nil
	}
	return -1, nil
}

func NewBoolean(value bool) *Boolean {
	if value {
		return TRUE
	}
	return FALSE
}
