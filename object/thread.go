package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

type Thread struct {
	callable Callable
	args     []Object
	done     chan bool
	result   Object
}

func (t *Thread) Type() Type {
	return THREAD
}

func (t *Thread) Inspect() string {
	switch obj := t.callable.(type) {
	case Object:
		return fmt.Sprintf("thread(%s)", obj.Inspect())
	default:
		return "thread()"
	}
}

func (t *Thread) Interface() interface{} {
	return nil
}

func (t *Thread) IsTruthy() bool {
	return true
}

func (t *Thread) Cost() int {
	return 8
}

func (t *Thread) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal %s", THREAD)
}

func (t *Thread) RunOperation(opType op.BinaryOpType, right Object) Object {
	return Errorf("eval error: unsupported operation for %s: %v", THREAD, opType)
}

func (t *Thread) Equals(other Object) Object {
	return NewBool(t == other)
}

func (t *Thread) SetAttr(name string, value Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", THREAD, name)
}

func (t *Thread) GetAttr(name string) (Object, bool) {
	switch name {
	case "wait":
		return NewBuiltin("thread.wait", func(ctx context.Context, args ...Object) Object {
			// Wait for the thread to finish or the context to be cancelled
			return t.Wait(ctx)
		}), true
	}
	return nil, false
}

func (t *Thread) Wait(ctx context.Context) Object {
	select {
	case <-ctx.Done():
		return EvalErrorf("eval error: %s", ctx.Err())
	case <-t.done:
		return t.result
	}
}

func NewThread(ctx context.Context, callable Callable, args []Object) *Thread {
	if callable == nil {
		panic("callable is nil")
	}

	t := &Thread{
		callable: callable,
		args:     args,
		done:     make(chan bool),
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.result = NewError(fmt.Errorf("panic: %v", r))
			}
			close(t.done)
		}()
		t.result = callable.Call(ctx, args...)
	}()

	return t
}
