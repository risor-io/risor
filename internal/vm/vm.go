package vm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/importer"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

const (
	MaxArgs       = 255
	MaxFrameDepth = 1024
	MaxStackDepth = 1024
)

type Options struct {
	Main              *object.Code
	InstructionOffset int
	Importer          importer.Importer
}

type VM struct {
	ip          int // instruction pointer
	sp          int // stack pointer
	fp          int // frame pointer
	stack       [MaxStackDepth]object.Object
	frames      [MaxFrameDepth]Frame
	tmp         [MaxArgs]object.Object
	activeFrame *Frame
	activeCode  *object.Code
	main        *object.Code
	globals     []object.Object
	builtins    []object.Object
	importer    importer.Importer
	modules     map[string]*object.Module
}

func New(opts Options) *VM {
	ipOffset := 0
	if opts.InstructionOffset > 0 {
		ipOffset = opts.InstructionOffset
	}
	vm := &VM{
		sp:       -1,
		ip:       ipOffset,
		main:     opts.Main,
		importer: opts.Importer,
		modules:  map[string]*object.Module{},
	}
	if opts.Main != nil && opts.Main.Symbols != nil {
		vm.globals = opts.Main.Symbols.Variables()
		vm.builtins = opts.Main.Symbols.Builtins()
	}
	return vm
}

func (vm *VM) Run(ctx context.Context) (err error) {

	// Translate any panic into an error so the caller has a good guarantee
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Activate the "main" entrypoint code in frame 0 and then run it
	vm.fp = 0
	vm.activeFrame = &vm.frames[vm.fp]
	vm.activeFrame.ActivateCode(vm.main)
	vm.activeCode = vm.main
	ctx = object.WithCallFunc(ctx, vm.callFunction)
	err = vm.eval(ctx)
	return
}

// Evaluate the active code. The caller must initialize the following variables
// before calling this function:
//   - vm.ip - instruction pointer within the active code
//   - vm.fp - frame pointer with the active code already set
//   - vm.activeCode - the code object to execute
//   - vm.activeFrame - the active call frame to use
//
// Assuming this function returns without error, the result of the evaluation
// will be on the top of the stack.
func (vm *VM) eval(ctx context.Context) (err error) {

	// fmt.Println("eval; ip @ ", vm.ip)

	// Run to the end of the active code
	for vm.ip < len(vm.activeCode.Instructions) {

		// The current instruction opcode
		opcode := vm.activeCode.Instructions[vm.ip]

		fmt.Println("ip", vm.ip)

		// Advance the instruction pointer to the next instruction. Note that
		// this is done before we actually execute the current instruction, so
		// relative jump instructions will need to take this into account.
		vm.ip++

		// Dispatch the instruction
		switch opcode {
		case op.Nop:
		case op.LoadAttr:
			obj := vm.Pop()
			name := vm.activeCode.Names[vm.fetch()]
			value, found := obj.GetAttr(name)
			if !found {
				return fmt.Errorf("attribute %q not found", name)
			}
			vm.Push(value)
		case op.LoadConst:
			vm.Push(vm.activeCode.Constants[vm.fetch()])
		case op.LoadFast:
			vm.Push(vm.activeFrame.Locals()[vm.fetch()])
		case op.LoadGlobal:
			vm.Push(vm.globals[vm.fetch()])
		case op.LoadFree:
			freeVars := vm.activeFrame.fn.FreeVars()
			vm.Push(freeVars[vm.fetch()].Value())
		case op.LoadBuiltin:
			vm.Push(vm.builtins[vm.fetch()])
		case op.StoreFast:
			vm.activeFrame.Locals()[vm.fetch()] = vm.Pop()
		case op.StoreGlobal:
			vm.globals[vm.fetch()] = vm.Pop()
		case op.StoreFree:
			freeVars := vm.activeFrame.fn.FreeVars()
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
					return errors.New("expected cell")
				}
			}
			fn := vm.activeCode.Constants[constIndex].(*object.Function)
			closure := object.NewClosure(fn, fn.Code(), free)
			vm.Push(closure)
		case op.MakeCell:
			symbolIndex := vm.fetch()
			framesBack := int(vm.fetch())
			frameIndex := vm.fp - framesBack
			if frameIndex < 0 {
				return fmt.Errorf("no frame at depth %d", framesBack)
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
			if err := vm.call(ctx, obj, argc); err != nil {
				return err
			}
		case op.Partial:
			argc := int(vm.fetch())
			args := make([]object.Object, argc)
			for i := argc - 1; i >= 0; i-- {
				args[i] = vm.Pop()
			}
			obj := vm.Pop()
			partial := object.NewPartial(obj, args)
			vm.Push(partial)
		case op.ReturnValue:
			returnAddr := vm.frames[vm.fp].returnAddr
			vm.fp--
			vm.activeFrame = &vm.frames[vm.fp]
			vm.activeCode = vm.activeFrame.Code()
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
				return fmt.Errorf("object is not a container: %T", obj)
			}
			result, err := container.GetItem(index)
			if err != nil {
				return err.Value()
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
				return fmt.Errorf("object is not a number: %T", obj)
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
				return fmt.Errorf("object is not a container: %T", container)
			}
		case op.Swap:
			vm.Swap(int(vm.fetch()))
		case op.BuildString:
			count := vm.fetch()
			items := make([]string, count)
			for i := uint16(0); i < count; i++ {
				dst := count - 1 - i
				obj := vm.Pop()
				switch obj := obj.(type) {
				case *object.Error:
					return obj.Value() // TODO: review this
				case *object.String:
					items[dst] = obj.Value()
				default:
					items[dst] = obj.Inspect()
				}
			}
			vm.Push(object.NewString(strings.Join(items, "")))
		case op.Range:
			container, ok := vm.Pop().(object.Container)
			if !ok {
				return fmt.Errorf("object is not a container: %T", container)
			}
			vm.Push(container.Iter())
		case op.Slice:
			start := vm.Pop()
			stop := vm.Pop()
			container, ok := vm.Pop().(object.Container)
			if !ok {
				return fmt.Errorf("object is not a container: %T", container)
			}
			slice := object.Slice{Start: start, Stop: stop}
			result, err := container.GetSlice(slice)
			if err != nil {
				return err.Value()
			}
			vm.Push(result)
		case op.Length:
			container, ok := vm.Pop().(object.Container)
			if !ok {
				return fmt.Errorf("object is not a container: %T", container)
			}
			vm.Push(container.Len())
		case op.Copy:
			offset := vm.fetch()
			vm.Push(vm.stack[vm.sp-int(offset)])
		case op.Import:
			name, ok := vm.Pop().(*object.String)
			if !ok {
				return fmt.Errorf("object is not a string: %T", name)
			}
			module, err := vm.loadModule(ctx, name.Value())
			if err != nil {
				return err
			}
			vm.Push(module)
		case op.Halt:
			return nil
		default:
			return fmt.Errorf("unknown opcode: %d", opcode)
		}
	}

	// If we reach this point and a return address is set, go there. This can
	// happen when importing a module completes, for example.
	if vm.activeFrame.returnAddr > 0 {
		fmt.Println("Deactivating module", vm.activeFrame.returnAddr)
		vm.fp--
		vm.ip = vm.activeFrame.returnAddr
		vm.activeFrame = &vm.frames[vm.fp]
		vm.activeCode = vm.activeFrame.code
		return nil
	}
	return nil
}

