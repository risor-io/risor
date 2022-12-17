package object

import (
	"fmt"
)

type Module struct {
	Name  string
	Scope Scope
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	return m.String()
}

func (m *Module) GetAttr(name string) (Object, bool) {
	return m.Scope.Get(name)
}

func (m *Module) ToInterface() interface{} {
	return nil
}

func (m *Module) String() string {
	return fmt.Sprintf("module(%s)", m.Name)
}

func (m *Module) Compare(other Object) (int, error) {
	typeComp := CompareTypes(m, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Module)
	if m.Name == otherStr.Name {
		return 0, nil
	}
	if m.Name > otherStr.Name {
		return 1, nil
	}
	return -1, nil
}

func (m *Module) Equals(other Object) Object {
	return NewBool(other.Type() == MODULE && m.Name == other.(*Module).Name)
}
