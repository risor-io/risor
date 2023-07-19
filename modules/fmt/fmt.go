package fmt

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
)

func printableValue(obj object.Object) interface{} {
	switch obj := obj.(type) {
	case *object.String:
		return obj.Value()
	case *object.NilType:
		return nil
	case fmt.Stringer:
		return obj.String()
	default:
		return obj.Inspect()
	}
}

func Printf(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: fmt.printf() takes 1 or more arguments (%d given)", len(args))
	}
	format, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var values []interface{}
	for _, arg := range args[1:] {
		values = append(values, printableValue(arg))
	}
	fmt.Printf(format, values...)
	return object.Nil
}

func Println(ctx context.Context, args ...object.Object) object.Object {
	var values []interface{}
	for _, arg := range args {
		values = append(values, printableValue(arg))
	}
	fmt.Println(values...)
	return object.Nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("fmt", map[string]object.Object{
		"printf":  object.NewBuiltin("print", Printf),
		"println": object.NewBuiltin("println", Println),
	})
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"print":  object.NewBuiltin("print", Println),
		"printf": object.NewBuiltin("printf", Printf),
	}
}
