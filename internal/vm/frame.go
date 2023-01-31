package vm

import "github.com/cloudcmds/tamarin/object"

type Frame struct {
	fn         *object.Function
	locals     []object.Object
	returnAddr int
	// ip          int
	// basePointer int
}

func NewFrame(fn *object.Function, locals []object.Object, returnAddr int) *Frame {
	return &Frame{
		fn:         fn,
		locals:     locals,
		returnAddr: returnAddr,
	}
}
