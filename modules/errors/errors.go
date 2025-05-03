package errors

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/risor-io/risor/object"
)

func getFormatAndValues(args []object.Object) (string, []interface{}, *object.Error) {
	numArgs := len(args)
	if numArgs == 0 {
		return "", nil, nil
	}
	format, err := object.AsString(args[0])
	if err != nil {
		return "", nil, err
	}
	var values []interface{}
	for _, arg := range args[1:] {
		values = append(values, object.PrintableValue(arg))
	}
	return format, values, nil
}

func New(ctx context.Context, args ...object.Object) object.Object {
	format, values, err := getFormatAndValues(args)
	if err != nil {
		return err
	}
	return object.NewError(fmt.Errorf(format, values...)).WithRaised(false)
}

func TypeError(ctx context.Context, args ...object.Object) object.Object {
	format, values, err := getFormatAndValues(args)
	if err != nil {
		return err
	}
	return object.TypeErrorf(format, values...).WithRaised(false)
}

func EvalError(ctx context.Context, args ...object.Object) object.Object {
	format, values, err := getFormatAndValues(args)
	if err != nil {
		return err
	}
	return object.EvalErrorf(format, values...).WithRaised(false)
}

func ArgsError(ctx context.Context, args ...object.Object) object.Object {
	format, values, err := getFormatAndValues(args)
	if err != nil {
		return err
	}
	return object.ArgsErrorf(format, values...).WithRaised(false)
}

func As(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.ArgsErrorf("args error: errors.as() takes exactly 2 arguments (%d given)", len(args))
	}
	err, typErr := object.AsError(args[0])
	if typErr != nil {
		return typErr
	}
	other, typErr := object.AsError(args[1])
	if typErr != nil {
		return typErr
	}
	otherType := reflect.TypeOf(other.Value())
	otherInst := reflect.New(otherType).Interface()
	return object.NewBool(errors.As(err.Value(), otherInst))
}

func Is(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.ArgsErrorf("args error: errors.is() takes exactly 2 arguments (%d given)", len(args))
	}
	err, typErr := object.AsError(args[0])
	if typErr != nil {
		return typErr
	}
	target, typErr := object.AsError(args[1])
	if typErr != nil {
		return typErr
	}
	return object.NewBool(errors.Is(err.Value(), target.Value()))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("errors", map[string]object.Object{
		"new":        object.NewBuiltin("new", New),
		"type_error": object.NewBuiltin("type_error", TypeError),
		"eval_error": object.NewBuiltin("eval_error", EvalError),
		"args_error": object.NewBuiltin("args_error", ArgsError),
		"as":         object.NewBuiltin("as", As),
		"is":         object.NewBuiltin("is", Is),
	})
}
