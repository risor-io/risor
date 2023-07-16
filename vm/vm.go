package vm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const (
	MaxArgs       = 255
	MaxFrameDepth = 1024
	MaxStackDepth = 1024
	StopSignal    = -1
	MB            = 1024 * 1024
)

type Options struct {
	Main              *object.Code
	InstructionOffset int
	Importer          importer.Importer
}

type VirtualMachine struct {
	ip          int // instruction pointer
	sp          int // stack pointer
	fp          int // frame pointer
	halt        int32
	stack       [MaxStackDepth]object.Object
	frames      [MaxFrameDepth]Frame
	tmp         [MaxArgs]object.Object
	activeFrame *Frame
	activeCode  *object.Code
	main        *object.Code
	importer    importer.Importer
	modules     map[string]*object.Module
	limits      limits.Limits
}

// Option is a configuration function for a Virtual Machine.
type Option func(*VirtualMachine)

// WithInstructionOffset sets the initial instruction offset.
func WithInstructionOffset(offset int) Option {
	return func(vm *VirtualMachine) {
		vm.ip = offset
	}
}

// WithImporter is used to supply an Importer to the Virtual Machine.
func WithImporter(importer importer.Importer) Option {
	return func(vm *VirtualMachine) {
		vm.importer = importer
	}
}

// WithLimits sets the limits for the Virtual Machine.
func WithLimits(limits limits.Limits) Option {
	return func(vm *VirtualMachine) {
		vm.limits = limits
	}
}

func defaultLimits() limits.Limits {
	return limits.New(limits.WithMaxBufferSize(100 * MB))
}

// New creates a new Virtual Machine.
func New(main *object.Code, options ...Option) *VirtualMachine {
	vm := &VirtualMachine{
		sp:      -1,
		ip:      0,
		main:    main,
		modules: map[string]*object.Module{},
	}
	for _, opt := range options {
		opt(vm)
	}
	if vm.limits == nil {
		vm.limits = defaultLimits()
	}
	return vm
}

