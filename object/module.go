package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

type Module struct {
	name  string
	attrs map[string]Object
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	return m.String()
}

func (m *Module) GetAttr(name string) (Object, bool) {
	obj, found := m.attrs[name]
	return obj, found
}

func (m *Module) Interface() interface{} {
	return nil
}

func (m *Module) String() string {
	return fmt.Sprintf("module(%s)", m.name)
}

func (m *Module) Name() *String {
	return NewString(m.name)
}

func (m *Module) Compare(other Object) (int, error) {
	typeComp := CompareTypes(m, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Module)
	if m.name == otherStr.name {
		return 0, nil
	}
	if m.name > otherStr.name {
		return 1, nil
	}
	return -1, nil
}

func (m *Module) IsTruthy() bool {
	return true
}

func (m *Module) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for module: %v", opType))
}

func (m *Module) Equals(other Object) Object {
	if m == other {
		return True
	}
	return False
}

func (m *Module) Register(name string, obj Object) {
	m.attrs[name] = obj
}

func NewModule(name string) *Module {
	return &Module{name: name, attrs: map[string]Object{}}
}
