package object

import (
	"fmt"
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) Type() Type {
	return TIME
}

func (t *Time) Inspect() string {
	return t.Value.Format(time.RFC3339)
}

func (t *Time) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (t *Time) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (t *Time) ToInterface() interface{} {
	return t.Value
}

func (t *Time) String() string {
	return fmt.Sprintf("time(%s)", t.Inspect())
}

func (t *Time) Compare(other Object) (int, error) {
	typeComp := CompareTypes(t, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Time)
	if t.Value == otherStr.Value {
		return 0, nil
	}
	if t.Value.After(otherStr.Value) {
		return 1, nil
	}
	return -1, nil
}

func (t *Time) Equals(other Object) Object {
	if other.Type() == TIME && t.Value == other.(*Time).Value {
		return True
	}
	return False
}

func NewTime(t time.Time) *Time {
	return &Time{Value: t}
}
