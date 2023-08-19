package object

import (
	"bufio"
	"context"
	"fmt"

	"github.com/risor-io/risor/op"
)

type FileIter struct {
	*base
	f       *File
	pos     int64
	done    bool
	scanner *bufio.Scanner
	current Object
}

func (iter *FileIter) Type() Type {
	return FILE_ITER
}

func (iter *FileIter) Inspect() string {
	return fmt.Sprintf("file_iter(%s)", iter.f.Inspect())
}

func (iter *FileIter) String() string {
	return iter.Inspect()
}

func (iter *FileIter) Interface() interface{} {
	var entries []any
	for {
		entry, ok := iter.Next()
		if !ok {
			break
		}
		entries = append(entries, entry.Interface())
	}
	return entries
}

func (iter *FileIter) Equals(other Object) Object {
	if iter == other {
		return True
	}
	return False
}

func (iter *FileIter) GetAttr(name string) (Object, bool) {
	switch name {
	case "next":
		return &Builtin{
			name: "file_iter.next",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("file_iter.next", 0, len(args))
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
			name: "file_iter.entry",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("file_iter.entry", 0, len(args))
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

func (iter *FileIter) IsTruthy() bool {
	return !iter.done
}

func (iter *FileIter) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for file_iter: %v", opType))
}

func (iter *FileIter) Next() (Object, bool) {
	if iter.done {
		return nil, false
	}
	if result := iter.scanner.Scan(); !result { // review: can panic
		iter.done = true
		return nil, false
	}
	iter.current = NewString(iter.scanner.Text())
	iter.pos++
	return iter.current, true
}

func (iter *FileIter) Entry() (IteratorEntry, bool) {
	if iter.current == nil {
		return nil, false
	}
	return NewEntry(NewInt(iter.pos), iter.current), true
}

func (iter *FileIter) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal file_iter")
}

func NewFileIter(f *File) *FileIter {
	return &FileIter{f: f, pos: -1, scanner: bufio.NewScanner(f)}
}
