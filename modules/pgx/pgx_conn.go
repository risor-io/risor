package pgx

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/op"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const PGX_CONN = object.Type("pgx_conn")

type PgxConn struct {
	ctx   context.Context
	conn  *pgx.Conn
	close chan bool
}

func (c *PgxConn) Type() object.Type {
	return PGX_CONN
}

func (c *PgxConn) Inspect() string {
	return "pgx_conn()"
}

func (c *PgxConn) Interface() interface{} {
	return c.conn
}

func (c *PgxConn) Value() *pgx.Conn {
	return c.conn
}

func (c *PgxConn) Equals(other object.Object) object.Object {
	value := other.Type() == PGX_CONN && c.conn == other.(*PgxConn).conn
	return object.NewBool(value)
}

func (c *PgxConn) IsTruthy() bool {
	return true
}

func (c *PgxConn) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "query":
		return object.NewBuiltin("pgx_conn.query", c.Query), true
	case "close":
		return object.NewBuiltin("pgx_conn.close", c.Close), true
	}
	return nil, false
}

func (c *PgxConn) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("unsupported operation for pgx_conn: %v", opType))
}

func (c *PgxConn) Close(ctx context.Context, args ...object.Object) object.Object {
	close(c.close)
	return object.Nil
}

func (c *PgxConn) waitToClose() {
	go func() {
		select {
		case <-c.close:
		case <-c.ctx.Done():
		}
		c.conn.Close(c.ctx)
	}()
}

func New(ctx context.Context, conn *pgx.Conn) *PgxConn {
	obj := &PgxConn{ctx: ctx, conn: conn, close: make(chan bool)}
	obj.waitToClose()
	return obj
}

func (c *PgxConn) Query(ctx context.Context, args ...object.Object) object.Object {

	// The arguments should include a query string and zero or more query args
	if len(args) < 1 {
		return object.Errorf("type error: pgx_conn.query() one or more arguments (%d given)", len(args))
	}
	query, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}

	// Build list of query args as their Go types
	var queryArgs []interface{}
	for _, queryArg := range args[1:] {
		queryArgs = append(queryArgs, queryArg.Interface())
	}

	// Start the query
	rows, err := c.conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return object.NewErrResult(object.NewError(err))
	}
	defer rows.Close()

	// The field descriptions will tell us how to decode the result values
	fields := rows.FieldDescriptions()
	var results []object.Object

	// Transform each result row into a Tamarin map object
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return object.NewErrResult(object.NewError(err))
		}
		row := map[string]object.Object{}
		for colIndex, value := range values {
			key := fields[colIndex].Name
			var val object.Object
			if timeVal, ok := value.(pgtype.Time); ok {
				usec := timeVal.Microseconds
				val = object.FromGoType(usec)
			} else {
				val = object.FromGoType(value)
			}
			if val == nil {
				return object.NewErrResult(
					object.Errorf("type error: pgx_conn.query() encountered unsupported type: %T", value))
			}
			if val != nil && !object.IsError(val) {
				row[key] = val
			} else {
				row[key] = object.NewString(fmt.Sprintf("__error__%s", value))
			}
		}
		results = append(results, object.NewMap(row))
	}
	return object.NewOkResult(object.NewList(results))
}
