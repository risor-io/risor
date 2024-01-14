package threads

import (
	"context"
	"fmt"
	"sync"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const THREAD object.Type = "thread"

type Thread struct {
	wg       sync.WaitGroup
	callable object.Callable
	args     []object.Object
	res      object.Object
	done     bool
	once     bool
}

func (t *Thread) Type() object.Type {
	return THREAD
}

func (t *Thread) Inspect() string {
	return "thread"
}

func (t *Thread) Interface() interface{} {
	return nil
}

func (t *Thread) IsTruthy() bool {
	return t.res != nil
}

func (t *Thread) Cost() int {
	return 8
}

func (t *Thread) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", THREAD)
}

func (t *Thread) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", THREAD, opType)
}

func (t *Thread) Equals(other object.Object) object.Object {
	if other.Type() != THREAD {
		return object.False
	}
	return object.False
}

func (t *Thread) SetAttr(name string, value object.Object) error {
	switch name {
	case "once":
		v, err := object.AsBool(value)
		if err != nil {
			return err.Value()
		}
		t.done = v
		return nil
	}
	return fmt.Errorf("attribute error: %s object has no attribute %q", THREAD, name)
}

func (t *Thread) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "start", "run":
		name = "threads." + name
		return object.NewBuiltin(name, func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require(name, 0, args); err != nil {
				return err
			}
			if t.once && t.done {
				if t.res == nil {
					return object.Nil
				}
				return t.res
			}
			t.done = true
			t.wg.Add(1)
			go func() {
				defer t.wg.Done()
				t.res = t.callable.Call(ctx, t.args...)
			}()
			return object.Nil
		}), true
	case "wait", "join":
		name = "threads." + name
		return object.NewBuiltin(name, func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require(name, 0, args); err != nil {
				return err
			}
			t.wg.Wait()
			return t.res
		}), true
	case "done":
		name = "threads." + name
		return object.NewBuiltin(name, func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require(name, 0, args); err != nil {
				return err
			}
			return object.NewBool(t.done)
		}), true
	}
	return nil, false
}

func New(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs == 0 {
		return object.Errorf("type error: threads.new() requires at least one argument")
	}

	t := &Thread{}

	switch obj := args[0].(type) {
	case *object.Builtin:
		t.callable = obj
	case *object.Function:
		callFunc, found := object.GetCallFunc(ctx)
		if !found {
			return object.Errorf("eval error: threads.new() context did not contain a call function")
		}
		t.callable = object.NewBuiltin(obj.Name(), func(ctx context.Context, args ...object.Object) object.Object {
			res, err := callFunc(ctx, obj, args)
			if err != nil {
				return object.NewError(err)
			}
			return res
		})
	default:
		return object.Errorf("%s is not callable", args[0].Inspect())
	}

	if numArgs > 1 {
		t.args = append(t.args, args[1:]...)
	}

	return t
}

func Module() *object.Module {
	return object.NewBuiltinsModule("threads", map[string]object.Object{
		"new": object.NewBuiltin("threads.new", New),
	})
}
