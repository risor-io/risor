package object

import (
	"bytes"
	"fmt"

	"github.com/risor-io/risor/op"
)

type Buffer struct {
	*base
	value *bytes.Buffer
}

func (b *Buffer) Inspect() string {
	return b.String()
}

func (b *Buffer) Type() Type {
	return BUFFER
}

func (b *Buffer) Value() *bytes.Buffer {
	return b.value
}

func (b *Buffer) Interface() interface{} {
	return b.value
}

func (b *Buffer) String() string {
	return fmt.Sprintf("buffer(%q)", b.value.String())
}

func (b *Buffer) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *Buffer:
		return bytes.Compare(b.value.Bytes(), other.value.Bytes()), nil
	case *String:
		return bytes.Compare(b.value.Bytes(), []byte(other.Value())), nil
	case *ByteSlice:
		return bytes.Compare(b.value.Bytes(), other.Value()), nil
	default:
		return 0, fmt.Errorf("type error: cannot compare buffer to type %s", other.Type())
	}
}

func (b *Buffer) Equals(other Object) Object {
	if b == other {
		return True
	}
	return False
}

func (b *Buffer) IsTruthy() bool {
	return b.value.Len() > 0
}

func (b *Buffer) RunOperation(opType op.BinaryOpType, right Object) Object {
	switch right := right.(type) {
	case *Buffer:
		return b.runOperationBytes(opType, right)
	case *String:
		return b.runOperationString(opType, right)
	case *ByteSlice:
		return b.runOperationByteSlice(opType, right)
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for buffer: %v on type %s", opType, right.Type()))
	}
}

func (b *Buffer) runOperationBytes(opType op.BinaryOpType, right *Buffer) Object {
	switch opType {
	case op.Add:
		if _, err := b.value.Write(right.value.Bytes()); err != nil {
			return NewError(err)
		}
		return Nil
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for buffer: %v on type %s", opType, right.Type()))
	}
}

func (b *Buffer) runOperationByteSlice(opType op.BinaryOpType, right *ByteSlice) Object {
	switch opType {
	case op.Add:
		if _, err := b.value.Write(right.value); err != nil {
			return NewError(err)
		}
		return Nil
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for buffer: %v on type %s", opType, right.Type()))
	}
}

func (b *Buffer) runOperationString(opType op.BinaryOpType, right *String) Object {
	switch opType {
	case op.Add:
		if _, err := b.value.Write([]byte(right.value)); err != nil {
			return NewError(err)
		}
		return Nil
	default:
		return NewError(fmt.Errorf("eval error: unsupported operation for buffer: %v on type %s", opType, right.Type()))
	}
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	return b.value.Read(p)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return b.value.Write(p)
}

func (b *Buffer) Cost() int {
	return len(b.value.Bytes())
}

func (b *Buffer) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", b.value.String())), nil
}

func NewBuffer(buf *bytes.Buffer) *Buffer {
	if buf == nil {
		buf = bytes.NewBuffer(nil)
	}
	return &Buffer{value: buf}
}

func NewBufferFromBytes(value []byte) *Buffer {
	return &Buffer{value: bytes.NewBuffer(value)}
}
