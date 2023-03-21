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
