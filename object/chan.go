package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/op"
)

var _ Iterable = (*Chan)(nil)

type Chan struct {
	value        chan Object
	capacity     int
	lastReceived Object
	rxCount      int64
}

func (c *Chan) Type() Type {
	return CHANNEL
}

func (c *Chan) Inspect() string {
	if c.capacity > 0 {
		return fmt.Sprintf("chan(%d)", c.capacity)
	}
	return "chan()"
}

func (c *Chan) Interface() interface{} {
	return c.value
}

func (c *Chan) IsTruthy() bool {
	return true
}

func (c *Chan) Cost() int {
	return 8
}

func (c *Chan) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", CHANNEL)
}

func (c *Chan) RunOperation(opType op.BinaryOpType, right Object) Object {
	return Errorf("eval error: unsupported operation for %s: %v", CHANNEL, opType)
}

func (c *Chan) Equals(other Object) Object {
	return NewBool(c == other)
}

func (c *Chan) SetAttr(name string, value Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", CHANNEL, name)
}

func (c *Chan) GetAttr(name string) (Object, bool) {
	switch name {
	case "send":
		return NewBuiltin("chan.send", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return Errorf("argument error: expected 1 argument, got %d", len(args))
			}
			if err := c.Send(ctx, args[0]); err != nil {
				return NewError(err)
			}
			return Nil
		}), true
	case "receive":
		return NewBuiltin("chan.receive", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return Errorf("argument error: expected 0 arguments, got %d", len(args))
			}
			value, err := c.Receive(ctx)
			if err != nil {
				return NewError(err)
			}
			return value
		}), true
	case "close":
		return NewBuiltin("chan.close", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return Errorf("argument error: expected 0 arguments, got %d", len(args))
			}
			if err := c.Close(); err != nil {
				return NewError(err)
			}
			return Nil
		}), true
	}
	return nil, false
}

func (c *Chan) Close() (err error) {
	// Translate a "close of closed channel" panic to an error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("exec error: %v", r)
		}
	}()
	close(c.value)
	return nil
}

func (c *Chan) Capacity() int {
	return c.capacity
}

func (c *Chan) Next(ctx context.Context) (Object, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case value, ok := <-c.value:
		if !ok {
			return nil, false
		}
		c.lastReceived = value
		c.rxCount++
		return value, true
	}
}

func (c *Chan) Entry() (IteratorEntry, bool) {
	if c.lastReceived != nil {
		return &Entry{
			key:     NewInt(c.rxCount),
			value:   c.lastReceived,
			primary: c.lastReceived,
		}, true
	}
	return nil, false
}

func (c *Chan) Iter() Iterator {
	return c
}

func (c *Chan) Send(ctx context.Context, value Object) (err error) {
	// Translate a "send on closed channel" panic to an error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("exec error: %v", r)
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.value <- value:
		return nil
	}
}

func (c *Chan) Receive(ctx context.Context) (Object, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case value, ok := <-c.value:
		if !ok {
			// TODO: return zero value of the underlying channel type, if any
			return Nil, nil
		}
		return value, nil
	}
}

func (c *Chan) Value() chan Object {
	return c.value
}

func NewChan(size int) *Chan {
	return &Chan{
		capacity: size,
		value:    make(chan Object, size),
	}
}
