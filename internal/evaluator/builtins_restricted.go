package evaluator

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/object"
)

// output a string to stdout
func printFun(ctx context.Context, args ...object.Object) object.Object {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		values[i] = arg.Inspect()
	}
	fmt.Println(values...)
	return object.NULL
}

// printfFun is the implementation of our `printf` function.
func printfFun(ctx context.Context, args ...object.Object) object.Object {
	// Convert to the formatted version, via our `sprintf` function
	out := sprintfFun(ctx, args...)
	// If that returned a string then we can print it
	if out.Type() == object.STRING_OBJ {
		fmt.Print(out.(*object.String).Value)
	}
	return object.NULL
}

// RegisterRestrictedBuiltins adds builtins that should not be present when
// running server side
func RegisterRestrictedBuiltins() {
	RegisterBuiltin("print", printFun)
	RegisterBuiltin("printf", printfFun)
}
