package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
)

type Module struct {
	name     string
	code     *Code
	globals  []Object // main.Symbols.Variables()
	builtins []Object // main.Symbols.Builtins()
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	return m.String()
}

func (m *Module) GetAttr(name string) (Object, bool) {
	resolution, found := m.code.Symbols.Lookup(name)
	if !found {
		return nil, false
	}
	switch resolution.Scope {
	case ScopeBuiltin:
		return m.builtins[resolution.Symbol.Index], true
	case ScopeGlobal:
		return m.globals[resolution.Symbol.Index], true
	default:
		panic("module attribute resolution scope not builtin or global")
		return nil, false
	}
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

func (m *Module) Code() *Code {
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

// func (m *Module) Register(name string, obj Object) {
// 	m.attrs[name] = obj
// }

func NewModule(name string, code *Code) *Module {
	return &Module{
		name:     name,
		code:     code,
		globals:  code.Symbols.Variables(),
		builtins: code.Symbols.Builtins(),
	}
}

func NewBuiltinsModule(name string, contents map[string]Object) *Module {
	code := NewCode(name)
	for name, obj := range contents {
		code.Symbols.InsertBuiltin(name, obj)
	}
	return &Module{
		name:     name,
		code:     code,
		globals:  code.Symbols.Variables(),
		builtins: code.Symbols.Builtins(),
	}
}
