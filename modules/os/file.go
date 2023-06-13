package os

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/op"
)

// File wraps os.File
type File struct {
	value *os.File
}

func (f *File) Inspect() string {
	return fmt.Sprintf("file(name=%q)", f.value.Name())
}

func (f *File) Type() object.Type {
	return "file"
}

func (f *File) Value() *os.File {
	return f.value
}

func (f *File) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "name":
		return object.NewString(f.value.Name()), true
	case "position":
		position, _ := f.value.Seek(0, io.SeekCurrent)
		return object.NewInt(int64(position)), true
	case "sync":
		return object.NewBuiltin("file.sync", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.NewArgsError("file.sync", 0, len(args))
			}
			if err := f.value.Sync(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "read":
		return object.NewBuiltin("file.read", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewArgsError("file.read", 1, len(args))
			}
			bufferSize, err := object.AsInt(args[0])
			if err != nil {
				return err
			}
			buffer := make([]byte, bufferSize)
			n, ioErr := f.value.Read(buffer)
			if ioErr != nil && ioErr != io.EOF {
				return object.NewError(ioErr)
			}
			return object.NewBSlice(buffer[:n])
		}), true
	case "write":
		return object.NewBuiltin("file.write", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewArgsError("file.write", 1, len(args))
			}
			switch obj := args[0].(type) {
			case *object.BSlice:
				n, ioErr := f.value.Write(obj.Value())
				if ioErr != nil {
					return object.NewError(ioErr)
				}
				return object.NewInt(int64(n))
			case *object.String:
				n, ioErr := f.value.WriteString(obj.Value())
				if ioErr != nil {
					return object.NewError(ioErr)
				}
				return object.NewInt(int64(n))
			default:
				return object.NewError(errors.New("type error: file.write expects bytes or string"))
			}
		}), true
	case "close":
		return object.NewBuiltin("file.close", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.NewArgsError("file.close", 0, len(args))
			}
			if err := f.value.Close(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "seek":
		return object.NewBuiltin("file.seek", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.NewArgsError("file.seek", 2, len(args))
			}
			offset, err := object.AsInt(args[0])
			if err != nil {
				return err
			}
			whence, err := object.AsInt(args[1])
			if err != nil {
				return err
			}
			newPosition, ioErr := f.value.Seek(offset, int(whence))
			if ioErr != nil {
				return object.NewError(ioErr)
			}
			return object.NewInt(newPosition)
		}), true
	}
	return nil, false
}

func (f *File) Interface() interface{} {
	return f.value
}

func (f *File) String() string {
	return fmt.Sprintf("file(%v)", f.value)
}

func (f *File) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare files")
}

func (f *File) Equals(other object.Object) object.Object {
	switch other := other.(type) {
	case *File:
		if f.value == other.value {
			return object.True
		}
	}
	return object.False
}

func (f *File) IsTruthy() bool {
	return true
}

func (f *File) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("unsupported operation for file: %v ", opType))
}

func NewFile(value *os.File) *File {
	return &File{value: value}
}
