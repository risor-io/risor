package object

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sync"

	"github.com/risor-io/risor/op"
	ros "github.com/risor-io/risor/os"
)

type File struct {
	*base
	ctx    context.Context
	value  ros.File
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

func (f *File) Value() ros.File {
	return f.value
}

func (f *File) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewBuiltin("file.name", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return NewArgsError("file.name", 0, len(args))
			}
			return NewString(filepath.Base(f.path))
		}), true
	case "stat":
		return NewBuiltin("file.stat", func(ctx context.Context, args ...Object) Object {
			if len(args) != 0 {
				return NewArgsError("file.stat", 0, len(args))
			}
			info, err := f.value.Stat()
			if err != nil {
				return NewError(err)
			}
			return NewFileInfo(info)
		}), true
	case "position":
		position, err := f.Position()
		if err != nil {
			return NewError(err), true
		}
		return NewInt(position), true
	case "read":
		return NewBuiltin("file.read", func(ctx context.Context, args ...Object) Object {
			if len(args) > 1 {
				return NewArgsRangeError("file.read", 0, 1, len(args))
			}
			if len(args) == 0 {
				bytes, err := io.ReadAll(f.value)
				if err != nil {
					return NewError(err)
				}
				return NewByteSlice(bytes)
			}
			switch obj := args[0].(type) {
			case *ByteSlice:
				slice := obj.Value()
				n, ioErr := f.Read(slice)
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				if n == len(slice) {
					return obj
				}
				return NewByteSlice(slice[:n])
			case *Buffer:
				stat, err := f.value.Stat()
				if err != nil {
					return NewError(err)
				}
				size := stat.Size()
				if size > math.MaxInt32 {
					return NewError(errors.New("file.read: file size exceeds maximum int32"))
				}
				buf := obj.Value()
				buf.Grow(int(size)) // review: this can panic
				n, ioErr := f.Read(buf.Bytes())
				if ioErr != nil && ioErr != io.EOF {
					return NewError(ioErr)
				}
				return NewByteSlice(buf.Bytes()[:n])
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
	case "read_lines":
		return NewBuiltin("file.read_lines", func(ctx context.Context, args ...Object) Object {
			if len(args) > 0 {
				return NewArgsError("file.read_lines", 0, len(args))
			}
			var lines []Object
			scanner := bufio.NewScanner(f.value)
			for scanner.Scan() {
				lines = append(lines, NewString(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return NewError(err)
			}
			return NewList(lines)
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
	return f.Inspect()
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

func (f *File) MarshalJSON() ([]byte, error) {
	return nil, errors.New("type error: unable to marshal file")
}

func (f *File) Iter() Iterator {
	return NewFileIter(f)
}

func NewFile(ctx context.Context, value ros.File, path string) *File {
	f := &File{
		ctx:    ctx,
		value:  value,
		path:   path,
		closed: make(chan bool),
	}
	f.waitToClose()
	return f
}
