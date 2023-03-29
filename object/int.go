package object

import (
	"fmt"
)

// Int wraps int64 and implements Object and Hashable interfaces.
type Int struct {
	// value holds the int64 wrapped by this object.
	value int64
}

func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.value)
}

func (i *Int) Type() Type {
	return INT
}

func (i *Int) Value() int64 {
	return i.value
}

func (i *Int) HashKey() HashKey {
	return HashKey{Type: i.Type(), IntValue: i.value}
}

func (i *Int) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (i *Int) Interface() interface{} {
	return i.value
}

func (i *Int) String() string {
	return fmt.Sprintf("int(%v)", i.value)
}

func (i *Int) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Float:
		if float64(i.value) == other.value {
			return 0, nil
		}
		if float64(i.value) > other.value {
			return 1, nil
		}
		return -1, nil
	case *Int:
		if i.value == other.value {
			return 0, nil
		}
		if i.value > other.value {
			return 1, nil
		}
		return -1, nil
	default:
		return CompareTypes(i, other), nil
	}
}

func (i *Int) Equals(other Object) Object {
	switch other := other.(type) {
	case *Int:
		if i.value == other.value {
			return True
		}
	case *Float:
		if float64(i.value) == other.value {
			return True
		}
	}
	return False
}

func (i *Int) IsTruthy() bool {
	return i.value != 0
}

func NewInt(value int64) *Int {
	if value >= 0 && value < tableSize {
		return intCache[value]
	}
	return &Int{value: value}
}

const tableSize = 1024

var intCache = []*Int{}

func init() {
	intCache = make([]*Int, tableSize)
	for i := 0; i < tableSize; i++ {
		intCache[i] = &Int{value: int64(i)}
	}
}
