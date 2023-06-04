package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type MapIter struct {
	m    *Map
	keys []string
	pos  int64
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
			name: "next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("map_iter.next", 0, len(args))
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

func (iter *MapIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.keys))
}

func (iter *MapIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for map_iter: %v", opType))
}

func (iter *MapIter) Next() (IteratorEntry, bool) {
	keys := iter.keys
	if iter.pos >= int64(len(keys)) {
		return nil, false
	}
	key := keys[iter.pos]
	iter.pos++
	value, ok := iter.m.items[key]
	if !ok {
		return nil, false
	}
	return NewEntry(NewString(key), value).WithKeyAsPrimary(), true
}

func NewMapIter(m *Map) *MapIter {
	return &MapIter{m: m, keys: m.SortedKeys(), pos: 0}
}
