package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/risor-io/risor/object"
)

func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.TypeErrorf("type error: pgx.connect() takes one or two arguments (%d given)", len(args))
	}
	url, ok := args[0].(*object.String)
	if !ok {
		return object.TypeErrorf("type error: pgx.connect() expected a string argument (got %s)", args[0].Type())
	}

	// Default options
	stream := false

	// Check for options map as second argument
	if len(args) == 2 {
		optMap, ok := args[1].(*object.Map)
		if !ok {
			return object.TypeErrorf("type error: pgx.connect() second argument must be a map (got %s)", args[1].Type())
		}

		// Process stream option if provided
		// When stream is true, query() returns a row iterator instead of loading all rows at once
		if streamVal, found := optMap.Value()["stream"]; found {
			streamBool, ok := streamVal.(*object.Bool)
			if !ok {
				return object.TypeErrorf("type error: pgx.connect() 'stream' option must be a boolean (got %s)", streamVal.Type())
			}
			stream = streamBool.Value()
		}
	}

	conn, err := pgx.Connect(ctx, url.Value())
	if err != nil {
		return object.NewError(err)
	}
	return New(ctx, conn, stream)
}

// Module returns the `pgx` module object
func Module() *object.Module {
	return object.NewBuiltinsModule("pgx", map[string]object.Object{
		"connect": object.NewBuiltin("connect", Connect),
	})
}
