package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

const (
	MaxArgs       = 255
	MaxFrameDepth = 1024
	MaxStackDepth = 1024
)

type VM struct {
	ip           int // instruction pointer
	sp           int // stack pointer
	fp           int // frame pointer
	stack        [MaxStackDepth]object.Object
	frames       [MaxFrameDepth]Frame
	tmp          [MaxArgs]object.Object
	currentFrame *Frame
	main         *object.Code
	currentScope *object.Code
	globals      []object.Object
	builtins     []object.Object
}

func NewWithOffset(main *object.Code, ofs int) *VM {
	v := New(main)
	v.ip = ofs - 1
	return v
}

func New(main *object.Code) *VM {
	vm := &VM{
		ip:           -1,
		sp:           -1,
		fp:           -1,
		main:         main,
		currentScope: main,
	}
	if main.Symbols != nil {
		vm.globals = main.Symbols.Variables()
		vm.builtins = main.Symbols.Builtins()
	}
	return vm
}

func (vm *VM) Run(ctx context.Context) error {
	fn := object.NewFunction(object.FunctionOpts{
		Code: vm.currentScope,
	})
	_, err := vm.Eval(ctx, fn, nil)
	return err
}

func (vm *VM) Eval(ctx context.Context, fn *object.Function, args []object.Object) (result object.Object, err error) {

	// Translate any panic into an error so the caller has a good guarantee
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = fmt.Errorf("panic: %v", r)
	// 	}
	// }()

	// Initialize the call frame with the main function
	vm.fp++
	vm.ip++
	vm.currentFrame = &vm.frames[vm.fp]
	vm.currentFrame.ActivateCode(vm.main)

	// Run the program until finished
	for vm.ip < len(vm.currentScope.Instructions) {

		// The current instruction opcode
		opcode := vm.currentScope.Instructions[vm.ip]
		vm.ip++

		switch opcode {
		case op.Nop:
		case op.LoadAttr:
			obj := vm.Pop()
			name := vm.currentScope.Names[vm.fetch()]
			value, found := obj.GetAttr(name)
			if !found {
				return nil, fmt.Errorf("attribute %q not found", name)
			}
			vm.Push(value)
		case op.LoadConst:
			vm.Push(vm.currentScope.Constants[vm.fetch()])
		case op.LoadFast:
			vm.Push(vm.currentFrame.Locals()[vm.fetch()])
		case op.LoadGlobal:
			vm.Push(vm.globals[vm.fetch()])
		case op.LoadFree:
			freeVars := vm.currentFrame.fn.FreeVars()
			vm.Push(freeVars[vm.fetch()].Value())
		case op.LoadBuiltin:
			vm.Push(vm.builtins[vm.fetch()])
		case op.StoreFast:
			vm.currentFrame.Locals()[vm.fetch()] = vm.Pop()
		case op.StoreGlobal:
			vm.globals[vm.fetch()] = vm.Pop()
		case op.StoreFree:
			freeVars := vm.currentFrame.fn.FreeVars()
			freeVars[vm.fetch()].Set(vm.Pop())
		case op.LoadClosure:
			constIndex := vm.fetch()
			freeCount := vm.fetch()
			free := make([]*object.Cell, freeCount)
			for i := uint16(0); i < freeCount; i++ {
				obj := vm.Pop()
				switch obj := obj.(type) {
				case *object.Cell:
					free[i] = obj
				default:
					return nil, errors.New("expected cell")
				}
			}
			fn := vm.currentScope.Constants[constIndex].(*object.Function)
			closure := object.NewClosure(fn, fn.Code(), free)
			vm.Push(closure)
		case op.MakeCell:
			symbolIndex := vm.fetch()
			framesBack := int(vm.fetch())
			frameIndex := vm.fp - framesBack
			if frameIndex < 0 {
				return nil, fmt.Errorf("no frame at depth %d", framesBack)
			}
			frame := &vm.frames[frameIndex]
			locals := frame.CaptureLocals()
			vm.Push(object.NewCell(&locals[symbolIndex]))
		case op.Nil:
			vm.Push(object.Nil)
		case op.True:
			vm.Push(object.True)
		case op.False:
			vm.Push(object.False)
		case op.CompareOp:
			opType := op.CompareOpType(vm.fetch())
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(compare(opType, a, b))
		case op.BinaryOp:
			opType := op.BinaryOpType(vm.fetch())
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(binaryOp(opType, a, b))
		case op.Call:
			argc := int(vm.fetch())
			for argIndex := argc - 1; argIndex >= 0; argIndex-- {
				vm.tmp[argIndex] = vm.Pop()
			}
			obj := vm.Pop()
			switch fn := obj.(type) {
			case *object.Builtin:
				result := fn.Call(ctx, vm.tmp[:argc]...)
				vm.Push(result)
			case *object.Function:
				vm.fp++
				frame := &vm.frames[vm.fp]
				paramsCount := len(fn.Parameters())
				if err := checkCallArgs(fn, argc); err != nil {
					return nil, err
				}
				if argc < paramsCount {
					defaults := fn.Defaults()
					for i := argc; i < len(defaults); i++ {
						vm.tmp[i] = defaults[i]
					}
				}
				code := fn.Code()
				if code.IsNamed {
					vm.tmp[paramsCount] = fn
					argc++
				}
				frame.ActivateFunction(fn, vm.ip, vm.tmp[:argc])
				vm.currentFrame = frame
				vm.currentScope = code
				vm.ip = 0
			default:
				return nil, fmt.Errorf("object is not callable: %T", obj)
			}
		case op.ReturnValue:
			returnAddr := vm.frames[vm.fp].returnAddr
			vm.fp--
			vm.currentFrame = &vm.frames[vm.fp]
			vm.currentScope = vm.currentFrame.Code()
			vm.ip = returnAddr
		case op.PopJumpForwardIfTrue:
			tos := vm.Pop()
			delta := int(vm.fetch()) - 2
			if tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpForwardIfFalse:
			tos := vm.Pop()
			delta := int(vm.fetch()) - 2
			if !tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpBackwardIfTrue:
			tos := vm.Pop()
			delta := int(vm.fetch()) - 2
			if tos.IsTruthy() {
				vm.ip -= delta
			}
		case op.PopJumpBackwardIfFalse:
			tos := vm.Pop()
			delta := int(vm.fetch()) - 2
			if !tos.IsTruthy() {
				vm.ip -= delta
			}
		case op.JumpForward:
			base := vm.ip - 1
			delta := int(vm.fetch())
			vm.ip = base + delta
		case op.JumpBackward:
			base := vm.ip - 1
			delta := int(vm.fetch())
			vm.ip = base - delta
		case op.BuildList:
			count := vm.fetch()
			items := make([]object.Object, count)
			for i := uint16(0); i < count; i++ {
				items[count-1-i] = vm.Pop()
			}
			vm.Push(object.NewList(items))
		case op.BuildMap:
			count := vm.fetch()
			items := make(map[string]object.Object, count)
			for i := uint16(0); i < count; i++ {
				v := vm.Pop()
				k := vm.Pop()
				items[k.(*object.String).Value()] = v
			}
			vm.Push(object.NewMap(items))
		case op.BuildSet:
			count := vm.fetch()
			items := make([]object.Object, count)
			for i := uint16(0); i < count; i++ {
				items[i] = vm.Pop()
			}
			vm.Push(object.NewSet(items))
		case op.BinarySubscr:
			index := vm.Pop()
			obj := vm.Pop()
			container, ok := obj.(object.Container)
			if !ok {
				return nil, fmt.Errorf("object is not a container: %T", obj)
			}
			result, err := container.GetItem(index)
			if err != nil {
				return nil, err.Value()
			}
			vm.Push(result)
		case op.UnaryNegative:
			obj := vm.Pop()
			switch obj := obj.(type) {
			case *object.Int:
				vm.Push(object.NewInt(-obj.Value()))
			case *object.Float:
				vm.Push(object.NewFloat(-obj.Value()))
			default:
				return nil, fmt.Errorf("object is not a number: %T", obj)
			}
		case op.UnaryNot:
			obj := vm.Pop()
			if obj.IsTruthy() {
				vm.Push(object.False)
			} else {
				vm.Push(object.True)
			}
		case op.ContainsOp:
			obj := vm.Pop()
			containerObj := vm.Pop()
			invert := vm.fetch() == 1
			if container, ok := containerObj.(object.Container); ok {
				value := container.Contains(obj)
				if invert {
					value = object.Not(value)
				}
				vm.Push(value)
			} else {
				return nil, fmt.Errorf("object is not a container: %T", container)
			}
		case op.Swap:
			vm.Swap(int(vm.fetch()))
		case op.Halt:
			return nil, nil
		default:
			return nil, fmt.Errorf("unknown opcode: %d", opcode)
		}
	}
	return nil, nil
}

