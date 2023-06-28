package object

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/risor-io/risor/op"
	tos "github.com/risor-io/risor/os"
)

type File struct {
	*base
	ctx    context.Context
	value  tos.File
	path   string
	once   sync.Once
	closed chan bool
}

func (f *File) Inspect() string {
	return fmt.Sprintf("file(path=%s)", f.path)
}

func (f *File) Type() Type {
	return FILE
}

func (f *File) Value() tos.File {
	return f.value
}

func (f *File) GetAttr(name string) (Object, bool) {
	switch name {
	case "position":
		position, err := f.Position()
		if err != nil {
			return NewError(err), true
		}
		return NewInt(position), true
	case "read":
		return NewBuiltin("file.read", func(ctx context.Context, args ...Object) Object {
			if len(args) != 1 {
				return NewArgsError("file.read", 1, len(args))
			}
			switch obj := args[0].(type) {
			case *ByteSlice:
				n, ioErr := f.Read(obj.Value())
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			case *Buffer:
				n, ioErr := f.Read(obj.Value())
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				return NewInt(int64(n))
			default:
				return Errorf("type error: file.read expects byte_slice or buffer (%s given)", obj.Type())
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
			n, ioErr := f.Write(bytes)
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
			if err := f.Close(); err != nil {
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
			newPosition, ioErr := f.Seek(offset, int(whence))
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
	seeker, ok := f.value.(io.Seeker)
	if !ok {
		return 0, errors.New("value error: this file does not support seeking")
	}
	return seeker.Seek(offset, whence)
}

func (f *File) Write(p []byte) (n int, err error) {
	writer, ok := f.value.(io.Writer)
	if !ok {
		return 0, errors.New("value error: this file does not support writing")
	}
	return writer.Write(p)
}

func (f *File) Close() error {
	var err error
	f.once.Do(func() {
		err = f.value.Close()
		close(f.closed)
	})
	return err
}

func (f *File) Position() (int64, error) {
	seeker, ok := f.value.(io.Seeker)
	if !ok {
		return 0, errors.New("value error: this file does not support seeking")
	}
	position, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	return position, nil
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

func NewFile(ctx context.Context, value tos.File, path string) *File {
	f := &File{
		ctx:    ctx,
		value:  value,
		path:   path,
		closed: make(chan bool),
	}
	f.waitToClose()
	return f
}
