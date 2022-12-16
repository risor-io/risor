package object

import (
	"context"
	"fmt"
)

// BuiltinFunction holds the type of a built-in function.
type BuiltinFunction func(ctx context.Context, args ...Object) Object

// Builtin wraps func and implements Object interface.
type Builtin struct {
	// The function that this object wraps.
	Fn BuiltinFunction

	// The name of the function.
	Name string

	// The module the function originates from (optional)
	Module *Module

	// The name of the module this function origiantes from.
	// This is only used for overriding builtins.
	ModuleName string
}

// Type returns the type of this object.
func (b *Builtin) Type() Type {
	return BUILTIN
}

// Inspect returns a string-representation of the given object.
func (b *Builtin) Inspect() string {
	if b.Module == nil {
		return fmt.Sprintf("builtin(%s)", b.Name)
	}
	return fmt.Sprintf("builtin(%s.%s)", b.Module.Name, b.Name)
}

func (b *Builtin) String() string {
	return b.Inspect()
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (b *Builtin) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", b.Type(), method)
}

func (b *Builtin) GetAttr(name string) (Object, bool) {
	return nil, false
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Builtin) ToInterface() interface{} {
	return b.Fn
}

// Returns a string that uniquely identifies this builtin function.
func (b *Builtin) Key() string {
	if b.Module == nil && b.ModuleName == "" {
		return b.Name
	} else if b.ModuleName != "" {
		return fmt.Sprintf("%s.%s", b.ModuleName, b.Name)
	}
	return fmt.Sprintf("%s.%s", b.Module.Name, b.Name)
}

func (b *Builtin) Equals(other Object) Object {
	if other.Type() != BUILTIN {
		return NewBool(false)
	}
	value := fmt.Sprintf("%v", b.Fn) == fmt.Sprintf("%v", other.(*Builtin).Fn)
	return NewBool(value)
}

// NewNoopBuiltin creates a builtin function that has no effect.
func NewNoopBuiltin(Name string, Module *Module) *Builtin {
	b := &Builtin{
		Fn: func(ctx context.Context, args ...Object) Object {
			return Nil
		},
		Name:   Name,
		Module: Module,
	}
	if Module != nil {
		b.ModuleName = Module.Name
	}
	return b
}