func (vm *VM) TOS() (object.Object, bool) {
	if vm.sp >= 0 {
		return vm.stack[vm.sp], true
	}
	return nil, false
}

func (vm *VM) Pop() object.Object {
	obj := vm.stack[vm.sp]
	vm.sp--
	return obj
}

func (vm *VM) Push(obj object.Object) {
	vm.sp++
	vm.stack[vm.sp] = obj
}

func (vm *VM) Swap(pos int) {
	otherIndex := vm.sp - pos
	tos := vm.stack[vm.sp]
	other := vm.stack[otherIndex]
	vm.stack[otherIndex] = tos
	vm.stack[vm.sp] = other
}

func (vm *VM) fetch() uint16 {
	ip := vm.ip
	vm.ip++
	return uint16(vm.currentScope.Instructions[ip])
}

func checkCallArgs(fn *object.Function, argc int) error {

	// Number of parameters in the function signature
	paramsCount := len(fn.Parameters())

	// Number of required args when the function is called (those without defaults)
	requiredArgsCount := fn.RequiredArgsCount()

	// Check if too many or too few arguments were passed
	if argc > paramsCount || argc < requiredArgsCount {
		switch paramsCount {
		case 0:
			return fmt.Errorf("type error: function takes no arguments (%d given)", argc)
		case 1:
			return fmt.Errorf("type error: function takes 1 argument (%d given)", argc)
		default:
			return fmt.Errorf("type error: function takes %d arguments (%d given)", paramsCount, argc)
		}
	}
	return nil
}
