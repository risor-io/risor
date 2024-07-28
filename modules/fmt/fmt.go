package fmt

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func printableValue(obj object.Object) interface{} {
	switch obj := obj.(type) {
	// Primitive types have their underlying Go value passed to fmt.Printf
	// so that Go's Printf-style formatting directives work as expected. Also,
	// with these types there's no good reason for the print format to differ.
	case *object.String,
		*object.Int,
		*object.Float,
		*object.Byte,
		*object.Error,
		*object.Bool,
		*object.NilType:
		return obj.Interface()
	// For time objects, as a personal preference, I'm using RFC3339 format
	// rather than Go's default time print format, which I find less readable.
	case *object.Time:
		return obj.Value().Format(time.RFC3339)
	}
	// For everything else, convert the object to a string directly, relying
	// on the object type's String() or Inspect() methods. This gives the author
	// of new types the ability to customize the object print string. Note that
	// Risor map and list objects fall into this category on purpose and the
	// print format for these is intentionally a bit different than the print
	// format for the equivalent Go type (maps and slices).
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
