package object

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type BSlice struct {
	value []byte
}

func (b *BSlice) Inspect() string {
	return fmt.Sprintf("bslice(%q)", b.value)
}

func (b *BSlice) Type() Type {
	return BSLICE
}

func (b *BSlice) Value() []byte {
	return b.value
}

func (b *BSlice) HashKey() HashKey {
	return HashKey{Type: b.Type(), StrValue: string(b.value)}
}

func (b *BSlice) GetAttr(name string) (Object, bool) {
	switch name {
	case "clone":
		return &Builtin{
			name: "bslice.clone",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("bslice.clone", 0, len(args))
				}
				return b.Clone()
			},
		}, true
	case "equals":
		return &Builtin{
			name: "bslice.equals",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.equals", 1, len(args))
				}
				return b.Equals(args[0])
			},
		}, true
	case "contains":
		return &Builtin{
			name: "bslice.contains",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.contains", 1, len(args))
				}
				return b.Contains(args[0])
			},
		}, true
	case "contains_any":
		return &Builtin{
			name: "bslice.contains_any",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.contains_any", 1, len(args))
				}
				return b.ContainsAny(args[0])
			},
		}, true
	case "contains_rune":
		return &Builtin{
			name: "bslice.contains_rune",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.contains_rune", 1, len(args))
				}
				return b.ContainsRune(args[0])
			},
		}, true
	case "count":
		return &Builtin{
			name: "bslice.count",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.count", 1, len(args))
				}
				return b.Count(args[0])
			},
		}, true
	case "has_prefix":
		return &Builtin{
			name: "bslice.has_prefix",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.has_prefix", 1, len(args))
				}
				return b.HasPrefix(args[0])
			},
		}, true
	case "has_suffix":
		return &Builtin{
			name: "bslice.has_suffix",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.has_suffix", 1, len(args))
				}
				return b.HasSuffix(args[0])
			},
		}, true
	case "index":
		return &Builtin{
			name: "bslice.index",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.index", 1, len(args))
				}
				return b.Index(args[0])
			},
		}, true
	case "index_any":
		return &Builtin{
			name: "bslice.index_any",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.index_any", 1, len(args))
				}
				return b.IndexAny(args[0])
			},
		}, true
	case "index_byte":
		return &Builtin{
			name: "bslice.index_byte",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.index_byte", 1, len(args))
				}
				return b.IndexByte(args[0])
			},
		}, true
	case "index_rune":
		return &Builtin{
			name: "bslice.index_rune",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.index_rune", 1, len(args))
				}
				return b.IndexRune(args[0])
			},
		}, true
	case "repeat":
		return &Builtin{
			name: "bslice.repeat",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("bslice.repeat", 1, len(args))
				}
				return b.Repeat(args[0])
			},
		}, true
	case "replace":
		return &Builtin{
			name: "bslice.replace",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 3 {
					return NewArgsError("bslice.replace", 3, len(args))
				}
				return b.Replace(args[0], args[1], args[2])
			},
		}, true
	case "replace_all":
		return &Builtin{
			name: "bslice.replace_all",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 2 {
					return NewArgsError("bslice.replace_all", 2, len(args))
				}
				return b.ReplaceAll(args[0], args[1])
			},
		}, true
	}
	return nil, false
}

func (b *BSlice) Interface() interface{} {
	return b.value
}

func (b *BSlice) String() string {
	return fmt.Sprintf("bslice(%v)", b.value)
}

func (b *BSlice) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *BSlice:
		return bytes.Compare(b.value, other.value), nil
	case *String:
		return bytes.Compare(b.value, []byte(other.value)), nil
	default:
		return 0, fmt.Errorf("type error: cannot compare bslice to type %s", other.Type())
	}
}

func (b *BSlice) Equals(other Object) Object {
	switch other := other.(type) {
	case *BSlice:
		cmp := bytes.Compare(b.value, other.value)
		if cmp == 0 {
			return True
		}
		return False
	case *String:
		cmp := bytes.Compare(b.value, []byte(other.value))
		if cmp == 0 {
			return True
		}
		return False
	}
	return False
}

func (b *BSlice) IsTruthy() bool {
	return len(b.value) > 0
}

func (b *BSlice) RunOperation(opType op.BinaryOpType, right Object) Object {
	switch right := right.(type) {
	case *BSlice:
		return b.runOperationBytes(opType, right)
	case *String:
		return b.runOperationString(opType, right)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for bslice: %v on type %s", opType, right.Type()))
	}
}

func (b *BSlice) runOperationBytes(opType op.BinaryOpType, right *BSlice) Object {
	switch opType {
	case op.Add:
		result := make([]byte, len(b.value)+len(right.value))
		copy(result, b.value)
		copy(result[len(b.value):], right.value)
		return NewBSlice(result)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for bslice: %v on type %s", opType, right.Type()))
	}
}

func (b *BSlice) runOperationString(opType op.BinaryOpType, right *String) Object {
	switch opType {
	case op.Add:
		rightBytes := []byte(right.value)
		result := make([]byte, len(b.value)+len(rightBytes))
		copy(result, b.value)
		copy(result[len(b.value):], rightBytes)
		return NewBSlice(result)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for bslice: %v on type %s", opType, right.Type()))
	}
}

