package pgx

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const PGX_ROW_ITERATOR = object.Type("pgx.row_iterator")

type RowIterator struct {
	ctx      context.Context
	rows     pgx.Rows
	once     sync.Once
	closed   chan bool
	isClosed bool
	fields   []pgconn.FieldDescription
	current  object.Object
	index    int64
}

func (ri *RowIterator) Type() object.Type {
	return PGX_ROW_ITERATOR
}

func (ri *RowIterator) Inspect() string {
	return "pgx.row_iterator()"
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
		return object.NewBuiltin("pgx.row_iterator.next", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("pgx.row_iterator.next", 0, args); err != nil {
				return err
			}
			obj, ok := ri.Next(ctx)
			if !ok {
				return object.Nil
			}
			return obj
		}), true
	case "close":
		return object.NewBuiltin("pgx.row_iterator.close", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("pgx.row_iterator.close", 0, args); err != nil {
				return err
			}
			ri.Close()
			return object.Nil
		}), true
	case "entry":
		return object.NewBuiltin("pgx.row_iterator.entry", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("pgx.row_iterator.entry", 0, args); err != nil {
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
	return object.TypeErrorf("type error: pgx.row_iterator object has no attribute %q", name)
}

func (ri *RowIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for pgx.row_iterator: %v", opType)
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
	return nil, errz.TypeErrorf("type error: unable to marshal pgx.row_iterator")
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
	values, err := ri.rows.Values()
	if err != nil {
		ri.current = object.NewError(err)
		return ri.current, true
	}

	// Transform the row into a Risor map object
	row := map[string]object.Object{}
	for colIndex, value := range values {
		key := ri.fields[colIndex].Name
		var val object.Object
		if timeVal, ok := value.(pgtype.Time); ok {
			usec := timeVal.Microseconds
			val = object.FromGoType(usec)
		} else {
			val = object.FromGoType(value)
		}
		if val == nil {
			ri.current = object.TypeErrorf("type error: pgx.row_iterator.next() encountered unsupported type: %T", value)
			return ri.current, true
		}
		if !object.IsError(val) {
			row[key] = val
		} else {
			row[key] = object.NewString(fmt.Sprintf("__error__%s", value))
		}
	}

	ri.index++
	ri.current = object.NewMap(row)
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

func NewRowIterator(ctx context.Context, rows pgx.Rows) *RowIterator {
	fields := rows.FieldDescriptions()
	return &RowIterator{
		ctx:    ctx,
		rows:   rows,
		closed: make(chan bool),
		fields: fields,
		index:  -1,
	}
}
