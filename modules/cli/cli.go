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

func FlagFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cli.flag", 1, args); err != nil {
		return err
	}
	optsMap, objErr := object.AsMap(args[0])
	if objErr != nil {
		return objErr
	}
	opts := optsMap.Value()

	// Infer flag type from the options
	// - If `type` is set, use that
	// - If `value` is set, match its type (string, int, bool, string slice)
	var flag *Flag
	flagValue, hasFlagValue := opts["value"]
	if typ, ok := opts["type"]; ok {
		typStr, ok := typ.(*object.String)
		if !ok {
			return object.Errorf("cli.flag type expected string (got %s)", typ.Type())
		}
		switch typStr.Value() {
		case "string":
			flag = NewFlag(&ucli.StringFlag{})
		case "int":
			flag = NewFlag(&ucli.Int64Flag{})
		case "float":
			flag = NewFlag(&ucli.Float64Flag{})
		case "bool":
			flag = NewFlag(&ucli.BoolFlag{})
		case "string_slice":
			flag = NewFlag(&ucli.StringSliceFlag{})
		case "int_slice":
			flag = NewFlag(&ucli.Int64SliceFlag{})
		case "float_slice":
			flag = NewFlag(&ucli.Float64SliceFlag{})
		default:
			return object.Errorf("unsupported cli.flag type: %s", typStr.Value())
		}
	} else if hasFlagValue {
		switch flagValue.(type) {
		case *object.String:
			flag = NewFlag(&ucli.StringFlag{})
		case *object.Int:
			flag = NewFlag(&ucli.Int64Flag{})
		case *object.Float:
			flag = NewFlag(&ucli.Float64Flag{})
		case *object.Bool:
			flag = NewFlag(&ucli.BoolFlag{})
		}
	}
	if flag == nil {
		return object.Errorf("cli.flag type must be specified")
	}
	if err := setFlagAttrs(flag, opts); err != nil {
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
		"app":     object.NewBuiltin("cli.app", AppFunc),
		"command": object.NewBuiltin("cli.command", CommandFunc),
		"flag":    object.NewBuiltin("cli.flag", FlagFunc),
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
