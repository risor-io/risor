package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/risor-io/risor/op"
	ros "github.com/risor-io/risor/os"
)

type FileInfo struct {
	value ros.FileInfo
	mode  *FileMode
}

func (f *FileInfo) Inspect() string {
	return f.String()
}

func (f *FileInfo) Type() Type {
	return FILE_INFO
}

func (f *FileInfo) Value() ros.FileInfo {
	return f.value
}

func (f *FileInfo) Interface() interface{} {
	return f.value
}

func (f *FileInfo) String() string {
	v := f.value
	return fmt.Sprintf("file_info(name=%s, mode=%s, size=%d, mod_time=%v)",
		v.Name(), v.Mode().String(), v.Size(), v.ModTime().Format(time.RFC3339))
}

func (f *FileInfo) Compare(other Object) (int, error) {
	return 0, errors.New("type error: unable to compare file_info objects")
}

func (f *FileInfo) Equals(other Object) Object {
	if f == other {
		return True
	}
	return False
}

func (f *FileInfo) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(f.value.Name()), true
	case "size":
		return NewInt(f.value.Size()), true
	case "mod_time":
		return NewTime(f.value.ModTime()), true
	case "mode":
		return f.mode, true
	case "is_dir":
		return NewBool(f.value.IsDir()), true
	}
	return nil, false
}

func (f *FileInfo) SetAttr(name string, value Object) error {
	return errors.New("type error: unable to set attributes on file_info objects")
}

func (f *FileInfo) IsTruthy() bool {
	return true
}

func (f *FileInfo) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for file_info: %v", opType))
}

func (f *FileInfo) Cost() int {
	return 0
}

func (f *FileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IsDir   bool      `json:"is_dir"`
		Mode    *FileMode `json:"mode"`
		ModTime time.Time `json:"mod_time"`
		Name    string    `json:"name"`
		Size    int64     `json:"size"`
	}{
		IsDir:   f.value.IsDir(),
		Mode:    f.mode,
		ModTime: f.value.ModTime(),
		Name:    f.value.Name(),
		Size:    f.value.Size(),
	})
}

func NewFileInfo(value ros.FileInfo) *FileInfo {
	return &FileInfo{value: value, mode: NewFileMode(value.Mode())}
}
