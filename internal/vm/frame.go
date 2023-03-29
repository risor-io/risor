package vm

import (
	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/object"
)

const DefaultFrameLocals = 4

type Frame struct {
	fn             *object.CompiledFunction
	scope          *compiler.Scope
	returnAddr     int
	localsCount    int
	locals         [DefaultFrameLocals]object.Object
	extendedLocals []object.Object
}

func (f *Frame) Init(fn *object.CompiledFunction, returnAddr int, localsCount int) {
	f.fn = fn
	if fn != nil {
		f.scope = fn.Scope().(*compiler.Scope)
	} else {
		f.scope = nil
	}
	f.returnAddr = returnAddr
	f.localsCount = localsCount
	if localsCount > DefaultFrameLocals {
		f.extendedLocals = make([]object.Object, localsCount)
	}
}

func (f *Frame) InitWithLocals(fn *object.CompiledFunction, returnAddr int, locals []object.Object) {
	count := len(locals)
	f.Init(fn, returnAddr, count)
	if count > DefaultFrameLocals {
		copy(f.extendedLocals, locals)
	} else {
		copy(f.locals[:], locals)
	}
}

func (f *Frame) Locals() []object.Object {
	if f.localsCount > DefaultFrameLocals {
		return f.extendedLocals
	}
	return f.locals[:f.localsCount]
}

func (f *Frame) Function() *object.CompiledFunction {
	return f.fn
}

func (f *Frame) Scope() *compiler.Scope {
	return f.scope
}

// func (f Frame) SetLocals(locals []object.Object) {
// 	f.localsCount = len(locals)
// 	if len(locals) <= DefaultFrameLocals {
// 		copy(f.locals[:], locals)
// 	} else {
// 		f.extendedLocals = locals
// 	}
// }

// func (f *Frame) NewChild(
// 	fn *object.CompiledFunction,
// 	locals []object.Object,
// 	returnAddr int,
// ) *Frame {
// 	return &Frame{
// 		fn: fn,
// 		// locals:     locals,
// 		returnAddr: returnAddr,
// 		// scope:      fn.Scope().(*compiler.Scope),
// 		// parent:     f,
// 	}
// }

// func NewFrame(
// 	fn *object.CompiledFunction,
// 	locals []object.Object,
// 	returnAddr int,
// 	scope *compiler.Scope,
// ) *Frame {
// 	return &Frame{
// 		fn:         fn,
// 		locals:     locals,
// 		returnAddr: returnAddr,
// 		scope:      scope,
// 	}
// }
