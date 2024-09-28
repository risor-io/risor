package object

import (
	"context"
	"fmt"
	"time"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

var _ Object = (*Context)(nil)

type Context struct {
	value  context.Context
	cancel context.CancelFunc
}

func (c *Context) Type() Type {
	return CONTEXT
}

func (c *Context) Inspect() string {
	return fmt.Sprintf("context(%p)", c.value)
}

func (c *Context) Interface() interface{} {
	return c.value
}

func (c *Context) IsTruthy() bool {
	return true
}

func (c *Context) Cost() int {
	return 8
}

func (c *Context) Value() context.Context {
	return c.value
}

func (c *Context) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Context) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal %s", CONTEXT)
}

func (c *Context) RunOperation(opType op.BinaryOpType, right Object) Object {
	return TypeErrorf("type error: unsupported operation for %s: %v", CONTEXT, opType)
}

func (c *Context) Equals(other Object) Object {
	return NewBool(c == other)
}

func (c *Context) SetAttr(name string, value Object) error {
	return TypeErrorf("type error: %s object has no attribute %q", CONTEXT, name)
}

func (c *Context) GetAttr(name string) (Object, bool) {
	switch name {
	case "done":
		return NewBuiltin("context.done", func(ctx context.Context, args ...Object) Object {
			select {
			case <-c.value.Done():
				return True
			default:
				return False
			}
		}), true
	case "err":
		return NewBuiltin("context.err", func(ctx context.Context, args ...Object) Object {
			if err := c.value.Err(); err != nil {
				return NewError(err)
			}
			return Nil
		}), true
	case "value":
		return NewBuiltin("context.value", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return Errorf("argument error: expected 1 argument, got %d", len(args))
			}
			key, ok := args[0].(*String)
			if !ok {
				return TypeErrorf("type error: context.value() key must be a string")
			}
			value := c.value.Value(getContextKey(key.Value()))
			if value == nil {
				return Nil
			}
			obj, ok := value.(Object)
			if !ok {
				return TypeErrorf("type error: context.value() value must be an object")
			}
			return obj
		}), true
	case "cancel":
		return NewBuiltin("context.cancel", func(ctx context.Context, args ...Object) Object {
			if c.cancel != nil {
				c.cancel()
			}
			return Nil
		}), true
	}
	return nil, false
}

func NewContext(ctx context.Context, cancel context.CancelFunc) *Context {
	return &Context{
		value:  ctx,
		cancel: cancel,
	}
}

func ContextBackground() *Context {
	return NewContext(context.Background(), nil)
}

func ContextWithCancel(parent *Context) *Context {
	ctx, cancel := context.WithCancel(parent.value)
	return NewContext(ctx, cancel)
}

func ContextWithTimeout(parent *Context, timeout float64) *Context {
	ctx, cancel := context.WithTimeout(parent.value, time.Duration(timeout*float64(time.Second)))
	return NewContext(ctx, cancel)
}

func ContextWithDeadline(parent *Context, deadline time.Time) *Context {
	ctx, cancel := context.WithDeadline(parent.value, deadline)
	return NewContext(ctx, cancel)
}

func ContextWithValue(parent *Context, key string, value Object) *Context {
	ctx := context.WithValue(parent.value, getContextKey(key), value)
	return NewContext(ctx, parent.cancel)
}

func getContextKey(key string) any {
	return contextKey(fmt.Sprintf("risor:ctx:%s", key))
}