func (vm *VirtualMachine) Run(ctx context.Context) (err error) {

	// Translate any panic into an error so the caller has a good guarantee
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Halt execution when the context is cancelled
	go func() {
		<-ctx.Done()
		atomic.StoreInt32(&vm.halt, 1)
	}()

	// Activate the "main" entrypoint code in frame 0 and then run it
	vm.fp = 0
	vm.activeFrame = &vm.frames[vm.fp]
	vm.activeFrame.ActivateCode(vm.main)
	vm.activeCode = vm.main
	ctx = object.WithCallFunc(ctx, vm.callFunction)
	ctx = object.WithCodeFunc(ctx, vm.codeFunction)
	ctx = limits.WithLimits(ctx, vm.limits)
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
func (vm *VirtualMachine) eval(ctx context.Context) error {

	// Run to the end of the active code
	for vm.ip < len(vm.activeCode.Instructions) {

		if atomic.LoadInt32(&vm.halt) == 1 {
			return ctx.Err()
		}

		// The current instruction opcode
		opcode := vm.activeCode.Instructions[vm.ip]

		// fmt.Println("ip", vm.ip, op.GetInfo(opcode).Name, "sp", vm.sp)

		// Advance the instruction pointer to the next instruction. Note that
		// this is done before we actually execute the current instruction, so
		// relative jump instructions will need to take this into account.
		vm.ip++

		// Dispatch the instruction
		switch opcode {
		case op.Nop:
		case op.LoadAttr:
			obj := vm.pop()
			name := vm.activeCode.Names[vm.fetch()]
			value, found := obj.GetAttr(name)
			if !found {
				return fmt.Errorf("exec error: attribute %q not found on %s object",
					name, obj.Type())
			}
			vm.push(value)
		case op.LoadConst:
			vm.push(vm.activeCode.Constants[vm.fetch()])
		case op.LoadFast:
			vm.push(vm.activeFrame.Locals()[vm.fetch()])
		case op.LoadGlobal:
			vm.push(vm.activeCode.Globals()[vm.fetch()])
		case op.LoadFree:
			freeVars := vm.activeFrame.fn.FreeVars()
			vm.push(freeVars[vm.fetch()].Value())
		case op.LoadBuiltin:
			vm.push(vm.activeCode.Builtins()[vm.fetch()])
		case op.StoreFast:
			vm.activeFrame.Locals()[vm.fetch()] = vm.pop()
		case op.StoreGlobal:
			vm.activeCode.Globals()[vm.fetch()] = vm.pop()
		case op.StoreFree:
			freeVars := vm.activeFrame.fn.FreeVars()
			freeVars[vm.fetch()].Set(vm.pop())
		case op.LoadClosure:
			constIndex := vm.fetch()
			freeCount := vm.fetch()
			free := make([]*object.Cell, freeCount)
			for i := uint16(0); i < freeCount; i++ {
				obj := vm.pop()
				switch obj := obj.(type) {
				case *object.Cell:
					free[i] = obj
				default:
					return errors.New("exec error: expected cell")
				}
			}
			fn := vm.activeCode.Constants[constIndex].(*object.Function)
			closure := object.NewClosure(fn, fn.Code(), free)
			vm.push(closure)
		case op.MakeCell:
			symbolIndex := vm.fetch()
			framesBack := int(vm.fetch())
			frameIndex := vm.fp - framesBack
			if frameIndex < 0 {
				return fmt.Errorf("exec error: no frame at depth %d", framesBack)
			}
			frame := &vm.frames[frameIndex]
			locals := frame.CaptureLocals()
			vm.push(object.NewCell(&locals[symbolIndex]))
		case op.Nil:
			vm.push(object.Nil)
		case op.True:
			vm.push(object.True)
		case op.False:
			vm.push(object.False)
		case op.CompareOp:
			opType := op.CompareOpType(vm.fetch())
			b := vm.pop()
			a := vm.pop()
			vm.push(object.Compare(opType, a, b))
		case op.BinaryOp:
			opType := op.BinaryOpType(vm.fetch())
			b := vm.pop()
			a := vm.pop()
			vm.push(object.BinaryOp(opType, a, b))
		case op.Call:
			argc := int(vm.fetch())
			for argIndex := argc - 1; argIndex >= 0; argIndex-- {
				vm.tmp[argIndex] = vm.pop()
			}
			obj := vm.pop()
			if err := vm.call(ctx, obj, argc); err != nil {
				return err
			}
		case op.Partial:
			argc := int(vm.fetch())
			args := make([]object.Object, argc)
			for i := argc - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			obj := vm.pop()
			partial := object.NewPartial(obj, args)
			vm.push(partial)
		case op.ReturnValue:
			returnAddr := vm.frames[vm.fp].returnAddr
			vm.fp--
			vm.activeFrame = &vm.frames[vm.fp]
			vm.activeCode = vm.activeFrame.Code()
			vm.ip = returnAddr
			if returnAddr == StopSignal {
				// If StopSignal is found as the return address, it means the
				// current eval call should stop.
				return nil
			}
		case op.PopJumpForwardIfTrue:
			tos := vm.pop()
			delta := int(vm.fetch()) - 2
			if tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpForwardIfFalse:
			tos := vm.pop()
			delta := int(vm.fetch()) - 2
			if !tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpBackwardIfTrue:
			tos := vm.pop()
			delta := int(vm.fetch()) - 2
			if tos.IsTruthy() {
				vm.ip -= delta
			}
		case op.PopJumpBackwardIfFalse:
			tos := vm.pop()
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
				items[count-1-i] = vm.pop()
			}
			vm.push(object.NewList(items))
		case op.BuildMap:
			count := vm.fetch()
			items := make(map[string]object.Object, count)
			for i := uint16(0); i < count; i++ {
				v := vm.pop()
				k := vm.pop()
				items[k.(*object.String).Value()] = v
			}
			vm.push(object.NewMap(items))
		case op.BuildSet:
			count := vm.fetch()
			items := make([]object.Object, count)
			for i := uint16(0); i < count; i++ {
				items[i] = vm.pop()
			}
			vm.push(object.NewSet(items))
		case op.BinarySubscr:
			idx := vm.pop()
			lhs := vm.pop()
			container, ok := lhs.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)", lhs.Type())
			}
			result, err := container.GetItem(idx)
			if err != nil {
				return err.Value()
			}
			vm.push(result)
		case op.StoreSubscr:
			idx := vm.pop()
			lhs := vm.pop()
			rhs := vm.pop()
			container, ok := lhs.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)", lhs.Type())
			}
			if err := container.SetItem(idx, rhs); err != nil {
				return err.Value()
			}
		case op.UnaryNegative:
			obj := vm.pop()
			switch obj := obj.(type) {
			case *object.Int:
				vm.push(object.NewInt(-obj.Value()))
			case *object.Float:
				vm.push(object.NewFloat(-obj.Value()))
			default:
				return fmt.Errorf("type error: object is not a number (got %s)", obj.Type())
			}
		case op.UnaryNot:
			obj := vm.pop()
			if obj.IsTruthy() {
				vm.push(object.False)
			} else {
				vm.push(object.True)
			}
		case op.ContainsOp:
			obj := vm.pop()
			containerObj := vm.pop()
			invert := vm.fetch() == 1
			if container, ok := containerObj.(object.Container); ok {
				value := container.Contains(obj)
				if invert {
					value = object.Not(value)
				}
				vm.push(value)
			} else {
				return fmt.Errorf("type error: object is not a container (got %s)",
					containerObj.Type())
			}
		case op.Swap:
			vm.swap(int(vm.fetch()))
		case op.BuildString:
			count := vm.fetch()
			items := make([]string, count)
			for i := uint16(0); i < count; i++ {
				dst := count - 1 - i
				obj := vm.pop()
				switch obj := obj.(type) {
				case *object.Error:
					return obj.Value() // TODO: review this
				case *object.String:
					items[dst] = obj.Value()
				default:
					items[dst] = obj.Inspect()
				}
			}
			vm.push(object.NewString(strings.Join(items, "")))
		case op.Range:
			containerObj := vm.pop()
			container, ok := containerObj.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)",
					containerObj.Type())
			}
			vm.push(container.Iter())
		case op.Slice:
			start := vm.pop()
			stop := vm.pop()
			containerObj := vm.pop()
			container, ok := containerObj.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)",
					containerObj.Type())
			}
			slice := object.Slice{Start: start, Stop: stop}
			result, err := container.GetSlice(slice)
			if err != nil {
				return err.Value()
			}
			vm.push(result)
		case op.Length:
			containerObj := vm.pop()
			container, ok := containerObj.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)",
					containerObj.Type())
			}
			vm.push(container.Len())
		case op.Copy:
			offset := vm.fetch()
			vm.push(vm.stack[vm.sp-int(offset)])
		case op.Import:
			name, ok := vm.pop().(*object.String)
			if !ok {
				return fmt.Errorf("type error: object is not a string (got %s)", name.Type())
			}
			module, err := vm.loadModule(ctx, name.Value())
			if err != nil {
				return err
			}
			vm.push(module)
		case op.PopTop:
			vm.pop()
		case op.Unpack:
			containerObj := vm.pop()
			nameCount := int64(vm.fetch())
			container, ok := containerObj.(object.Container)
			if !ok {
				return fmt.Errorf("type error: object is not a container (got %s)",
					containerObj.Type())
			}
			containerSize := container.Len().Value()
			if containerSize != nameCount {
				return fmt.Errorf("exec error: unpack count mismatch: %d != %d",
					containerSize, nameCount)
			}
			iter := container.Iter()
			for {
				val, ok := iter.Next()
				if !ok {
					break
				}
				vm.push(val)
			}
		case op.GetIter:
			obj := vm.pop()
			switch obj := obj.(type) {
			case object.Container:
				vm.push(obj.Iter())
			case object.Iterator:
				vm.push(obj)
			default:
				return fmt.Errorf("type error: object is not iterable (got %s)", obj.Type())
			}
		case op.ForIter:
			base := vm.ip - 1
			jumpAmount := vm.fetch()
			nameCount := vm.fetch()
			iter := vm.pop().(object.Iterator)
			if _, ok := iter.Next(); !ok {
				vm.ip = base + int(jumpAmount)
			} else {
				obj, _ := iter.Entry()
				vm.push(iter)
				if nameCount == 1 {
					vm.push(obj.Key())
				} else if nameCount == 2 {
					vm.push(obj.Value())
					vm.push(obj.Key())
				} else if nameCount != 0 {
					return fmt.Errorf("exec error: invalid iteration")
				}
			}
		case op.Halt:
			return nil
		default:
			return fmt.Errorf("exec error: unknown opcode: %d", opcode)
		}
	}
	return nil
}

