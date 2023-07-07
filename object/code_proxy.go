package object

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/op"
)

type CodeProxy struct {
	*base
	name     string
	builtins *List
	code     *Code
}

func (c *CodeProxy) Inspect() string {
	return c.String()
}

func (c *CodeProxy) String() string {
	return fmt.Sprintf("code(%s)", c.name)
}

func (c *CodeProxy) Type() Type {
	return CODE
}

func (c *CodeProxy) Interface() interface{} {
	return c.code
}

func (c *CodeProxy) Equals(other Object) Object {
	if c == other {
		return True
	}
	return False
}

func (c *CodeProxy) GetAttr(name string) (Object, bool) {
	switch name {
	case "builtins":
		return &Builtin{
			name: "code.builtins",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 0 {
					return NewArgsError("code.builtins", 0, len(args))
				}
				return c.builtins
			},
		}, true
	default:
		return nil, false
	}
}

func (c *CodeProxy) IsTruthy() bool {
	return len(c.code.Instructions) > 0
}

func (c *CodeProxy) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for code: %v", opType))
}

func (c *CodeProxy) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal code")
}

func NewCodeProxy(c *Code) *CodeProxy {

	// Provide some isolation
	builtins := c.Symbols.Root().Builtins()
	copiedBuiltins := make([]Object, len(builtins))
	copy(copiedBuiltins, builtins)

	return &CodeProxy{
		name:     c.Name,
		builtins: NewList(copiedBuiltins),
		code:     c,
	}
}
