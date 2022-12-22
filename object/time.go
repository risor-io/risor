package object

import (
	"fmt"
	"time"
)

type Time struct {
	value time.Time
}

func (t *Time) Type() Type {
	return TIME
}

func (t *Time) Value() time.Time {
	return t.value
}

func (t *Time) Inspect() string {
	return t.value.Format(time.RFC3339)
}

func (t *Time) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (t *Time) Interface() interface{} {
	return t.value
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
	if t.value == otherStr.value {
		return 0, nil
	}
	if t.value.After(otherStr.value) {
		return 1, nil
	}
	return -1, nil
}

func (t *Time) Equals(other Object) Object {
	if other.Type() == TIME && t.value == other.(*Time).value {
		return True
	}
	return False
}

func NewTime(t time.Time) *Time {
	return &Time{value: t}
}
