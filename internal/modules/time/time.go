package time

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// Name of this module
const Name = "time"

func Now(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.now", 0, args); err != nil {
		return err
	}
	return &object.Time{Value: time.Now()}
}

func Parse(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.parse", 1, args); err != nil {
		return err
	}
	layout, err := AsString(args[0])
	if err != nil {
		return err
	}
	value, err := AsString(args[1])
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
	if err := RequireArgs("time.after", 2, args); err != nil {
		return err
	}
	x, err := AsTime(args[0])
	if err != nil {
		return err
	}
	y, err := AsTime(args[1])
	if err != nil {
		return err
	}
	return ToBoolean(x.After(y))
}

func Before(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.before", 2, args); err != nil {
		return err
	}
	x, err := AsTime(args[0])
	if err != nil {
		return err
	}
	y, err := AsTime(args[1])
	if err != nil {
		return err
	}
	return ToBoolean(x.Before(y))
}

func Format(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.format", 2, args); err != nil {
		return err
	}
	t, err := AsTime(args[0])
	if err != nil {
		return err
	}
	layout, err := AsString(args[1])
	if err != nil {
		return err
	}
	return &object.String{Value: t.Format(layout)}
}

func UTC(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.utc", 1, args); err != nil {
		return err
	}
	t, err := AsTime(args[0])
	if err != nil {
		return err
	}
	return &object.Time{Value: t.UTC()}
}

func Unix(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("time.unix", 1, args); err != nil {
		return err
	}
	t, err := AsTime(args[0])
	if err != nil {
		return err
	}
	return &object.Integer{Value: t.Unix()}
}

func RequireArgs(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		return object.NewError(
			fmt.Sprintf("type error: %s() takes exactly %d argument (%d given)", funcName, count, nArgs))
	}
	return nil
}

func AsString(obj object.Object) (result string, err *object.Error) {
	s, ok := obj.(*object.String)
	if !ok {
		return "", object.NewError("type error: expected a string (got %v)", obj.Type())
	}
	return s.Value, nil
}

func AsTime(obj object.Object) (result time.Time, err *object.Error) {
	s, ok := obj.(*object.Time)
	if !ok {
		return time.Time{}, object.NewError("type error: expected a time (got %v)", obj.Type())
	}
	return s.Value, nil
}

func ToBoolean(value bool) *object.Boolean {
	if value {
		return object.TRUE
	}
	return object.FALSE
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
