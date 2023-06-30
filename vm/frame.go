package vm

import (
	"github.com/risor-io/risor/object"
)

const DefaultFrameLocals = 8

type Frame struct {
	returnAddr     int
	localsCount    uint16
	fn             *object.Function
	code           *object.Code
	storage        [DefaultFrameLocals]object.Object
	locals         []object.Object
	extendedLocals []object.Object
	capturedLocals []object.Object
}

func (f *Frame) ActivateCode(code *object.Code) {
	f.code = code
	f.fn = nil
	f.returnAddr = 0
	f.localsCount = code.Symbols.Size()
	f.capturedLocals = nil
	for i := 0; i < DefaultFrameLocals; i++ {
		f.storage[i] = nil
	}
	// Decide where to store local variables. If the frame storage has enough
	// space, use that. Otherwise, allocate a new slice as extendedLocals.
	// After this, f.locals will always point to the correct storage.
	if f.localsCount > DefaultFrameLocals {
		f.extendedLocals = make([]object.Object, f.localsCount)
		f.locals = f.extendedLocals
	} else {
		f.extendedLocals = nil
		f.locals = f.storage[:f.localsCount]
	}
}

func (f *Frame) ActivateFunction(fn *object.Function, returnAddr int, localValues []object.Object) {
	// Activate the function's code
	f.ActivateCode(fn.Code())
	f.fn = fn
	// Save the instruction pointer of the caller
	f.returnAddr = returnAddr
	// Initialize any local variables that were provided.
	// Note the copy builtin is slower than this loop.
	for i := 0; i < len(localValues); i++ {
		f.locals[i] = localValues[i]
	}
}

func (f *Frame) SetReturnAddr(addr int) {
	f.returnAddr = addr
}

func (f *Frame) Locals() []object.Object {
	return f.locals
}

func (f *Frame) CaptureLocals() []object.Object {
	if f.capturedLocals != nil {
		return f.capturedLocals
	}
	if f.extendedLocals != nil {
		f.capturedLocals = f.extendedLocals
		return f.capturedLocals
	}
	newStorage := make([]object.Object, len(f.locals))
	copy(newStorage, f.locals)
	f.capturedLocals = newStorage
	f.locals = newStorage
	return newStorage
}

func (f *Frame) Function() *object.Function {
	return f.fn
}

func (f *Frame) Code() *object.Code {
	return f.code
}
