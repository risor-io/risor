package errors

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
)

func New(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: errors.new() takes 1 or more arguments (%d given)", len(args))
	}
	format, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var values []interface{}
	for _, arg := range args[1:] {
		values = append(values, object.PrintableValue(arg))
	}
	return object.NewError(fmt.Errorf(format, values...)).WithRaised(false)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("errors", map[string]object.Object{
		"new": object.NewBuiltin("new", New),
	})
}
