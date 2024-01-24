package os

import (
	"errors"
	"io"
)

// InMemoryFile is an in-memory file backed by a slice of bytes
// with support for seek
type InMemoryFile struct {
	data []byte
	pos  int
}

func NewInMemoryFile(data []byte) *InMemoryFile {
	return &InMemoryFile{
		data: data,
	}
}

// Read reads the next len(p) bytes from the file or until the end is reached.
func (f *InMemoryFile) Read(p []byte) (int, error) {
	if f.pos >= len(f.data) {
		return 0, io.EOF
	}

	n := copy(p, f.data[f.pos:])
	f.pos += n
	return n, nil
}

// ReadAt reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the read pointer is not modified.
func (f *InMemoryFile) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(f.data)) {
		return 0, io.EOF
	}

	if off < 0 {
		return 0, errors.New("negative offset")
	}

	n := copy(p, f.data[off:])
	return n, nil
}

// Write appends the contents of p to the file from the current position.
// It moves the read pointer by the length of the written data.
// The return value n is the length of p; err is always nil.
func (f *InMemoryFile) Write(p []byte) (int, error) {
	if f.pos > len(f.data) {
		f.pos = len(f.data)
	}
	if len(p) <= len(f.data[f.pos:]) {
		n := copy(f.data[f.pos:], p)
		f.pos += n
		return n, nil
	}
	f.data = append(f.data[:f.pos], p...)
	f.pos = len(f.data)
	return len(p), nil
}

// Seek sets the read pointer to pos.
func (f *InMemoryFile) Seek(pos int) {
	if pos < 0 {
		pos = 0
	}
	f.pos = pos
}

// Rewind resets the read pointer to 0.
func (f *InMemoryFile) Rewind() {
	f.Seek(0)
}

// Close clears all the data out of the file and sets the read position to 0.
func (f *InMemoryFile) Close() error {
	f.data = nil
	f.pos = 0
	return nil
}

// Len returns the length of data remaining to be read.
func (f *InMemoryFile) Len() int {
	return len(f.data[f.pos:])
}

// Bytes returns the underlying bytes from the current position.
func (b *InMemoryFile) Bytes() []byte {
	return b.data[b.pos:]
}

func (f *InMemoryFile) Stat() (FileInfo, error) {
	return NewFileInfo(GenericFileInfoOpts{
		Name: "",
		Size: int64(len(f.data)),
	}), nil
}
