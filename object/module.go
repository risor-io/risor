package object

import (
	"errors"
	"fmt"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/op"
)

type Module struct {
	*base
	name         string
	code         *compiler.Code
	builtins     map[string]Object
	globals      []Object
	globalsIndex map[string]int
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	return m.String()
}

func (m *Module) GetAttr(name string) (Object, bool) {
	switch name {
	case "__name__":
		return NewString(m.name), true
	}
	if builtin, found := m.builtins[name]; found {
		return builtin, true
	}
	if index, found := m.globalsIndex[name]; found {
		return m.globals[index], true
	}
	return nil, false
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

func (m *Module) Code() *compiler.Code {
	return m.code
}

func (m *Module) Compare(other Object) (int, error) {
	typeComp := CompareTypes(m, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherMod := other.(*Module)
	if m.name == otherMod.name {
		return 0, nil
	}
	if m.name > otherMod.name {
		return 1, nil
	}
	return -1, nil
}

func (m *Module) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for module: %v", opType))
}

func (m *Module) Equals(other Object) Object {
	if m == other {
		return True
	}
	return False
}

func (m *Module) MarshalJSON() ([]byte, error) {
	return nil, errors.New("type error: unable to marshal module")
}

func NewModule(name string, code *compiler.Code) *Module {
	globalsIndex := map[string]int{}
	globalsCount := code.GlobalsCount()
	globals := make([]Object, globalsCount)
	for i := 0; i < globalsCount; i++ {
		symbol := code.Global(i)
		globalsIndex[symbol.Name()] = int(i)
		value := symbol.Value()
		switch value := value.(type) {
		case int64:
			globals[i] = NewInt(value)
		case float64:
			globals[i] = NewFloat(value)
		case string:
			globals[i] = NewString(value)
		case bool:
			globals[i] = NewBool(value)
		case nil:
			globals[i] = Nil
		// TODO: functions, others?
		default:
			panic(fmt.Sprintf("unsupported global type: %T", value))
		}
	}
	return &Module{
		name:         name,
		code:         code,
		globals:      globals,
		globalsIndex: globalsIndex,
	}
}

func NewBuiltinsModule(name string, contents map[string]Object) *Module {
	builtins := map[string]Object{}
	for k, v := range contents {
		builtins[k] = v
	}
	return &Module{
		name:     name,
		builtins: builtins,
	}
}
