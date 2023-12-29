package pgx

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/risor-io/risor/object"
)

func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.Errorf("type error: pgx.connect() takes exactly one argument (%d given)", len(args))
	}
	url, ok := args[0].(*object.String)
	if !ok {
		return object.Errorf("type error: pgx.connect() expected a string argument (got %s)", args[0].Type())
	}
	conn, err := pgx.Connect(ctx, url.Value())
	if err != nil {
		return object.NewError(err)
	}
	return New(ctx, conn)
}

//go:embed pgx.md
var docs string

// Module returns the `pgx` module object
func Module() *object.Module {
	return object.NewBuiltinsModule("pgx", map[string]object.Object{
		"connect": object.NewBuiltin("connect", Connect),
	}).WithDocstring(docs)
}
