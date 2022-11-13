package object

import (
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) Type() Type {
	return TIME_OBJ
}

func (t *Time) Inspect() string {
	return t.Value.Format(time.RFC3339)
}

func (t *Time) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (t *Time) ToInterface() interface{} {
	return t.Value
}

func (t *Time) String() string {
	return t.Inspect()
}
