package object

import (
	"context"
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
	callable     BuiltinFunction
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

func (m *Module) SetAttr(name string, value Object) error {
	if name == "__name__" {
		return fmt.Errorf("attribute error: cannot set attribute %q", name)
	}
	if _, found := m.builtins[name]; found {
		return fmt.Errorf("attribute error: cannot set attribute %q", name)
	}
	if index, found := m.globalsIndex[name]; found {
		m.globals[index] = value
		return nil
	}
	return fmt.Errorf("attribute error: module has no attribute %q", name)
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

func (m *Module) UseGlobals(globals []Object) {
	if len(globals) != len(m.globals) {
		panic(fmt.Sprintf("invalid module globals length: %d, expected: %d",
			len(globals), len(m.globals)))
	}
	m.globals = globals
}

func (m *Module) Call(ctx context.Context, args ...Object) Object {
	if m.callable == nil {
		return NewError(fmt.Errorf("exec error: module %q is not callable", m.name))
	}
	return m.callable(ctx, args...)
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
		builtins:     map[string]Object{},
		code:         code,
		globals:      globals,
		globalsIndex: globalsIndex,
	}
}

func NewBuiltinsModule(name string, contents map[string]Object, callableOption ...BuiltinFunction) *Module {
	builtins := map[string]Object{}
	for k, v := range contents {
		builtins[k] = v
	}
	var callable BuiltinFunction
	if len(callableOption) > 0 {
		callable = callableOption[0]
	}
	return &Module{
		name:         name,
		builtins:     builtins,
		callable:     callable,
		globalsIndex: map[string]int{},
	}
}
