package carbon

import (
	"context"

	"github.com/golang-module/carbon/v2"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func Now(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("carbon.now", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 1 {
		tz, err := object.AsString(args[0])
		if err != nil {
			return err
		}
		return NewCarbon(carbon.Now(tz))
	}
	return NewCarbon(carbon.Now())
}

func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("carbon.carbon", 0, 2, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return NewCarbon(carbon.Now())
	}
	switch arg := args[0].(type) {
	case *object.String:
		var tz []string
		if len(args) == 2 {
			tzArg, err := object.AsString(args[1])
			if err != nil {
				return err
			}
			tz = []string{tzArg}
		}
		c := carbon.Parse(arg.Value(), tz...)
		if c.IsInvalid() {
			return object.Errorf("value error: invalid time string (got %q)", arg.Value())
		}
		return NewCarbon(c)
	case *object.Time:
		return NewCarbon(carbon.CreateFromStdTime(arg.Value()))
	default:
		return object.Errorf("type error: expected string or time (got %s)", arg.Type())
	}
}

func Yesterday(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("carbon.yesterday", 0, 1, args); err != nil {
		return err
	}
	var tz []string
	if len(args) == 1 {
		tzArg, err := object.AsString(args[0])
		if err != nil {
			return err
		}
		tz = []string{tzArg}
	}
	return NewCarbon(carbon.Yesterday(tz...))
}

func Tomorrow(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("carbon.tomorrow", 0, 1, args); err != nil {
		return err
	}
	var tz []string
	if len(args) == 1 {
		tzArg, err := object.AsString(args[0])
		if err != nil {
			return err
		}
		tz = []string{tzArg}
	}
	return NewCarbon(carbon.Tomorrow(tz...))
}

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("carbon.parse", 1, 2, args); err != nil {
		return err
	}
	carbonStr, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var tz []string
	if len(args) == 2 {
		tzArg, err := object.AsString(args[1])
		if err != nil {
			return err
		}
		tz = []string{tzArg}
	}
	c := carbon.Parse(carbonStr, tz...)
	if c.IsInvalid() {
		return object.Errorf("value error: invalid time string (got %q)", carbonStr)
	}
	return NewCarbon(c)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("carbon", map[string]object.Object{
		"now":       object.NewBuiltin("now", Now),
		"yesterday": object.NewBuiltin("yesterday", Yesterday),
		"tomorrow":  object.NewBuiltin("tomorrow", Tomorrow),
		"parse":     object.NewBuiltin("parse", Parse),
	}, Create)
}