func (vm *VirtualMachine) loadModule(ctx context.Context, name string) (*object.Module, error) {
	if module, ok := vm.modules[name]; ok {
		return module, nil
	}
	if vm.importer == nil {
		return nil, fmt.Errorf("exec error: imports are disabled")
	}
	// Load and compile the module code
	module, err := vm.importer.Import(ctx, name)
	if err != nil {
		return nil, err
	}
	code := module.Code()

	// Allocate a new frame to evaluate the module code
	baseFrame := vm.fp
	vm.fp++
	frame := &vm.frames[vm.fp]
	frame.ActivateCode(code)
	frame.SetReturnAddr(vm.ip)
	vm.activeFrame = &vm.frames[vm.fp]
	vm.activeCode = vm.activeFrame.code
	vm.ip = 0

	// Evaluate the module code
	if err := vm.eval(ctx); err != nil {
		// Unwind the stack
		vm.fp = baseFrame
		vm.activeFrame = &vm.frames[vm.fp]
		vm.activeCode = vm.activeFrame.code
		vm.ip = frame.returnAddr
		return nil, err
	}

	// Resume the previous frame
	vm.fp--
	vm.ip = vm.activeFrame.returnAddr
	vm.activeFrame = &vm.frames[vm.fp]
	vm.activeCode = vm.activeFrame.code

	// Cache the module
	vm.modules[name] = module
	return module, nil
}

