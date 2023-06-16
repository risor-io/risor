package object

import (
	"fmt"

	"github.com/cloudcmds/tamarin/v2/op"
)

type Module struct {
	name string
	code *Code
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
	case "__builtins__":
		builtins := m.code.Builtins()
		copied := make([]Object, len(builtins))
		copy(copied, builtins)
		return NewList(copied), true
	}
	resolution, found := m.code.Symbols.Lookup(name)
	if !found {
		return nil, false
	}
	switch resolution.Scope {
	case ScopeBuiltin:
		return m.code.Builtins()[resolution.Symbol.Index], true
	case ScopeGlobal:
		return m.code.Globals()[resolution.Symbol.Index], true
	default:
		panic("module attribute resolution scope not builtin or global")
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

func NewModule(name string, code *Code) *Module {
	return &Module{
		name: name,
		code: code,
	}
}

func NewBuiltinsModule(name string, contents map[string]Object) *Module {
	code := NewCode(name)
	for name, obj := range contents {
		code.Symbols.InsertBuiltin(name, obj)
	}
	return &Module{
		name: name,
		code: code,
	}
}
