package os

import (
	"bytes"
	"errors"
)

// BufferFile is an in memory file backed by a bytes buffer
// Writes to this file are append only and seek is not supported
type BufferFile struct {
	buf *bytes.Buffer
}

func (f *BufferFile) Close() error {
	return nil
}

func (f *BufferFile) Read(p []byte) (n int, err error) {
	return f.buf.Read(p)
}

func (f *BufferFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("io error: read at not supported")
}

func (f *BufferFile) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("io error: seek not supported")
}

func (f *BufferFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

func (f *BufferFile) Stat() (FileInfo, error) {
	return NewFileInfo(GenericFileInfoOpts{
		Name: "",
		Size: int64(f.buf.Len()),
	}), nil
}

func (f *BufferFile) Bytes() []byte {
	return f.buf.Bytes()
}

func NewBufferFile(data []byte) *BufferFile {
	return &BufferFile{bytes.NewBuffer(data)}
}
