package os

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudcmds/tamarin/v2/object"
)

func Exit(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: exit() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		os.Exit(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		os.Exit(int(obj.Value()))
	case *object.Error:
		os.Exit(1)
	}
	return object.Errorf("type error: exit() argument must be an int or error (%s given)", args[0].Type())
}

func Printf(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: printf() takes 1 or more arguments (%d given)", len(args))
	}
	format, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var values []interface{}
	for _, arg := range args[1:] {
		switch arg := arg.(type) {
		case *object.String:
			values = append(values, arg.Value())
		default:
			values = append(values, arg.Interface())
		}
	}
	fmt.Printf(format, values...)
	return object.Nil
}

func Print(ctx context.Context, args ...object.Object) object.Object {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			values[i] = arg.Value()
		default:
			values[i] = arg.Inspect()
		}
	}
	fmt.Println(values...)
	return object.Nil
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"exit":   object.NewBuiltin("exit", Exit),
		"print":  object.NewBuiltin("print", Print),
		"printf": object.NewBuiltin("printf", Printf),
	}
}
