package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type StringIter struct {
	s       *String
	runes   []rune
	pos     int64
	current *String
}

func (iter *StringIter) Type() Type {
	return STRING_ITER
}

func (iter *StringIter) Inspect() string {
	return fmt.Sprintf("string_iter(%s)", iter.s.Inspect())
}

func (iter *StringIter) String() string {
	return iter.Inspect()
}

func (iter *StringIter) Interface() interface{} {
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

func (iter *StringIter) Equals(other Object) Object {
	switch other := other.(type) {
	case *StringIter:
		return NewBool(iter == other)
	default:
		return False
	}
}

func (iter *StringIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "string_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("string_iter.next", 0, len(args))
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
			name: "string_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("string_iter.entry", 0, len(args))
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

func (iter *StringIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.runes))
}

func (iter *StringIter) Next() (Object, bool) {
	if iter.pos >= int64(len(iter.runes)-1) {
		iter.current = nil
		return nil, false
	}
	iter.pos++
	iter.current = NewString(string(iter.runes[iter.pos]))
	return iter.current, true
}

func (iter *StringIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(iter.pos), iter.current), true
}

func (iter *StringIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for string_iter: %v", opType))
}

func (iter *StringIter) Cost() int {
	return 1
}

func NewStringIter(s *String) *StringIter {
	return &StringIter{s: s, runes: []rune(s.value), pos: -1}
}
