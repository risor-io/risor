package vm

import (
	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/object"
)

type Frame struct {
	fn         *object.CompiledFunction
	locals     []object.Object
	returnAddr int
	scope      *compiler.Scope
	parent     *Frame
}

func (f *Frame) NewChild(
	fn *object.CompiledFunction,
	locals []object.Object,
	returnAddr int,
) *Frame {
	return &Frame{
		fn:         fn,
		locals:     locals,
		returnAddr: returnAddr,
		scope:      fn.Scope().(*compiler.Scope),
		parent:     f,
	}
}

func NewFrame(
	fn *object.CompiledFunction,
	locals []object.Object,
	returnAddr int,
	scope *compiler.Scope,
) *Frame {
	return &Frame{
		fn:         fn,
		locals:     locals,
		returnAddr: returnAddr,
		scope:      scope,
	}
}
