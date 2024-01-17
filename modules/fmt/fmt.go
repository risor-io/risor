package fmt

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func printableValue(obj object.Object) interface{} {
	iface := obj.Interface()
	if iface != nil {
		return iface
	}
	switch obj := obj.(type) {
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
	stdout := os.GetDefaultOS(ctx).Stdout()
	if _, ioErr := fmt.Fprintf(stdout, format, values...); ioErr != nil {
		return object.Errorf("io error: %v", ioErr)
	}
	return object.Nil
}

func Println(ctx context.Context, args ...object.Object) object.Object {
	var values []interface{}
	for _, arg := range args {
		values = append(values, printableValue(arg))
	}
	stdout := os.GetDefaultOS(ctx).Stdout()
	if _, ioErr := fmt.Fprintln(stdout, values...); ioErr != nil {
		return object.Errorf("io error: %v", ioErr)
	}
	return object.Nil
}

//go:embed fmt.md
var docs string

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"print":  object.NewBuiltin("print", Println),
		"printf": object.NewBuiltin("printf", Printf),
	}
}

func Module() *object.Module {
	return object.NewBuiltinsModule("fmt", map[string]object.Object{
		"printf":  object.NewBuiltin("print", Printf),
		"println": object.NewBuiltin("println", Println),
	})
}
