package object

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/risor-io/risor/op"
)

type Duration struct {
	value time.Duration
}

func (d *Duration) Type() Type {
	return DURATION
}

func (d *Duration) Value() time.Duration {
	return d.value
}

func (d *Duration) Inspect() string {
	return fmt.Sprintf("duration(%s)", d.value)
}

func (d *Duration) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (d *Duration) Interface() interface{} {
	return d.value
}

func (d *Duration) String() string {
	return d.value.String()
}

func (d *Duration) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Duration:
		if d.value < other.value {
			return -1, nil
		}
		if d.value > other.value {
			return 1, nil
		}
		return 0, nil
	case *Int:
		if d.value < time.Duration(other.value) {
			return -1, nil
		}
		if d.value > time.Duration(other.value) {
			return 1, nil
		}
		return 0, nil
	case *Float:
		if d.value < time.Duration(other.value) {
			return -1, nil
		}
		if d.value > time.Duration(other.value) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("type error: unable to compare duration with %s", other.Type())
	}
}

func (d *Duration) Equals(other Object) Object {
	switch other := other.(type) {
	case *Duration:
		if d.value == other.value {
			return True
		}
	case *Int:
		if d.value == time.Duration(other.value) {
			return True
		}
	}
	return False
}

func (d *Duration) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for duration: %v", opType))
}

func NewDuration(d time.Duration) *Duration {
	return &Duration{value: d}
}

func (d *Duration) IsTruthy() bool {
	return d.value > 0
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.value)
}

func (d *Duration) Cost() int {
	return 0
}

func (d *Duration) SetAttr(name string, value Object) error {
	return fmt.Errorf("attribute error: duration object has no attribute %q", name)
}