func (b *BSlice) GetItem(key Object) (Object, *Error) {
	indexObj, ok := key.(*Int)
	if !ok {
		return nil, Errorf("index error: bslice index must be an int (got %s)", key.Type())
	}
	index, err := ResolveIndex(indexObj.value, int64(len(b.value)))
	if err != nil {
		return nil, NewError(err)
	}
	return NewBSlice([]byte{b.value[index]}), nil
}

func (b *BSlice) GetSlice(slice Slice) (Object, *Error) {
	start, stop, err := ResolveIntSlice(slice, int64(len(b.value)))
	if err != nil {
		return nil, NewError(err)
	}
	return NewBSlice(b.value[start:stop]), nil
}

func (b *BSlice) SetItem(key, value Object) *Error {
	indexObj, ok := key.(*Int)
	if !ok {
		return Errorf("index error: index must be an int (got %s)", key.Type())
	}
	index, err := ResolveIndex(indexObj.value, int64(len(b.value)))
	if err != nil {
		return NewError(err)
	}
	data, convErr := AsBytes(value)
	if convErr != nil {
		return convErr
	}
	if len(data) != 1 {
		return Errorf("value error: value must be a single byte (got %d)", len(data))
	}
	b.value[index] = data[0]
	return nil
}

func (b *BSlice) DelItem(key Object) *Error {
	return Errorf("type error: cannot delete from bslice")
}

func (b *BSlice) Contains(obj Object) *Bool {
	data, err := AsInt(obj)
	if err != nil {
		return False
	}
	if data < 0 || data > 255 {
		return False
	}
	return NewBool(bytes.Contains(b.value, []byte{byte(data)}))
}

func (b *BSlice) Len() *Int {
	return NewInt(int64(len(b.value)))
}

func (b *BSlice) Iter() Iterator {
	return NewBytesIter(b)
}

func (b *BSlice) Clone() *BSlice {
	value := make([]byte, len(b.value))
	copy(value, b.value)
	return NewBSlice(value)
}

func (b *BSlice) Reversed() *BSlice {
	value := make([]byte, len(b.value))
	for i := 0; i < len(b.value); i++ {
		value[i] = b.value[len(b.value)-i-1]
	}
	return NewBSlice(value)
}

func (b *BSlice) Integers() []Object {
	result := make([]Object, len(b.value))
	for i, v := range b.value {
		result[i] = NewInt(int64(v))
	}
	return result
}

func (b *BSlice) ContainsAny(obj Object) Object {
	chars, err := AsString(obj)
	if err != nil {
		return err
	}
	return NewBool(bytes.ContainsAny(b.value, chars))
}

func (b *BSlice) ContainsRune(obj Object) Object {
	s, err := AsString(obj)
	if err != nil {
		return err
	}
	if len(s) != 1 {
		return Errorf("bslice.contains_rune: argument must be a single character")
	}
	return NewBool(bytes.ContainsRune(b.value, rune(s[0])))
}

func (b *BSlice) Count(obj Object) Object {
	data, err := AsBytes(obj)
	if err != nil {
		return err
	}
	return NewInt(int64(bytes.Count(b.value, data)))
}

func (b *BSlice) HasPrefix(obj Object) Object {
	data, err := AsBytes(obj)
	if err != nil {
		return err
	}
	return NewBool(bytes.HasPrefix(b.value, data))
}

func (b *BSlice) HasSuffix(obj Object) Object {
	data, err := AsBytes(obj)
	if err != nil {
		return err
	}
	return NewBool(bytes.HasSuffix(b.value, data))
}

func (b *BSlice) Index(obj Object) Object {
	data, err := AsBytes(obj)
	if err != nil {
		return err
	}
	return NewInt(int64(bytes.Index(b.value, data)))
}

func (b *BSlice) IndexAny(obj Object) Object {
	chars, err := AsString(obj)
	if err != nil {
		return err
	}
	return NewInt(int64(bytes.IndexAny(b.value, chars)))
}

func (b *BSlice) IndexByte(obj Object) Object {
	data, err := AsBytes(obj)
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return Errorf("bslice.index_byte: argument must be a single byte")
	}
	return NewInt(int64(bytes.IndexByte(b.value, data[0])))
}

func (b *BSlice) IndexRune(obj Object) Object {
	s, err := AsString(obj)
	if err != nil {
		return err
	}
	if len(s) != 1 {
		return Errorf("bslice.index_rune: argument must be a single character")
	}
	return NewInt(int64(bytes.IndexRune(b.value, rune(s[0]))))
}

func (b *BSlice) Repeat(obj Object) Object {
	count, err := AsInt(obj)
	if err != nil {
		return err
	}
	return NewBSlice(bytes.Repeat(b.value, int(count)))
}

func (b *BSlice) Replace(old, new, count Object) Object {
	oldBytes, err := AsBytes(old)
	if err != nil {
		return err
	}
	newBytes, err := AsBytes(new)
	if err != nil {
		return err
	}
	n, err := AsInt(count)
	if err != nil {
		return err
	}
	return NewBSlice(bytes.Replace(b.value, oldBytes, newBytes, int(n)))
}

func (b *BSlice) ReplaceAll(old, new Object) Object {
	oldBytes, err := AsBytes(old)
	if err != nil {
		return err
	}
	newBytes, err := AsBytes(new)
	if err != nil {
		return err
	}
	return NewBSlice(bytes.ReplaceAll(b.value, oldBytes, newBytes))
}

func (b *BSlice) Cost() int {
	return len(b.value)
}

func NewBSlice(value []byte) *BSlice {
	return &BSlice{value: value}
}
