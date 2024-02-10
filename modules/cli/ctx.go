package cli

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	ucli "github.com/urfave/cli/v2"
)

const CTX object.Type = "cli.ctx"

type Ctx struct {
	value    *ucli.Context
	callFunc object.CallFunc
}

func (c *Ctx) Type() object.Type {
	return "cli.ctx"
}

func (c *Ctx) Inspect() string {
	return "cli.ctx"
}

func (c *Ctx) Interface() interface{} {
	return c.value
}

func (c *Ctx) IsTruthy() bool {
	return true
}

func (c *Ctx) Cost() int {
	return 0
}

func (c *Ctx) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", CTX)
}

func (c *Ctx) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", CTX, opType)
}

func (c *Ctx) Equals(other object.Object) object.Object {
	if other.Type() != "cli.ctx" {
		return object.False
	}
	return object.NewBool(c.value == other.(*Ctx).value)
}

func (c *Ctx) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", CTX, name)
}

func (c *Ctx) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "args":
		return object.NewBuiltin("cli.ctx.args", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewStringList(c.value.Args().Slice())
		}), true
	case "narg":
		return object.NewBuiltin("cli.ctx.narg", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewInt(int64(c.value.NArg()))
		}), true
	case "value":
		return object.NewBuiltin("cli.ctx.value", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("cli.ctx.value", 1, args); err != nil {
				return err
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			val := c.value.Value(name)
			return object.FromGoType(val)
		}), true
	case "count":
		return object.NewBuiltin("cli.ctx.count", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("cli.ctx.count", 1, args); err != nil {
				return err
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Count(name)))
		}), true
	case "flag_names":
		return object.NewBuiltin("cli.ctx.flag_names", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewStringList(c.value.FlagNames())
		}), true
	case "local_flag_names":
		return object.NewBuiltin("cli.ctx.local_flag_names", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewStringList(c.value.LocalFlagNames())
		}), true
	case "is_set":
		return object.NewBuiltin("cli.ctx.is_set", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("cli.ctx.is_set", 1, args); err != nil {
				return err
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return object.NewBool(c.value.IsSet(name))
		}), true
	case "set":
		return object.NewBuiltin("cli.ctx.set", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("cli.ctx.set", 2, args); err != nil {
				return err
			}
			name, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			value, err := object.AsString(args[1])
			if err != nil {
				return err
			}
			if err := c.value.Set(name, value); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "num_flags":
		return object.NewBuiltin("cli.ctx.num_flags", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewInt(int64(c.value.NumFlags()))
		}), true
	}
	return nil, false
}

func NewContext(c *ucli.Context, callFunc object.CallFunc) *Ctx {
	return &Ctx{c, callFunc}
}