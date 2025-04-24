package sql

import (
	"context"

	"github.com/risor-io/risor/object"
)

func Connect(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)

	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsError("sql.connect", 1, numArgs)
	}

	connStr, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	// Default options
	stream := false

	// Check for options map as second argument
	if numArgs == 2 {
		optMap, ok := args[1].(*object.Map)
		if !ok {
			return object.TypeErrorf("type error: sql.connect() second argument must be a map (got %s)", args[1].Type())
		}

		// Process stream option if provided
		if streamVal, found := optMap.Value()["stream"]; found {
			streamBool, ok := streamVal.(*object.Bool)
			if !ok {
				return object.TypeErrorf("type error: sql.connect() 'stream' option must be a boolean (got %s)", streamVal.Type())
			}
			stream = streamBool.Value()
		}
	}

	db, connErr := New(ctx, connStr, stream)
	if connErr != nil {
		return object.NewError(connErr)
	}

	return db
}

func Module() *object.Module {
	return object.NewBuiltinsModule("sql", map[string]object.Object{
		"connect": object.NewBuiltin("sql.connect", Connect),
	})
}