func (vm *VirtualMachine) call(ctx context.Context, fn object.Object, argc int) error {
	// The arguments are understood to be stored in vm.tmp here
	args := vm.tmp[:argc]
	switch fn := fn.(type) {
	case *object.Builtin:
		result := fn.Call(ctx, args...)
		if err, ok := result.(*object.Error); ok {
			return err.Value()
		}
		vm.push(result)
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
			argc = paramsCount
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
			return fmt.Errorf("exec error: max arguments limit of %d exceeded (got %d)", MaxArgs, expandedCount)
		}
		// We can just append arguments from the partial into vm.tmp
		copy(vm.tmp[argc:], fn.Args())
		return vm.call(ctx, fn.Function(), expandedCount)
	default:
		return fmt.Errorf("type error: object is not callable (got %s)", fn.Type())
	}
	return nil
}

func (vm *VirtualMachine) TOS() (object.Object, bool) {
	if vm.sp >= 0 {
		return vm.stack[vm.sp], true
	}
	return nil, false
}

func (vm *VirtualMachine) pop() object.Object {
	obj := vm.stack[vm.sp]
	vm.sp--
	return obj
}

func (vm *VirtualMachine) push(obj object.Object) {
	vm.sp++
	vm.stack[vm.sp] = obj
}

func (vm *VirtualMachine) swap(pos int) {
	otherIndex := vm.sp - pos
	tos := vm.stack[vm.sp]
	other := vm.stack[otherIndex]
	vm.stack[otherIndex] = tos
	vm.stack[vm.sp] = other
}

func (vm *VirtualMachine) fetch() uint16 {
	ip := vm.ip
	vm.ip++
	return uint16(vm.activeCode.Instructions[ip])
}

func (vm *VirtualMachine) codeFunction(ctx context.Context) (*object.Code, error) {
	return vm.activeCode, nil
}

// Calls a compiled function with the given arguments. This is used internally
// when a Risor object calls a function, e.g. [1, 2, 3].map(func(x) { x + 1 }).
func (vm *VirtualMachine) callFunction(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
	baseFrame := vm.fp
	baseIP := vm.ip
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
	frame.ActivateFunction(fn, StopSignal, vm.tmp[:argc])
	vm.activeFrame = frame
	vm.activeCode = code
	vm.ip = 0
	// Evaluate the function code then return the result from TOS
	if err := vm.eval(ctx); err != nil {
		// Unwind the stack
		vm.fp = baseFrame
		vm.activeFrame = &vm.frames[vm.fp]
		vm.activeCode = vm.activeFrame.code
		vm.ip = baseIP
		return nil, err
	}
	vm.ip = baseIP
	value := vm.pop()
	return value, nil
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
