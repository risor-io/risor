package sql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const SQL_ROW_ITERATOR = object.Type("sql.row_iterator")

type RowIterator struct {
	ctx      context.Context
	rows     *sql.Rows
	once     sync.Once
	closed   chan bool
	isClosed bool
	columns  []string
	current  object.Object
	index    int64
}

func (ri *RowIterator) Type() object.Type {
	return SQL_ROW_ITERATOR
}

func (ri *RowIterator) Inspect() string {
	return "sql.row_iterator()"
}

func (ri *RowIterator) Interface() interface{} {
	return ri.rows
}

func (ri *RowIterator) Equals(other object.Object) object.Object {
	return object.NewBool(ri == other)
}

func (ri *RowIterator) IsTruthy() bool {
	return !ri.isClosed
}

func (ri *RowIterator) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "next":
		return object.NewBuiltin("sql.row_iterator.next", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("sql.row_iterator.next", 0, args); err != nil {
				return err
			}
			obj, ok := ri.Next(ctx)
			if !ok {
				return object.Nil
			}
			return obj
		}), true
	case "close":
		return object.NewBuiltin("sql.row_iterator.close", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("sql.row_iterator.close", 0, args); err != nil {
				return err
			}
			ri.Close()
			return object.Nil
		}), true
	case "entry":
		return object.NewBuiltin("sql.row_iterator.entry", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("sql.row_iterator.entry", 0, args); err != nil {
				return err
			}
			entry, ok := ri.Entry()
			if !ok {
				return object.Nil
			}
			return entry
		}), true
	}
	return nil, false
}

func (ri *RowIterator) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: sql.row_iterator object has no attribute %q", name)
}

func (ri *RowIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for sql.row_iterator: %v", opType)
}

func (ri *RowIterator) Close() {
	ri.once.Do(func() {
		ri.isClosed = true
		ri.rows.Close()
		close(ri.closed)
	})
}

func (ri *RowIterator) Cost() int {
	return 8
}

func (ri *RowIterator) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal sql.row_iterator")
}

// Next implements the object.Iterator interface.
// Advances the iterator and returns the current row and a bool indicating success.
func (ri *RowIterator) Next(ctx context.Context) (object.Object, bool) {
	// Check if there are more rows
	if !ri.rows.Next() {
		ri.current = nil
		if ri.rows.Err() != nil {
			return object.NewError(ri.rows.Err()), true
		}
		return nil, false
	}

	// Get the values for the current row
	rowValues := make([]interface{}, len(ri.columns))
	for i := range rowValues {
		var s interface{}
		rowValues[i] = &s
	}

	if err := ri.rows.Scan(rowValues...); err != nil {
		ri.current = object.NewError(err)
		return ri.current, true
	}

	// Transform the row into a Risor map object
	row := object.NewMap(make(map[string]object.Object))
	for i := range rowValues {
		val := *(rowValues[i].(*interface{}))
		switch val := val.(type) {
		case []byte:
			row.Set(ri.columns[i], object.NewString(string(val)))
		default:
			row.Set(ri.columns[i], object.FromGoType(val))
		}
	}

	ri.index++
	ri.current = row
	return ri.current, true
}

// Entry implements the object.Iterator interface.
// Returns an IteratorEntry for the current row and a bool indicating success.
func (ri *RowIterator) Entry() (object.IteratorEntry, bool) {
	if ri.current == nil {
		return nil, false
	}
	return object.NewEntry(object.NewInt(ri.index), ri.current).
		WithValueAsPrimary(), true
}

func NewRowIterator(ctx context.Context, rows *sql.Rows) *RowIterator {
	columns, _ := rows.Columns()
	return &RowIterator{
		ctx:     ctx,
		rows:    rows,
		closed:  make(chan bool),
		columns: columns,
		index:   -1,
	}
}
