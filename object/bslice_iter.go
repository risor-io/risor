package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type BSliceIter struct {
	*base
	b       *BSlice
	pos     int64
	current Object
}

func (iter *BSliceIter) Type() Type {
	return BSLICE_ITER
}

func (iter *BSliceIter) Inspect() string {
	return fmt.Sprintf("bslice_iter(%s)", iter.b.Inspect())
}

func (iter *BSliceIter) String() string {
	return iter.Inspect()
}

func (iter *BSliceIter) Interface() interface{} {
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

func (iter *BSliceIter) Equals(other Object) Object {
	switch other := other.(type) {
	case *BSliceIter:
		return NewBool(iter == other)
	default:
		return False
	}
}

func (iter *BSliceIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "bslice_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("bslice_iter.next", 0, len(args))
				}
				value, ok := iter.Next()
				if !ok {
					return Nil
				}
				return value
			},
		}, true
	case "entry":
		return &Builtin{
			name: "bslice_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("bslice_iter.entry", 0, len(args))
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

func (iter *BSliceIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.b.value))
}

func (iter *BSliceIter) Next() (Object, bool) {
	data := iter.b.value
	if iter.pos >= int64(len(data)-1) {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	value := data[iter.pos]
	iter.current = NewBSlice([]byte{value})
	return iter.current, true
}

func (iter *BSliceIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(iter.pos), iter.current), true
}

func (iter *BSliceIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for bslice_iter: %v", opType))
}

func NewBytesIter(b *BSlice) *BSliceIter {
	return &BSliceIter{b: b, pos: -1}
}
