package object

import (
	"fmt"
)

// Bool wraps bool and implements Object and Hashable interface.
type Bool struct {
	// Value holds the boolean value we wrap.
	Value bool
}

func (b *Bool) Type() Type {
	return BOOL
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Bool) HashKey() Key {
	var value int64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return Key{Type: b.Type(), IntValue: value}
}

func (b *Bool) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (b *Bool) ToInterface() interface{} {
	return b.Value
}

func (b *Bool) String() string {
	return fmt.Sprintf("Bool(%v)", b.Value)
}

func (b *Bool) Compare(other Object) (int, error) {
	typeComp := CompareTypes(b, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherBool := other.(*Bool)
	if b.Value == otherBool.Value {
		return 0, nil
	}
	if b.Value {
		return 1, nil
	}
	return -1, nil
}

func (b *Bool) Equals(other Object) Object {
	if other.Type() == BOOL && b.Value == other.(*Bool).Value {
		return True
	}
	return False
}

func NewBool(value bool) *Bool {
	if value {
		return True
	}
	return False
}

func Not(b *Bool) *Bool {
	if b.Value {
		return False
	}
	return True
}

func Equals(a, b Object) bool {
	return a.Equals(b).(*Bool).Value
}

func IsTruthy(obj Object) bool {
	switch obj {
	case Nil:
		return false
	case True:
		return true
	case False:
		return false
	default:
		switch obj := obj.(type) {
		case *Int:
			return obj.Value != 0
		case *Float:
			return obj.Value != 0.0
		case *String:
			return obj.Value != ""
		case *List:
			return len(obj.Items) > 0
		case *Map:
			return len(obj.Items) > 0
		case *Set:
			return len(obj.Items) > 0
		case *Bool:
			return obj.Value
		}
		return true
	}
}
