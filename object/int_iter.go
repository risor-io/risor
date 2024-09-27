package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/op"
)

type IntIter struct {
	*base
	done      bool
	target    int64
	pos       int64
	increment int64
	current   *Int
}

func (iter *IntIter) Type() Type {
	return INT_ITER
}

func (iter *IntIter) Inspect() string {
	return fmt.Sprintf("int_iter(%d)", iter.target)
}

func (iter *IntIter) String() string {
	return iter.Inspect()
}

func (iter *IntIter) Interface() interface{} {
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

func (iter *IntIter) Equals(other Object) Object {
	if iter == other {
		return True
	}
	return False
}

func (iter *IntIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "int_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("int_iter.next", 0, len(args))
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
			name: "int_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("int_iter.entry", 0, len(args))
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

func (iter *IntIter) IsTruthy() bool {
	return !iter.done
}

func (iter *IntIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return TypeErrorf("type error: unsupported operation for int_iter: %v", opType)
}

func (iter *IntIter) Next(ctx context.Context) (Object, bool) {
	if iter.done {
		return nil, false
	}
	absTarget := iter.target
	if absTarget < 0 {
		absTarget = -absTarget
	}
	if iter.pos >= absTarget-1 {
		iter.done = true
		return nil, false
	}
	iter.pos++
	iter.current = NewInt(iter.pos * iter.increment)
	return iter.current, true
}

func (iter *IntIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(iter.pos), iter.current), true
}

func (iter *IntIter) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal int_iter")
}

func NewIntIter(i *Int) *IntIter {
	increment := int64(1)
	if i.value < 0 {
		increment = -1
	}
	return &IntIter{target: i.value, increment: increment, pos: -1}
}
