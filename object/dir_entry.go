package object

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/risor-io/risor/op"
	ros "github.com/risor-io/risor/os"
)

type DirEntry struct {
	value    ros.DirEntry
	fileInfo *FileInfo
}

func (d *DirEntry) Inspect() string {
	return fmt.Sprintf("dir_entry(name=%s, type=%s)", d.value.Name(),
		fileModeTypeString(d.value.Type()))
}

func (d *DirEntry) Type() Type {
	return DIR_ENTRY
}

func (d *DirEntry) Value() ros.DirEntry {
	return d.value
}

func (d *DirEntry) Interface() interface{} {
	return d.value
}

func (d *DirEntry) String() string {
	return fmt.Sprintf("dir_entry(%v)", d.value)
}

func (d *DirEntry) Compare(other Object) (int, error) {
	return 0, errors.New("type error: unable to compare dir_entry objects")
}

func (d *DirEntry) Equals(other Object) Object {
	if d == other {
		return True
	}
	return False
}

func (d *DirEntry) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(d.value.Name()), true
	case "type":
		return NewString(fileModeTypeString(d.value.Type())), true
	case "is_dir":
		return NewBool(d.value.IsDir()), true
	case "info":
		return NewBuiltin("dir_entry.info",
			func(ctx context.Context, args ...Object) Object {
				info, err := d.value.Info()
				if err != nil {
					return NewError(err)
				}
				return NewFileInfo(info)
			}), true
	}
	return nil, false
}

func (d *DirEntry) SetAttr(name string, value Object) error {
	return errors.New("type error: unable to set attributes on dir_entry objects")
}

func (d *DirEntry) IsTruthy() bool {
	return true
}

func (d *DirEntry) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for dir_entry: %v", opType))
}

func (d *DirEntry) Cost() int {
	return 0
}

func (d *DirEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		IsDir bool   `json:"is_dir"`
		Info  Object `json:"info"`
	}{
		Name:  d.value.Name(),
		Type:  fileModeTypeString(d.value.Type()),
		IsDir: d.value.IsDir(),
		Info:  d.fileInfo,
	})
}

func (d *DirEntry) FileInfo() (*FileInfo, bool) {
	return d.fileInfo, d.fileInfo != nil
}

func NewDirEntry(value ros.DirEntry, fileInfo ...*FileInfo) *DirEntry {
	var info *FileInfo
	if len(fileInfo) > 0 {
		info = fileInfo[0]
	}
	return &DirEntry{value: value, fileInfo: info}
}
