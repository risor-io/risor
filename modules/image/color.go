package image

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/op"
)

type Color struct {
	c color.Color
}

func (c *Color) Inspect() string {
	return c.String()
}

func (c *Color) Type() object.Type {
	return "image.color"
}

func (c *Color) Value() color.Color {
	return c.c
}

func (c *Color) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "rgba":
		return object.NewBuiltin("color.rgba",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("color.rgba", 0, len(args))
				}
				r, g, b, a := c.c.RGBA()
				return object.NewList([]object.Object{
					object.NewInt(int64(r)),
					object.NewInt(int64(g)),
					object.NewInt(int64(b)),
					object.NewInt(int64(a)),
				})
			}), true
	}
	return nil, false
}

func (c *Color) Interface() interface{} {
	return c.c
}

func (c *Color) String() string {
	r, g, b, a := c.c.RGBA()
	return fmt.Sprintf("color(r=%d g=%d b=%d a=%d)", r, g, b, a)
}

func (c *Color) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare colors")
}

func (c *Color) Equals(other object.Object) object.Object {
	switch other := other.(type) {
	case *Color:
		if c.c == other.c {
			return object.True
		}
	}
	return object.False
}

func (c *Color) IsTruthy() bool {
	return true
}

func (c *Color) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for color: %v ", opType))
}

func (c *Color) Cost() int {
	return 8
}

func NewColor(c color.Color) *Color {
	return &Color{c: c}
}
