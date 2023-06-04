package object

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type StringIter struct {
	s     *String
	runes []rune
	pos   int64
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
			name: "next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("string_iter.next", 0, len(args))
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

func (iter *StringIter) IsTruthy() bool {
	return iter.pos < int64(len(iter.runes))
}

func (iter *StringIter) Next() (IteratorEntry, bool) {
	if iter.pos >= int64(len(iter.runes)) {
		return nil, false
	}
	r := iter.runes[iter.pos]
	entry := NewEntry(NewInt(iter.pos), NewString(string(r)))
	iter.pos++
	return entry, true
}

func (iter *StringIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for string_iter: %v", opType))
}

func NewStringIter(s *String) *StringIter {
	return &StringIter{s: s, runes: []rune(s.value), pos: 0}
}
