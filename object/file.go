package object

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cloudcmds/tamarin/v2/op"
)

type File struct {
	value *os.File
}

func (f *File) Inspect() string {
	return fmt.Sprintf("file(name=%q)", f.value.Name())
}

func (f *File) Type() Type {
	return FILE
}

func (f *File) Value() *os.File {
	return f.value
}

func (f *File) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(f.value.Name()), true
	case "position":
		position, _ := f.value.Seek(0, io.SeekCurrent)
		return NewInt(int64(position)), true
	case "sync":
		return NewBuiltin("file.sync", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return NewArgsError("file.sync", 0, len(args))
			}
			if err := f.value.Sync(); err != nil {
				return NewError(err)
			}
			return Nil
		}), true
	case "read":
		return NewBuiltin("file.read", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return NewArgsError("file.read", 1, len(args))
			}
			bufferSize, err := AsInt(args[0])
			if err != nil {
				return err
			}
			buffer := make([]byte, bufferSize)
			n, ioErr := f.value.Read(buffer)
			if ioErr != nil && ioErr != io.EOF {
				return NewError(ioErr)
			}
			return NewBSlice(buffer[:n])
		}), true
	case "write":
		return NewBuiltin("file.write", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return NewArgsError("file.write", 1, len(args))
			}
			switch obj := args[0].(type) {
			case *BSlice:
				n, ioErr := f.value.Write(obj.Value())
				if ioErr != nil {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			case *String:
				n, ioErr := f.value.WriteString(obj.Value())
				if ioErr != nil {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			default:
				return NewError(errors.New("type error: file.write expects bytes or string"))
			}
		}), true
	case "close":
		return NewBuiltin("file.close", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return NewArgsError("file.close", 0, len(args))
			}
			if err := f.value.Close(); err != nil {
				return NewError(err)
			}
			return Nil
		}), true
	case "seek":
		return NewBuiltin("file.seek", func(ctx context.Context, args ...Object) Object {
			if len(args) != 2 {
				return NewArgsError("file.seek", 2, len(args))
			}
			offset, err := AsInt(args[0])
			if err != nil {
				return err
			}
			whence, err := AsInt(args[1])
			if err != nil {
				return err
			}
			newPosition, ioErr := f.value.Seek(offset, int(whence))
			if ioErr != nil {
				return NewError(ioErr)
			}
			return NewInt(newPosition)
		}), true
	}
	return nil, false
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.value.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.value.Seek(offset, whence)
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.value.Write(p)
}

func (f *File) Close() error {
	return f.value.Close()
}

func (f *File) Interface() interface{} {
	return f.value
}

func (f *File) String() string {
	return fmt.Sprintf("file(%v)", f.value)
}

func (f *File) Compare(other Object) (int, error) {
	return 0, errors.New("type error: unable to compare files")
}

func (f *File) Equals(other Object) Object {
	switch other := other.(type) {
	case *File:
		if f.value == other.value {
			return True
		}
	}
	return False
}

func (f *File) IsTruthy() bool {
	return true
}

func (f *File) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for file: %v ", opType))
}

func NewFile(value *os.File) *File {
	return &File{value: value}
}
