package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/risor-io/risor/object"
)

func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.TypeErrorf("type error: pgx.connect() takes exactly one argument (%d given)", len(args))
	}
	url, ok := args[0].(*object.String)
	if !ok {
		return object.TypeErrorf("type error: pgx.connect() expected a string argument (got %s)", args[0].Type())
	}
	conn, err := pgx.Connect(ctx, url.Value())
	if err != nil {
		return object.NewError(err)
	}
	return New(ctx, conn)
}

// Module returns the `pgx` module object
func Module() *object.Module {
	return object.NewBuiltinsModule("pgx", map[string]object.Object{
		"connect": object.NewBuiltin("connect", Connect),
	})
}
