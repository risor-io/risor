package object

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cloudcmds/tamarin/v2/op"
)

type SliceIter struct {
	*base
	s         interface{}
	pos       int
	size      int
	current   Object
	converter TypeConverter
}

func (iter *SliceIter) Type() Type {
	return SLICE_ITER
}

func (iter *SliceIter) Inspect() string {
	return fmt.Sprintf("slice_iter(pos=%d size=%d)", iter.pos, iter.size)
}

func (iter *SliceIter) String() string {
	return iter.Inspect()
}

func (iter *SliceIter) Interface() interface{} {
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

func (iter *SliceIter) Equals(other Object) Object {
	if iter == other {
		return True
	}
	return False
}

func (iter *SliceIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "slice_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("slice_iter.next", 0, len(args))
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
			name: "slice_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("slice_iter.entry", 0, len(args))
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

func (iter *SliceIter) IsTruthy() bool {
	return iter.pos < iter.size
}

func (iter *SliceIter) Next() (Object, bool) {
	if iter.pos >= iter.size-1 {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	value := reflect.ValueOf(iter.s).Index(iter.pos).Interface()
	obj, err := iter.converter.From(value)
	if err != nil {
		// This shouldn't happen, but consider what to do here...
		return nil, false
	}
	iter.current = obj
	return iter.current, true
}

func (iter *SliceIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(int64(iter.pos)), iter.current), true
}

func (iter *SliceIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for slice_iter: %v", opType))
}

func NewSliceIter(s interface{}) (*SliceIter, error) {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("type error: cannot create slice_iter (%T given)", s)
	}
	conv, err := NewTypeConverter(typ.Elem())
	if err != nil {
		return nil, err
	}
	return &SliceIter{
		s:         s,
		size:      reflect.ValueOf(s).Len(),
		pos:       -1,
		converter: conv,
	}, nil
}
