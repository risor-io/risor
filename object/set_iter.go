package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

type SetIter struct {
	*base
	set     *Set
	keys    []HashKey
	pos     int64
	current Object
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
				value, ok := iter.Next(ctx)
				if !ok {
					return Nil
				}
				return value
			},
		}, true
	case "entry":
		return &Builtin{
			name: "set_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("set_iter.entry", 0, len(args))
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

func (iter *SetIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.keys))
}

func (iter *SetIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return EvalErrorf("eval error: unsupported operation for set_iter: %v", opType)
}

func (iter *SetIter) Next(ctx context.Context) (Object, bool) {
	hashKeys := iter.keys
	if iter.pos >= int64(len(hashKeys)-1) {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	key := hashKeys[iter.pos]
	value, ok := iter.set.items[key]
	if !ok {
		iter.current = nil
		return nil, false
	}
	iter.current = value
	return value, true
}

func (iter *SetIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(iter.current, True).WithKeyAsPrimary(), true
}

func (iter *SetIter) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal set_iter")
}

func NewSetIter(set *Set) *SetIter {
	return &SetIter{set: set, keys: set.Keys(), pos: -1}
}
