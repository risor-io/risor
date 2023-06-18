package object

import (
	"context"
	"errors"
	"fmt"
	"image/color"

	"github.com/cloudcmds/tamarin/v2/op"
)

type Color struct {
	c color.Color
}

func (c *Color) Inspect() string {
	return c.String()
}

func (c *Color) Type() Type {
	return COLOR
}

func (c *Color) Value() color.Color {
	return c.c
}

func (c *Color) GetAttr(name string) (Object, bool) {
	switch name {
	case "rgba":
		return &Builtin{
			name: "color.rgba",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("color.rgba", 0, len(args))
				}
				r, g, b, a := c.c.RGBA()
				return NewList([]Object{
					NewInt(int64(r)),
					NewInt(int64(g)),
					NewInt(int64(b)),
					NewInt(int64(a)),
				})
			},
		}, true
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

func (c *Color) Compare(other Object) (int, error) {
	return 0, errors.New("type error: unable to compare colors")
}

func (c *Color) Equals(other Object) Object {
	switch other := other.(type) {
	case *Color:
		if c.c == other.c {
			return True
		}
	}
	return False
}

func (c *Color) IsTruthy() bool {
	return true
}

func (c *Color) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for color: %v ", opType))
}

func NewColor(c color.Color) *Color {
	return &Color{c: c}
}
