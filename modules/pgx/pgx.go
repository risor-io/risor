package pgx

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/jackc/pgx/v5"
)

// Name of this module
const Name = "pgx"

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
		return object.NewErrResult(object.NewError(err))
	}
	return object.NewOkResult(New(ctx, conn))
}

// Module returns the `pgx` module object
func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("connect", Connect, m),
	}); err != nil {
		return nil, err
	}
	return m, nil
}
