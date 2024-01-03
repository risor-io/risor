package sql

import (
	"context"

	"github.com/risor-io/risor/object"
)

func Connect(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)

	if numArgs < 1 {
		return object.NewArgsError("sql.connect", 1, numArgs)
	}

	connStr, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	db, connErr := New(ctx, connStr)
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
