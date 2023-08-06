package os

import (
	"errors"
	"io"
)

type NilFile struct{}

func (f *NilFile) Close() error {
	return nil
}

func (f *NilFile) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (f *NilFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func (f *NilFile) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("io error: unable to seek in nil file")
}

func (f *NilFile) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (f *NilFile) Stat() (FileInfo, error) {
	return NewFileInfo(GenericFileInfoOpts{
		Name: "",
		Size: 0,
	}), nil
}
