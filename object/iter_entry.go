package object

import (
	"fmt"

	"github.com/risor-io/risor/op"
)

type Entry struct {
	*base
	key     Object
	value   Object
	primary Object
}

func (e *Entry) Type() Type {
	return ITER_ENTRY
}

func (e *Entry) Inspect() string {
	return fmt.Sprintf("iter_entry(%s, %s)", e.key.Inspect(), e.value.Inspect())
}

func (e *Entry) Interface() interface{} {
	return map[string]interface{}{
		"key":   e.key.Interface(),
		"value": e.value.Interface(),
	}
}

func (e *Entry) Equals(other Object) Object {
	switch other := other.(type) {
	case *Entry:
		if e.key.Equals(other.key) != True {
			return False
		}
		if e.value.Equals(other.value) != True {
			return False
		}
		return True
	default:
		return False
	}
}

func (e *Entry) GetAttr(name string) (Object, bool) {
	switch name {
	case "key":
		return e.key, true
	case "value":
		return e.value, true
	}
	return nil, false
}

func (e *Entry) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for entry: %v", opType))
}

func (e *Entry) Key() Object {
	return e.key
}

func (e *Entry) Value() Object {
	return e.value
}

func (e *Entry) Primary() Object {
	if e.primary != nil {
		return e.primary
	}
	return e.value
}

func (e *Entry) WithKeyAsPrimary() *Entry {
	e.primary = e.key
	return e
}

func (e *Entry) WithValueAsPrimary() *Entry {
	e.primary = e.value
	return e
}

func NewEntry(key, value Object) *Entry {
	return &Entry{key: key, value: value}
}
