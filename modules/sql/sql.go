package sql

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Name of this module
const Name = "sql"

func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: sql.connect() takes exactly one argument (%d given)", len(args))
	}
	url, ok := args[0].(*object.String)
	if !ok {
		return object.NewError("type error: argument to sql.connect not supported, got=%s", args[0].Type())
	}
	conn, err := pgx.Connect(ctx, url.Value)
	if err != nil {
		return &object.Result{Err: &object.Error{Message: err.Error()}}
	}
	return &object.Result{Ok: &object.DatabaseConnection{Conn: conn}}
}

func Query(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError("type error: sql.query() takes at least two arguments (%d given)", len(args))
	}
	var conn *object.DatabaseConnection
	switch arg := args[0].(type) {
	case *object.DatabaseConnection:
		conn = arg
	case *object.Result:
		if !arg.IsOk() {
			return arg
		}
		var ok bool
		conn, ok = arg.Ok.(*object.DatabaseConnection)
		if !ok {
			return object.NewError("type error: argument to sql.query not supported, got=%s", arg.Ok.Type())
		}
	default:
		return object.NewError("type error: argument to sql.query not supported, got=%s", args[0].Type())
	}
	pgxConn := conn.Conn.(*pgx.Conn)
	queryString, ok := args[1].(*object.String)
	if !ok {
		return object.NewError("type error: expected query string, got=%s", args[1].Type())
	}
	var queryArgs []interface{}
	for _, queryArg := range args[2:] {
		queryArgs = append(queryArgs, queryArg.Interface())
	}
	rows, err := pgxConn.Query(ctx, queryString.Value, queryArgs...)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	defer rows.Close()
	fields := rows.FieldDescriptions()
	var results []object.Object
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return object.NewErrorResult(err.Error())
		}
		row := object.NewMap(nil)
		for colIndex, val := range values {
			key := fields[colIndex].Name
			var hashVal object.Object
			if timeVal, ok := val.(pgtype.Time); ok {
				usec := timeVal.Microseconds
				hashVal = object.FromGoType(usec)
			} else {
				hashVal = object.FromGoType(val)
			}
			if hashVal == nil {
				return object.NewErrorResult("type error: unsupported type in sql results: %T", val)
			}
			if !object.IsError(hashVal) {
				row.Items[key] = hashVal
			} else {
				row.Items[key] = object.Nil
			}
		}
		results = append(results, row)
	}
	return &object.Result{Ok: &object.List{Items: results}}
}

// Module returns the `sql` module object
func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := &object.Module{Name: Name, Scope: s}

	if err := s.AddBuiltins([]*object.Builtin{
		{Module: m, Name: "connect", Fn: Connect},
		{Module: m, Name: "query", Fn: Query},
	}); err != nil {
		return nil, err
	}
	return m, nil
}
