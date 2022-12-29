package object

import "fmt"

type Entry struct {
	key   Object
	value Object
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

func (e *Entry) IsTruthy() bool {
	return true
}

func (e *Entry) Key() Object {
	return e.key
}

func (e *Entry) Value() Object {
	return e.value
}

func NewEntry(key, value Object) *Entry {
	return &Entry{key: key, value: value}
}
