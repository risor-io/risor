package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

type ListIter struct {
	*base
	l       *List
	pos     int64
	current Object
}

func (iter *ListIter) Type() Type {
	return LIST_ITER
}

func (iter *ListIter) Inspect() string {
	return fmt.Sprintf("list_iter(%s)", iter.l.Inspect())
}

func (iter *ListIter) String() string {
	return iter.Inspect()
}

func (iter *ListIter) Interface() interface{} {
	ctx := context.Background()
	var entries []any
	for {
		entry, ok := iter.Next(ctx)
		if !ok {
			break
		}
		entries = append(entries, entry.Interface())
	}
	return entries
}

func (iter *ListIter) Equals(other Object) Object {
	switch other := other.(type) {
	case *ListIter:
		return NewBool(iter == other)
	default:
		return False
	}
}

func (iter *ListIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "list_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list_iter.next", 0, len(args))
				}
				value, ok := iter.Next(ctx)
				if !ok {
					return Nil
				}
				return value
			},
		}, true
	case "entry":
		return &Builtin{
			name: "list_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list_iter.entry", 0, len(args))
				}
				entry, ok := iter.Entry()
				if !ok {
					return Nil
				}
				return entry
			},
		}, true
	}
	return nil, false
}

func (iter *ListIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.l.items))
}

func (iter *ListIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return EvalErrorf("eval error: unsupported operation for list_iter: %v", opType)
}

func (iter *ListIter) Next(ctx context.Context) (Object, bool) {
	items := iter.l.items
	if iter.pos >= int64(len(items)-1) {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	iter.current = items[iter.pos]
	return iter.current, true
}

func (iter *ListIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(iter.pos), iter.current), true
}

func (iter *ListIter) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal list_iter")
}

func NewListIter(l *List) *ListIter {
	return &ListIter{l: l, pos: -1}
}
