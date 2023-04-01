package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

type ListIter struct {
	l   *List
	pos int64
}

func (iter *ListIter) Type() Type {
	return LIST_ITER
}

func (iter *ListIter) Inspect() string {
	return fmt.Sprintf("list_iter(%s)", iter.l.Inspect())
}

func (iter *ListIter) Interface() interface{} {
	var entries []map[string]interface{}
	for {
		entry, ok := iter.Next()
		if !ok {
			break
		}
		entries = append(entries, entry.Interface().(map[string]interface{}))
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
			name: "next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("list_iter.next", 0, len(args))
				}
				entry, ok := iter.Next()
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
	return NewError(fmt.Errorf("unsupported operation for list_iter: %v", opType))
}

func (iter *ListIter) Next() (IteratorEntry, bool) {
	items := iter.l.items
	if iter.pos >= int64(len(items)) {
		return nil, false
	}
	r := items[iter.pos]
	entry := NewEntry(NewInt(iter.pos), r)
	iter.pos++
	return entry, true
}

func NewListIter(l *List) *ListIter {
	return &ListIter{l: l, pos: 0}
}
