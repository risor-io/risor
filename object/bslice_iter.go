package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type BSliceIter struct {
	b   *BSlice
	pos int64
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
			name: "next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("bslice_iter.next", 0, len(args))
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

func (iter *BSliceIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.b.value))
}

func (iter *BSliceIter) Next() (IteratorEntry, bool) {
	data := iter.b.value
	if iter.pos >= int64(len(data)) {
		return nil, false
	}
	r := data[iter.pos]
	entry := NewEntry(NewInt(iter.pos), NewBSlice([]byte{r}))
	iter.pos++
	return entry, true
}

func (iter *BSliceIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for bslice_iter: %v", opType))
}

func NewBytesIter(b *BSlice) *BSliceIter {
	return &BSliceIter{b: b, pos: 0}
}
