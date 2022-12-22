package object

import (
	"fmt"
)

type Module struct {
	name  string
	scope Scope
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	return m.String()
}

func (m *Module) GetAttr(name string) (Object, bool) {
	return m.scope.Get(name)
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

func (m *Module) Equals(other Object) Object {
	return NewBool(other.Type() == MODULE && m.name == other.(*Module).name)
}

func NewModule(name string, scope Scope) *Module {
	return &Module{name, scope}
}
