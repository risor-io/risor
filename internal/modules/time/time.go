package time

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// Name of this module
const Name = "time"

func Now(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.now", 0, args); err != nil {
		return err
	}
	return &object.Time{Value: time.Now()}
}

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.parse", 1, args); err != nil {
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
		return object.NewErrorResult(parseErr.Error())
	}
	return &object.Result{Ok: &object.Time{Value: t}}
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
	return object.NewBoolean(x.After(y))
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
	return object.NewBoolean(x.Before(y))
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
	return &object.String{Value: t.Format(layout)}
}

func UTC(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.utc", 1, args); err != nil {
		return err
	}
	t, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	return &object.Time{Value: t.UTC()}
}

func Unix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("time.unix", 1, args); err != nil {
		return err
	}
	t, err := object.AsTime(args[0])
	if err != nil {
		return err
	}
	return &object.Integer{Value: t.Unix()}
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "now", Func: Now},
		{Name: "parse", Func: Parse},
		{Name: "after", Func: After},
		{Name: "before", Func: Before},
		{Name: "format", Func: Format},
		{Name: "utc", Func: UTC},
		{Name: "unix", Func: Unix},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
