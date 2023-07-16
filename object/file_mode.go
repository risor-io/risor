package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/risor-io/risor/op"
	ros "github.com/risor-io/risor/os"
)

type FileMode struct {
	value ros.FileMode
}

func (m *FileMode) Inspect() string {
	return fmt.Sprintf("file_mode(%s)", m.value)
}

func (m *FileMode) Type() Type {
	return FILE_MODE
}

func (m *FileMode) Value() ros.FileMode {
	return m.value
}

func (m *FileMode) Interface() interface{} {
	return m.value
}

func (m *FileMode) String() string {
	return m.value.String()
}

func (m *FileMode) Compare(other Object) (int, error) {
	switch other := other.(type) {
	case *FileMode:
		if m.value < other.value {
			return -1, nil
		} else if m.value > other.value {
			return 1, nil
		}
		return 0, nil
	case *Int:
		if m.value < ros.FileMode(other.value) {
			return -1, nil
		} else if m.value > ros.FileMode(other.value) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("type error: unable to compare file_mode to %s", other.Type())
	}
}

func (m *FileMode) Equals(other Object) Object {
	switch other := other.(type) {
	case *FileMode:
		if m.value == other.value {
			return True
		}
	case *Int:
		if m.value == ros.FileMode(other.value) {
			return True
		}
	}
	return False
}

func (m *FileMode) GetAttr(name string) (Object, bool) {
	switch name {
	case "is_dir":
		return NewBool(m.value.IsDir()), true
	case "is_regular":
		return NewBool(m.value.IsRegular()), true
	case "perm":
		return NewString(m.value.String()), true
	case "type":
		return NewString(fileModeTypeString(m.value)), true
	}
	return nil, false
}

func (m *FileMode) SetAttr(name string, value Object) error {
	return errors.New("type error: unable to set attributes on file_mode objects")
}

func (m *FileMode) IsTruthy() bool {
	return m.value != ros.FileMode(0)
}

func (m *FileMode) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for file_mode: %v", opType))
}

func (m *FileMode) Cost() int {
	return 0
}

func (m *FileMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IsDir     bool   `json:"is_dir"`
		IsRegular bool   `json:"is_regular"`
		Perm      string `json:"perm"`
		Type      string `json:"type"`
	}{
		IsDir:     m.value.IsDir(),
		IsRegular: m.value.IsRegular(),
		Perm:      m.value.String(),
		Type:      fileModeTypeString(m.value),
	})
}

func NewFileMode(value ros.FileMode) *FileMode {
	return &FileMode{value: value}
}

func fileModeTypeString(m ros.FileMode) string {
	switch {
	case m.IsDir():
		return "dir"
	case m.IsRegular():
		return "regular"
	case m&fs.ModeSymlink != 0:
		return "symlink"
	case m&fs.ModeNamedPipe != 0:
		return "named_pipe"
	case m&fs.ModeSocket != 0:
		return "socket"
	case m&fs.ModeDevice != 0:
		return "device"
	case m&fs.ModeCharDevice != 0:
		return "char_device"
	case m&fs.ModeIrregular != 0:
		return "irregular"
	default:
		return "unknown"
	}
}