func (vm *VM) loadModule(ctx context.Context, name string) (*object.Module, error) {
	if module, ok := vm.modules[name]; ok {
		return module, nil
	}
	if vm.importer == nil {
		return nil, fmt.Errorf("imports are disabled")
	}
	module, err := vm.importer.Import(ctx, name)
	if err != nil {
		return nil, err
	}
	code := module.Code()
	vm.fp++
	frame := &vm.frames[vm.fp]
	frame.ActivateCode(code)
	frame.SetReturnAddr(vm.ip)
	fmt.Println("Activated module", name)
	vm.ip = 0
	vm.modules[name] = module
	return module, nil
}

func (vm *VM) call(ctx context.Context, fn object.Object, argc int) error {
	// The arguments are understood to be stored in vm.tmp here
	args := vm.tmp[:argc]
	switch fn := fn.(type) {
	case *object.Builtin:
		result := fn.Call(ctx, args...)
		vm.Push(result)
	case *object.Function:
		vm.fp++
		frame := &vm.frames[vm.fp]
		paramsCount := len(fn.Parameters())
		if err := checkCallArgs(fn, argc); err != nil {
			return err
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
		vm.activeFrame = frame
		vm.activeCode = code
		vm.ip = 0
	case *object.Partial:
		// Combine the current arguments with the partial's arguments
		expandedCount := argc + len(fn.Args())
		if expandedCount > MaxArgs {
			return fmt.Errorf("max arguments limit of %d exceeded (got %d)", MaxArgs, expandedCount)
		}
		// We can just append arguments from the partial into vm.tmp
		copy(vm.tmp[argc:], fn.Args())
		return vm.call(ctx, fn.Function(), expandedCount)
	default:
		return fmt.Errorf("object is not callable: %T", fn)
	}
	return nil
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
	return uint16(vm.activeCode.Instructions[ip])
}

// Calls a compiled function with the given arguments. This is used internally
// when a Tamarin object calls a function, e.g. [1, 2, 3].map(func(x) { x + 1 }).
func (vm *VM) callFunction(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
	// Advance to the next frame
	vm.fp++
	frame := &vm.frames[vm.fp]
	// Check that the argument count is appropriate
	paramsCount := len(fn.Parameters())
	argc := len(args)
	if err := checkCallArgs(fn, argc); err != nil {
		return nil, err
	}
	// Assemble frame local variables in vm.tmp. The local variable order is:
	// 1. Function parameters
	// 2. Function name (if the function is named)
	copy(vm.tmp[:argc], args)
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
	// Activate this new frame with the function code and local variables
	frame.ActivateFunction(fn, vm.ip, vm.tmp[:argc])
	vm.activeFrame = frame
	vm.activeCode = code
	vm.ip = 0
	// Evaluate the function code then return the result from TOS
	if err := vm.eval(ctx); err != nil {
		return nil, err
	}
	return vm.Pop(), nil
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

func concat(a, b []object.Object) []object.Object {
	aLen := len(a)
	bLen := len(b)
	result := make([]object.Object, aLen+bLen)
	copy(result, a)
	copy(result[aLen:], b)
	return result
}
