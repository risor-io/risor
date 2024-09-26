package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

type MapIter struct {
	*base
	m       *Map
	keys    []string
	pos     int64
	current *String
}

func (iter *MapIter) Type() Type {
	return MAP_ITER
}

func (iter *MapIter) Inspect() string {
	return fmt.Sprintf("map_iter(%s)", iter.m.Inspect())
}

func (iter *MapIter) String() string {
	return iter.Inspect()
}

func (iter *MapIter) Interface() interface{} {
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

func (iter *MapIter) Equals(other Object) Object {
	switch other := other.(type) {
	case *MapIter:
		return NewBool(iter == other)
	default:
		return False
	}
}

func (iter *MapIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "map_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map_iter.next", 0, len(args))
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
			name: "map_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map_iter.entry", 0, len(args))
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

func (iter *MapIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.keys))
}

func (iter *MapIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return EvalErrorf("eval error: unsupported operation for map_iter: %v", opType)
}

func (iter *MapIter) Next(ctx context.Context) (Object, bool) {
	keys := iter.keys
	if iter.pos >= int64(len(keys)-1) {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	iter.current = NewString(keys[iter.pos])
	return iter.current, true
}

func (iter *MapIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	value, ok := iter.m.items[iter.current.value]
	if !ok {
		iter.current = nil
		return nil, false
	}
	return NewEntry(iter.current, value).WithKeyAsPrimary(), true
}

func (iter *MapIter) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal map_iter")
}

func NewMapIter(m *Map) *MapIter {
	return &MapIter{m: m, keys: m.SortedKeys(), pos: -1}
}
