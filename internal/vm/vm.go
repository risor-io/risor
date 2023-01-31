package vm

import (
	"errors"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

const MaxFrameCount = 1024

type VM struct {
	ip         int
	sp         int
	constants  []object.Object
	stack      *Stack[object.Object]
	frameStack *Stack[*Frame]
	code       []op.Code
}

func New(constants []object.Object, code []op.Code) *VM {
	return &VM{
		stack:      NewStack[object.Object](1024),
		frameStack: NewStack[*Frame](1024),
		constants:  constants,
		sp:         -1,
		code:       code,
	}
}

func (vm *VM) Run() error {
	for vm.ip < len(vm.code) {
		opcode := vm.code[vm.ip]
		fmt.Println("IP:", vm.ip, "OPCODE:", opcode)
		vm.ip++
		switch opcode {
		case op.Nop:
		case op.LoadConst:
			constIndex := vm.fetch2()
			vm.stack.Push(vm.constants[constIndex])
			fmt.Println("PUSHED CONSTANT:", constIndex, vm.constants[constIndex])
		case op.StoreFast:
			obj := vm.pop()
			idx := vm.fetch()
			frame, ok := vm.frameStack.Top()
			if !ok {
				return errors.New("no frame")
			}
			frame.locals[idx] = obj
		case op.Nil:
			vm.stack.Push(object.Nil)
		case op.True:
			vm.stack.Push(object.True)
		case op.False:
			vm.stack.Push(object.False)
		case op.BinaryOp:
			opType := vm.fetch()
			b := vm.pop()
			a := vm.pop()
			fmt.Println("BINOP A:", a, "B:", b, "OP:", opType)
			switch opType {
			case int(op.Add):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() + b.(*object.Int).Value()))
			case int(op.Subtract):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() - b.(*object.Int).Value()))
			case int(op.Multiply):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() * b.(*object.Int).Value()))
			case int(op.Divide):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() / b.(*object.Int).Value()))
			case int(op.Modulo):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() % b.(*object.Int).Value()))
			case int(op.And):
				vm.stack.Push(object.NewInt(a.(*object.Int).Value() & b.(*object.Int).Value()))
			}
		case op.Call:
			if vm.frameStack.Size() >= MaxFrameCount {
				return errors.New("stack overflow")
			}
			fn, ok := vm.pop().(*object.Function)
			if !ok {
				return fmt.Errorf("not a function: %T", fn)
			}
			fnCode := fn.GetCodeStart()
			argc := vm.fetch()
			args := make([]object.Object, argc)
			for i := 0; i < argc; i++ {
				args[len(args)-1-i] = vm.pop()
			}
			frame := NewFrame(fn, args, vm.ip)
			vm.frameStack.Push(frame)
			vm.ip = fnCode
			fmt.Println("CALL IP", fnCode)
		case op.ReturnValue:
			frame, ok := vm.frameStack.Pop()
			if !ok {
				return errors.New("invalid return")
			}
			if vm.fetch() == 0 {
				// Ensure that a return value is on top of the stack
				vm.stack.Push(object.Nil)
			}
			vm.ip = frame.returnAddr
		case op.JumpForward:
			jump := vm.fetch2()
			vm.ip = jump
			fmt.Println("JUMPED TO", vm.ip)
		case op.Print:
			fmt.Println("PRINT", vm.top())
		case op.LoadFast:
			frame, ok := vm.frameStack.Top()
			if !ok {
				return errors.New("invalid frame")
			}
			index := vm.fetch()
			vm.stack.Push(frame.locals[index])
		case op.Halt:
			return nil
		default:
			return fmt.Errorf("unknown opcode: %d", opcode)
		}
	}
	return nil
}

func (vm *VM) TOS() (object.Object, bool) {
	return vm.stack.Top()
}

func (vm *VM) Frame() (*Frame, bool) {
	return vm.frameStack.Top()
}

func (vm *VM) pop() object.Object {
	obj, ok := vm.stack.Pop()
	if !ok {
		return nil
	}
	return obj
}

func (vm *VM) top() object.Object {
	obj, ok := vm.stack.Top()
	if !ok {
		return nil
	}
	return obj
}

func (vm *VM) fetch() int {
	ip := vm.ip
	vm.ip++
	return int(vm.code[ip])
}

func (vm *VM) fetch2() int {
	v1 := vm.fetch()
	v2 := vm.fetch()
	return v1 | v2<<8
}

func (vm *VM) fetch4() int {
	v1 := vm.fetch()
	v2 := vm.fetch()
	v3 := vm.fetch()
	v4 := vm.fetch()
	return v1 | v2<<8 | v3<<16 | v4<<24
}
