package time

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// Name of this module
const Name = "time"

func Now(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.now", 0, args); err != nil {
		return err
	}
	return object.NewTime(time.Now())
}

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.parse", 2, args); err != nil {
		return err
	}
	layout, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	value, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	t, parseErr := time.Parse(layout, value)
	if parseErr != nil {
		return object.NewErrResult(object.NewError(parseErr))
	}
	return object.NewOkResult(object.NewTime(t))
}

func After(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.after", 2, args); err != nil {
		return err
	}
	x, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	y, err := object.AsTime(args[1])
	if err != nil {
		return err
	}
	return object.NewBool(x.After(y))
}

func Before(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.before", 2, args); err != nil {
		return err
	}
	x, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	y, err := object.AsTime(args[1])
	if err != nil {
		return err
	}
	return object.NewBool(x.Before(y))
}

func Format(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.format", 2, args); err != nil {
		return err
	}
	t, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	layout, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewString(t.Format(layout))
}

func UTC(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.utc", 1, args); err != nil {
		return err
	}
	t, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	return object.NewTime(t.UTC())
}

func Unix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.unix", 1, args); err != nil {
		return err
	}
	t, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	return object.NewInt(t.Unix())
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("now", Now, m),
		object.NewBuiltin("parse", Parse, m),
		object.NewBuiltin("after", After, m),
		object.NewBuiltin("before", Before, m),
		object.NewBuiltin("format", Format, m),
		object.NewBuiltin("utc", UTC, m),
		object.NewBuiltin("unix", Unix, m),
	}); err != nil {
		return nil, err
	}
	return m, nil
}
