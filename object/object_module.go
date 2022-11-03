package object

import (
	"bytes"
	"fmt"
)

type Module struct {
	Name  string
	Scope interface{}
}

func (m *Module) Type() Type {
	return MODULE_OBJ
}

func (m *Module) Inspect() string {
	var out bytes.Buffer
	out.WriteString("/")
	out.WriteString(m.Name)
	out.WriteString("/")
	return out.String()
}

func (m *Module) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (m *Module) ToInterface() interface{} {
	return "<MODULE>"
}

func (m *Module) String() string {
	return fmt.Sprintf("Module(%s)", m.Name)
}
