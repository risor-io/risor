package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

type SetIter struct {
	set  *Set
	keys []HashKey
	pos  int64
}

func (iter *SetIter) Type() Type {
	return SET_ITER
}

func (iter *SetIter) Inspect() string {
	return fmt.Sprintf("set_iter(%s)", iter.set.Inspect())
}

func (iter *SetIter) String() string {
	return iter.Inspect()
}

func (iter *SetIter) Interface() interface{} {
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

func (iter *SetIter) Equals(other Object) Object {
	switch other := other.(type) {
	case *SetIter:
		return NewBool(iter == other)
	default:
		return False
	}
}

func (iter *SetIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("set_iter.next", 0, len(args))
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

func (iter *SetIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.keys))
}

func (iter *SetIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for set_iter: %v", opType))
}

func (iter *SetIter) Next() (IteratorEntry, bool) {
	hashKeys := iter.keys
	if iter.pos >= int64(len(hashKeys)) {
		return nil, false
	}
	key := hashKeys[iter.pos]
	iter.pos++
	value, ok := iter.set.items[key]
	if !ok {
		return nil, false
	}
	return NewEntry(value, True).WithKeyAsPrimary(), true
}

func NewSetIter(set *Set) *SetIter {
	return &SetIter{set: set, keys: set.Keys(), pos: 0}
}
