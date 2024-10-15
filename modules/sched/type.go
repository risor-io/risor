//go:build sched
// +build sched

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

func (t *task) IsTruthy() bool {
	return true
}

func (t *task) Cost() int {
	return 0
}

func (t *task) Equals(other object.Object) object.Object {
	so, ok := other.(*task)
	if !ok {
		return object.False
	}
	ok = (*t == *so)
	return object.NewBool(ok)
}

func (t *task) Inspect() string {
	return "sched.task"
}

func (t *task) Type() object.Type {
	return object.Type("sched.task")
}

func (t *task) Interface() any {
	return t
}

func (t *task) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for job: %v", opType)
}

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
