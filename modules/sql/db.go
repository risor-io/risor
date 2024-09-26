package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	// TODO: we can add more drivers from this list:
	//  https://github.com/xo/dburl?tab=readme-ov-file#database-schemes-aliases-and-drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/xo/dburl"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const DB_CONN object.Type = "db.conn"

type DB struct {
	conn   *sql.DB
	once   sync.Once
	closed chan bool
}

func (db *DB) Type() object.Type {
	return DB_CONN
}

func (db *DB) Inspect() string {
	return "sql.conn"
}

func (db *DB) Interface() interface{} {
	return db.conn
}

func (db *DB) IsTruthy() bool {
	return db.conn != nil
}

func (db *DB) Cost() int {
	return 8
}

func (db *DB) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal db.conn")
}

func (db *DB) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", DB_CONN, opType)
}

func (db *DB) Equals(other object.Object) object.Object {
	if other.Type() != DB_CONN {
		return object.False
	}
	return object.NewBool(db.conn == other.(*DB).conn)
}

func (db *DB) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", DB_CONN, name)
}

func (db *DB) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "query":
		return object.NewBuiltin("sql.query", db.Query), true
	case "exec":
		return object.NewBuiltin("sql.exec", db.Exec), true
	case "close":
		return object.NewBuiltin("sql.close", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("sql.close", 0, args); err != nil {
				return err
			}
			if err := db.Close(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	}
	return nil, false
}

func (db *DB) Exec(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: sql.exec() requires at least one argument")
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
	_, err := db.conn.Exec(query, queryArgs...)
	if err != nil {
		return object.NewError(err)
	}

	return nil
}

func (db *DB) Query(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: sql.query() requires at least one argument")
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
	rows, err := db.conn.Query(query, queryArgs...)
	if err != nil || rows.Err() != nil {
		return object.Errorf("failed to query db: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return object.Errorf("failed to get columns: %w", err)
	}

	rowList := object.NewList(make([]object.Object, 0))
	for rows.Next() {
		rowValues := make([]interface{}, len(columns))
		for i := range rowValues {
			var s interface{}
			rowValues[i] = &s
		}
		if err := rows.Scan(rowValues...); err != nil {
			return object.NewError(err)
		}

		row := object.NewMap(make(map[string]object.Object))
		for i := range rowValues {
			val := *(rowValues[i].(*interface{}))
			switch val := val.(type) {
			case []byte:
				row.Set(columns[i], object.NewString(string(val)))
			default:
				row.Set(columns[i], object.FromGoType(val))
			}
		}

		rowList.Append(row)
	}

	return rowList
}

func (db *DB) Close() error {
	var err error
	db.once.Do(func() {
		err = db.conn.Close()
		close(db.closed)
	})
	return err
}

func (db *DB) waitToClose(ctx context.Context) {
	go func() {
		select {
		case <-db.closed:
		case <-ctx.Done():
			_ = db.conn.Close()
		}
	}()
}

func New(ctx context.Context, connection string) (*DB, error) {
	db, err := dburl.Open(connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	obj := &DB{
		conn:   db,
		closed: make(chan bool),
	}
	obj.waitToClose(ctx)
	return obj, nil
}
