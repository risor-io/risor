package vm

import (
	"fmt"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type code struct {
	*compiler.Code
	Instructions []op.Code
	Constants    []object.Object
	Globals      []object.Object
	Names        []string
}

func wrapCode(cc *compiler.Code) *code {
	// Note that this does NOT set the Globals field.
	c := &code{
		Code:         cc,
		Instructions: make([]op.Code, cc.InstructionCount()),
		Constants:    make([]object.Object, cc.ConstantsCount()),
		Names:        make([]string, cc.NameCount()),
	}
	for i := 0; i < cc.InstructionCount(); i++ {
		c.Instructions[i] = cc.Instruction(i)
	}
	for i := 0; i < cc.NameCount(); i++ {
		c.Names[i] = cc.Name(i)
	}
	for i := 0; i < cc.ConstantsCount(); i++ {
		constant := cc.Constant(i)
		switch constant := constant.(type) {
		case int:
			c.Constants[i] = object.NewInt(int64(constant))
		case int64:
			c.Constants[i] = object.NewInt(constant)
		case float64:
			c.Constants[i] = object.NewFloat(constant)
		case string:
			c.Constants[i] = object.NewString(constant)
		case bool:
			c.Constants[i] = object.NewBool(constant)
		case *compiler.Function:
			c.Constants[i] = object.NewFunction(constant)
		case nil:
			c.Constants[i] = object.Nil
		default:
			panic(fmt.Sprintf("unsupported constant type: %T", constant))
		}
	}
	return c
}

func (c *code) InstructionCount() int {
	return len(c.Instructions)
}

func (c *code) ConstantsCount() int {
	return len(c.Constants)
}

func (c *code) GlobalsCount() int {
	return len(c.Globals)
}

func (c *code) Clone() *code {
	clone := &code{
		Code:         c.Code,
		Instructions: make([]op.Code, len(c.Instructions)),
		Constants:    make([]object.Object, len(c.Constants)),
		Globals:      make([]object.Object, len(c.Globals)),
		Names:        make([]string, len(c.Names)),
	}
	copy(clone.Instructions, c.Instructions)
	copy(clone.Constants, c.Constants)
	copy(clone.Globals, c.Globals)
	copy(clone.Names, c.Names)
	return clone
}

func loadChildCode(root *code, cc *compiler.Code) *code {
	c := wrapCode(cc)
	c.Globals = root.Globals
	return c
}

func loadRootCode(cc *compiler.Code, globals map[string]object.Object) *code {
	c := wrapCode(cc)
	globalNames := cc.GlobalNames()
	c.Globals = make([]object.Object, len(globalNames))
	for i, name := range globalNames {
		if value, found := globals[name]; found {
			c.Globals[i] = value
		}
	}
	return c
}
