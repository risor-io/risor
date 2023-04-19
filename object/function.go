package object

import (
	"github.com/cloudcmds/tamarin/internal/op"
)

// Function is a function that has been compiled to bytecode.
type Function struct {
	*DefaultImpl
	name          string
	parameters    []string
	defaults      []Object
	defaultsCount int
	scope         *Code
	freeVars      []*Cell
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Inspect() string {
	return "function()"
}

func (f *Function) Instructions() []op.Code {
	return f.scope.Instructions
}

func (f *Function) FreeVars() []*Cell {
	return f.freeVars
}

func (f *Function) Code() *Code {
	return f.scope
}

func (f *Function) Parameters() []string {
	return f.parameters
}

func (f *Function) Defaults() []Object {
	return f.defaults
}

func (f *Function) RequiredArgsCount() int {
	return len(f.parameters) - f.defaultsCount
}

type FunctionOpts struct {
	Name           string
	ParameterNames []string
	Defaults       []Object
	Code           *Code
}

func NewFunction(opts FunctionOpts) *Function {
	var defaultsCount int
	for _, value := range opts.Defaults {
		if value != nil {
			defaultsCount++
		}
	}
	return &Function{
		name:          opts.Name,
		parameters:    opts.ParameterNames,
		defaults:      opts.Defaults,
		defaultsCount: defaultsCount,
		scope:         opts.Code,
	}
}

func NewClosure(
	fn *Function,
	scope *Code,
	freeVars []*Cell,
) *Function {
	return &Function{
		name:          fn.name,
		parameters:    fn.parameters,
		defaults:      fn.defaults,
		defaultsCount: fn.defaultsCount,
		scope:         scope,
		freeVars:      freeVars,
	}
}
