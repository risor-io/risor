package color

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/os"
)

var _ object.Object = (*Color)(nil)

const COLOR object.Type = "color.color"

type Color struct {
	value *color.Color
}

func (c *Color) IsTruthy() bool {
	return true
}

func (c *Color) Type() object.Type {
	return COLOR
}

func (c *Color) Inspect() string {
	return fmt.Sprintf("%s()", COLOR)
}

func (c *Color) Value() *color.Color {
	return c.value
}

func (c *Color) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set %q on %s object", name, COLOR)
}

func (c *Color) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "sprintf":
		return object.NewBuiltin("sprintf", func(ctx context.Context, args ...object.Object) object.Object {
			numArgs := len(args)
			if numArgs < 1 {
				return object.NewArgsRangeError("color.sprintf", 1, 64, numArgs)
			}
			format, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			var items []interface{}
			for i := 1; i < numArgs; i++ {
				items = append(items, args[i].Interface())
			}
			return object.NewString(c.value.Sprintf(format, items...))
		}), true
	case "fprintf":
		return object.NewBuiltin("fprintf", func(ctx context.Context, args ...object.Object) object.Object {
			numArgs := len(args)
			if numArgs < 2 {
				return object.NewArgsRangeError("color.fprintf", 2, 64, numArgs)
			}
			writer, err := object.AsWriter(args[0])
			if err != nil {
				return err
			}
			format, err := object.AsString(args[1])
			if err != nil {
				return err
			}
			var items []interface{}
			for i := 2; i < numArgs; i++ {
				items = append(items, args[i].Interface())
			}
			c.value.Fprintf(writer, format, items...)
			return object.Nil
		}), true
	case "printf":
		return object.NewBuiltin("printf", func(ctx context.Context, args ...object.Object) object.Object {
			numArgs := len(args)
			if numArgs < 1 {
				return object.TypeErrorf("type error: color.printf() takes 1 or more arguments (%d given)", len(args))
			}
			format, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			var values []interface{}
			for _, arg := range args[1:] {
				values = append(values, arg.Interface())
			}
			stdout := os.GetDefaultOS(ctx).Stdout()
			if _, ioErr := c.value.Fprintf(stdout, format, values...); ioErr != nil {
				return object.Errorf("io error: %v", ioErr)
			}
			return object.Nil
		}), true
	case "print":
		return object.NewBuiltin("print", func(ctx context.Context, args ...object.Object) object.Object {
			var values []interface{}
			for _, arg := range args {
				values = append(values, arg.Interface())
			}
			stdout := os.GetDefaultOS(ctx).Stdout()
			if _, ioErr := c.value.Fprintln(stdout, values...); ioErr != nil {
				return object.Errorf("io error: %v", ioErr)
			}
			return object.Nil
		}), true
	default:
		return nil, false
	}
}

func (c *Color) Interface() interface{} {
	return c.value
}

func (c *Color) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Color) Cost() int {
	return 0
}

func (c *Color) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", COLOR, opType)
}

func NewColor(v *color.Color) *Color {
	return &Color{value: v}
}
