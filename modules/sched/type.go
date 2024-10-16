package sched

import (
	"context"

	"codnect.io/chrono"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type task struct {
	Name string
	t    chrono.ScheduledTask
}

// IsTruthy returns true if the task is not cancelled.
func (t *task) IsTruthy() bool {
	return !t.t.IsCancelled()
}

// Cost returns the cost of the task.
func (t *task) Cost() int {
	return 0
}

// Equals returns true if the task is equal to the other object.
func (t *task) Equals(other object.Object) object.Object {
	so, ok := other.(*task)
	if !ok {
		return object.False
	}
	ok = (*t == *so)
	return object.NewBool(ok)
}

// Inspect returns the string representation of the task.
func (t *task) Inspect() string {
	return "sched.task"
}

// Type returns the type of the task.
func (t *task) Type() object.Type {
	return object.Type("sched.task")
}

func (t *task) Interface() any {
	return t.t
}

// RunOperation returns a type error for unsupported operations.
func (t *task) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for job: %v", opType)
}

// GetAttr returns the attribute of the task.
func (t *task) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "cancel":
		return object.NewBuiltin("sched.task.stop", func(ctx context.Context, args ...object.Object) object.Object {
			t.t.Cancel()
			return nil
		}), true
	case "is_running":
		return object.NewBuiltin("sched.task.is_running", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewBool(!t.t.IsCancelled())
		}), true
	}
	return nil, false
}

func (t *task) SetAttr(name string, value object.Object) error {
	return errz.TypeErrorf("type error: object has no attribute %q", name)
}
