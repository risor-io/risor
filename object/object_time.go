package object

import (
	"fmt"
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) Type() Type {
	return TIME_OBJ
}

func (t *Time) Inspect() string {
	return fmt.Sprintf("%v", t.Value)
}

func (t *Time) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (t *Time) ToInterface() interface{} {
	return t.Value
}

func (t *Time) String() interface{} {
	return t.Value
}
