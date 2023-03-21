package object

import (
	"github.com/cloudcmds/tamarin/internal/op"
)

// CompiledFunction is a function that has been compiled to bytecode.
type CompiledFunction struct {
	*DefaultImpl
	name         string
	parameters   []string
	defaults     []Object
	instructions []op.Code
	scope        interface{}
	freeVars     []*Cell
}

func (f *CompiledFunction) Type() Type {
	return COMPILED_FUNCTION
}

func (f *CompiledFunction) Name() string {
	if f.name == "" {
		return "anonymous"
	}
	return f.name
}

func (f *CompiledFunction) Inspect() string {
	return "compiled_function()"
}

func (f *CompiledFunction) Instructions() []op.Code {
	return f.instructions
}

func (f *CompiledFunction) FreeVars() []*Cell {
	return f.freeVars
}

func (f *CompiledFunction) Scope() interface{} {
	return f.scope
}

func (f *CompiledFunction) Parameters() []string {
	return f.parameters
}

func (f *CompiledFunction) Defaults() []Object {
	return f.defaults
}

func NewCompiledFunction(
	name string,
	parameters []string,
	defaults []Object,
	instructions []op.Code,
	scope interface{},
) *CompiledFunction {
	return &CompiledFunction{
		name:         name,
		parameters:   parameters,
		defaults:     defaults,
		instructions: instructions,
		scope:        scope,
	}
}

func NewClosure(
	fn *CompiledFunction,
	scope interface{},
	freeVars []*Cell,
) *CompiledFunction {
	return &CompiledFunction{
		name:         fn.name,
		parameters:   fn.parameters,
		defaults:     fn.defaults,
		instructions: fn.instructions,
		scope:        scope,
		freeVars:     freeVars,
	}
}
