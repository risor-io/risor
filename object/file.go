package object

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/cloudcmds/tamarin/v2/op"
)

type File struct {
	*base
	ctx    context.Context
	value  *os.File
	once   sync.Once
	closed chan bool
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
			switch obj := args[0].(type) {
			case *BSlice:
				n, ioErr := f.value.Read(obj.Value())
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			case *Buffer:
				n, ioErr := f.value.Read(obj.Value())
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			default:
				return Errorf("type error: file.read expects bslice or buffer (%s given)", obj.Type())
			}
		}), true
	case "write":
		return NewBuiltin("file.write", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return NewArgsError("file.write", 1, len(args))
			}
			bytes, err := AsBytes(args[0])
			if err != nil {
				return err
			}
			n, ioErr := f.value.Write(bytes)
			if ioErr != nil {
				return NewError(ioErr)
			}
			return NewInt(int64(n))
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
	var err error
	f.once.Do(func() {
		err = f.value.Close()
		close(f.closed)
	})
	return err
}

func (f *File) waitToClose() {
	go func() {
		select {
		case <-f.closed:
		case <-f.ctx.Done():
			f.value.Close()
		}
	}()
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
	if f == other {
		return True
	}
	return False
}

func (f *File) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for file: %v ", opType))
}

func (f *File) Cost() int {
	return 8
}

func NewFile(ctx context.Context, value *os.File) *File {
	f := &File{
		ctx:    ctx,
		value:  value,
		closed: make(chan bool),
	}
	f.waitToClose()
	return f
}
