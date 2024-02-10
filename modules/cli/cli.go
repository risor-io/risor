package cli

import (
	"context"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	ucli "github.com/urfave/cli/v2"
)

func AppFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.app", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	app, err := NewApp(opts)
	if err != nil {
		return object.NewError(err)
	}
	return app
}

func NewStringFlag(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.string_flag", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	flag := NewFlag(&ucli.StringFlag{})
	if err := setFlagAttrs(flag, opts.Value()); err != nil {
		return object.NewError(err)
	}
	return flag
}

func NewIntFlag(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.int_flag", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	flag := NewFlag(&ucli.Int64Flag{})
	if err := setFlagAttrs(flag, opts.Value()); err != nil {
		return object.NewError(err)
	}
	return flag
}

func NewBoolFlag(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.bool_flag", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	flag := NewFlag(&ucli.BoolFlag{})
	if err := setFlagAttrs(flag, opts.Value()); err != nil {
		return object.NewError(err)
	}
	return flag
}

func NewStringSliceFlag(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.string_slice_flag", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	flag := NewFlag(&ucli.StringSliceFlag{})
	if err := setFlagAttrs(flag, opts.Value()); err != nil {
		return object.NewError(err)
	}
	return flag
}

func CommandFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.command", 1, args); err != nil {
		return err
	}
	opts, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	cmd := NewCommand(&ucli.Command{})
	for k, v := range opts.Value() {
		if err := cmd.SetAttr(k, v); err != nil {
			return object.NewError(err)
		}
	}
	return cmd
}

func Module() *object.Module {
	return object.NewBuiltinsModule("cli", map[string]object.Object{
		"app":               object.NewBuiltin("cli.app", AppFunc),
		"command":           object.NewBuiltin("cli.command", CommandFunc),
		"string_flag":       object.NewBuiltin("cli.string_flag", NewStringFlag),
		"int_flag":          object.NewBuiltin("cli.int_flag", NewIntFlag),
		"bool_flag":         object.NewBuiltin("cli.bool_flag", NewBoolFlag),
		"string_slice_flag": object.NewBuiltin("cli.string_slice_flag", NewStringSliceFlag),
	})
}

func setFlagAttrs(flag *Flag, attrs map[string]object.Object) error {
	for k, v := range attrs {
		if err := flag.SetAttr(k, v); err != nil {
			return err
		}
	}
	return nil
}
