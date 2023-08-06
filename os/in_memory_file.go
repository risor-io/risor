package os

import (
	"bytes"
	"errors"
)

type InMemoryFile struct {
	buf *bytes.Buffer
}

func (f *InMemoryFile) Close() error {
	return nil
}

func (f *InMemoryFile) Read(p []byte) (n int, err error) {
	return f.buf.Read(p)
}

func (f *InMemoryFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("io error: read at not supported")
}

func (f *InMemoryFile) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("io error: seek not supported")
}

func (f *InMemoryFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

func (f *InMemoryFile) Stat() (FileInfo, error) {
	return NewFileInfo(GenericFileInfoOpts{
		Name: "",
		Size: int64(f.buf.Len()),
	}), nil
}

func (f *InMemoryFile) Bytes() []byte {
	return f.buf.Bytes()
}

func NewInMemoryFile(data []byte) *InMemoryFile {
	return &InMemoryFile{bytes.NewBuffer(data)}
}
