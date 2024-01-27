// Package vm provides a VirtualMachine that executes compiled Risor code.
package vm

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/risor-io/risor/compiler"
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

type VirtualMachine struct {
	ip           int // instruction pointer
	sp           int // stack pointer
	fp           int // frame pointer
	halt         int32
	stack        [MaxStackDepth]object.Object
	frames       [MaxFrameDepth]frame
	tmp          [MaxArgs]object.Object
	activeFrame  *frame
	activeCode   *code
	main         *compiler.Code
	importer     importer.Importer
	modules      map[string]*object.Module
	inputGlobals map[string]any
	globals      map[string]object.Object
	limits       limits.Limits
	loadedCode   map[*compiler.Code]*code
	running      bool
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

// WithGlobals provides global variables with the given names.
func WithGlobals(globals map[string]any) Option {
	return func(vm *VirtualMachine) {
		for name, value := range globals {
			vm.inputGlobals[name] = value
		}
	}
}

func defaultLimits() limits.Limits {
	return limits.New(limits.WithMaxBufferSize(100 * MB))
}

// Run the given code in a new Virtual Machine and return the result.
func Run(ctx context.Context, main *compiler.Code, options ...Option) (object.Object, error) {
	machine := New(main, options...)
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

// New creates a new Virtual Machine.
func New(main *compiler.Code, options ...Option) *VirtualMachine {
	vm := &VirtualMachine{
		sp:           -1,
		ip:           0,
		main:         main,
		modules:      map[string]*object.Module{},
		inputGlobals: map[string]any{},
		globals:      map[string]object.Object{},
		loadedCode:   map[*compiler.Code]*code{},
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
	if doneChan := ctx.Done(); doneChan != nil {
		go func() {
			<-doneChan
			atomic.StoreInt32(&vm.halt, 1)
		}()
	}

	// Convert globals
	vm.globals, err = object.AsObjects(vm.inputGlobals)
	if err != nil {
		return err
	}

	// Add any globals that are modules cache
	for name, value := range vm.globals {
		if module, ok := value.(*object.Module); ok {
			vm.modules[name] = module
		}
	}

	// Keep `running` flag up-to-date
	vm.running = true
	defer func() { vm.running = false }()

	// Activate the "main" entrypoint code in frame 0 and then run it
	var code *code
	if len(vm.loadedCode) > 0 {
		code = vm.reload(vm.main)
	} else {
		code = vm.load(vm.main)
	}
	vm.activateCode(0, vm.ip, code)
	ctx = object.WithCallFunc(ctx, vm.callFunction)
	ctx = limits.WithLimits(ctx, vm.limits)
	err = vm.eval(ctx)
	return
}

// Get a global variable by name as a Risor Object. Returns an error if the
// variable can't be found.
func (vm *VirtualMachine) Get(name string) (object.Object, error) {
	code := vm.activeCode
	if code == nil {
		return nil, errors.New("no active code")
	}
	for i := 0; i < code.GlobalsCount(); i++ {
		if g := code.Global(i); g.Name() == name {
			return code.Globals[g.Index()], nil
		}
	}
	return nil, fmt.Errorf("global with name %q not found", name)
}

// GlobalNames returns the names of all global variables in the active code.
func (vm *VirtualMachine) GlobalNames() []string {
	if vm.activeCode == nil {
		return nil
	}
	count := vm.activeCode.GlobalsCount()
	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		names = append(names, vm.activeCode.Global(i).Name())
	}
	return names
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
			switch value := value.(type) {
			case object.AttrResolver:
				attr, err := value.ResolveAttr(ctx, name)
				if err != nil {
					return err
				}
				vm.push(attr)
			default:
				vm.push(value)
			}
		case op.LoadConst:
			vm.push(vm.activeCode.Constants[vm.fetch()])
		case op.LoadFast:
			vm.push(vm.activeFrame.Locals()[vm.fetch()])
		case op.LoadGlobal:
			vm.push(vm.activeCode.Globals[vm.fetch()])
		case op.LoadFree:
			freeVars := vm.activeFrame.fn.FreeVars()
			vm.push(freeVars[vm.fetch()].Value())
		case op.StoreFast:
			vm.activeFrame.Locals()[vm.fetch()] = vm.pop()
		case op.StoreGlobal:
			vm.activeCode.Globals[vm.fetch()] = vm.pop()
		case op.StoreFree:
			freeVars := vm.activeFrame.fn.FreeVars()
			freeVars[vm.fetch()].Set(vm.pop())
		case op.StoreAttr:
			idx := vm.fetch()
			obj := vm.pop()
			value := vm.pop()
			name := vm.activeCode.Names[idx]
			if err := obj.SetAttr(name, value); err != nil {
				return err
			}
		case op.LoadClosure:
			constIndex := vm.fetch()
			freeCount := vm.fetch()
			free := make([]*object.Cell, freeCount)
			for i := uint16(0); i < freeCount; i++ {
				obj := vm.pop()
				switch obj := obj.(type) {
				case *object.Cell:
					free[freeCount-i-1] = obj
				default:
					return errors.New("exec error: expected cell")
				}
			}
			fn := vm.activeCode.Constants[constIndex].(*object.Function)
			vm.push(object.NewClosure(fn, free))
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
			activeFrame := vm.activeFrame
			returnAddr := activeFrame.returnAddr
			returnSp := activeFrame.returnSp
			vm.resumeFrame(vm.fp-1, returnAddr, returnSp)
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
			iterableObj := vm.pop()
			iterable, ok := iterableObj.(object.Iterable)
			if !ok {
				return fmt.Errorf("type error: object is not an iterable (got %s)",
					iterableObj.Type())
			}
			vm.push(iterable.Iter())
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
		case op.FromImport:
			parentLen := vm.fetch()
			nameLen := vm.fetch()
			if nameLen != 1 {
				return fmt.Errorf("exec error: from-import name length is not 1: %d", nameLen)
			}
			name, ok := vm.pop().(*object.String)
			if !ok {
				return fmt.Errorf("type error: object is not a string (got %s)", name.Type())
			}
			from := make([]string, parentLen)
			for i := int(parentLen - 1); i >= 0; i-- {
				val, ok := vm.pop().(*object.String)
				if !ok {
					return fmt.Errorf("type error: object is not a string (got %s)", val.Type())
				}
				from[i] = val.Value()
			}
			// name is a real module name
			module, err := vm.loadModule(ctx, filepath.Join(filepath.Join(from...), name.Value()))
			if err == nil {
				vm.push(module)
			} else {
				// name is a symbol
				module, err := vm.loadModule(ctx, filepath.Join(from...))
				if err != nil {
					return err
				}
				attr, found := module.GetAttr(name.String())
				if !found {
					return fmt.Errorf("import error: cannot import name %q from %q",
						name.String(), module.Name())
				}
				vm.push(attr)
			}
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
			case object.Iterable:
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
	// Activate a new frame to evaluate the module code
	baseFP := vm.fp
	baseIP := vm.ip
	baseSP := vm.sp
	code := vm.load(module.Code())
	vm.activateCode(vm.fp+1, 0, code)
	// Restore the previous frame when done
	defer vm.resumeFrame(baseFP, baseIP, baseSP)
	// Evaluate the module code
	if err := vm.eval(ctx); err != nil {
		return nil, err
	}
	module.UseGlobals(code.Globals)
	// Cache the module
	vm.modules[name] = module
	return module, nil
}

func (vm *VirtualMachine) call(ctx context.Context, fn object.Object, argc int) error {
	// The arguments are understood to be stored in vm.tmp here
	args := vm.tmp[:argc]
	switch fn := fn.(type) {
	case *object.Function:
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
		if code.IsNamed() {
			vm.tmp[paramsCount] = fn
			argc++
		}
		vm.activateFunction(vm.fp+1, 0, fn, vm.tmp[:argc])
	case *object.Partial:
		// Combine the current arguments with the partial's arguments
		expandedCount := argc + len(fn.Args())
		if expandedCount > MaxArgs {
			return fmt.Errorf("exec error: max arguments limit of %d exceeded (got %d)", MaxArgs, expandedCount)
		}
		// We can just append arguments from the partial into vm.tmp
		copy(vm.tmp[argc:], fn.Args())
		return vm.call(ctx, fn.Function(), expandedCount)
	case object.Callable:
		result := fn.Call(ctx, args...)
		if err, ok := result.(*object.Error); ok {
			return err.Value()
		}
		vm.push(result)
	default:
		return fmt.Errorf("type error: object is not callable (got %s)", fn.Type())
	}
	return nil
}

// GetIP returns the current instruction pointer.
func (vm *VirtualMachine) GetIP() int {
	return vm.ip
}

// SetIP sets the current instruction pointer.
func (vm *VirtualMachine) SetIP(value int) {
	vm.ip = value
}

func (vm *VirtualMachine) TOS() (object.Object, bool) {
	if vm.sp >= 0 {
		return vm.stack[vm.sp], true
	}
	return nil, false
}

func (vm *VirtualMachine) pop() object.Object {
	obj := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil
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

// Call a function with the supplied arguments. If isolation between VMs is
// important to you, do not provide a function here that was obtained from
// another VM, since it could be a closure over variables in that VM. This
// method should only be called after this VM stops running. Otherwise, an
// error is returned.
func (vm *VirtualMachine) Call(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
	if vm.running {
		return nil, errors.New("exec error: cannot call function while the vm is running")
	}
	return vm.callFunction(ctx, fn, args)
}

// Calls a compiled function with the given arguments. This is used internally
// when a Risor object calls a function, e.g. [1, 2, 3].map(func(x) { x + 1 }).
func (vm *VirtualMachine) callFunction(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
	baseFP := vm.fp
	baseIP := vm.ip
	baseSP := vm.sp

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
		argc = paramsCount
	}
	code := fn.Code()
	if code.IsNamed() {
		vm.tmp[paramsCount] = fn
		argc++
	}

	// Activate a frame for the function call
	vm.activateFunction(vm.fp+1, 0, fn, vm.tmp[:argc])

	// Setting StopSignal as the return address will cause the eval function to
	// stop execution when it reaches the end of the active code.
	vm.activeFrame.returnAddr = StopSignal

	// Restore the previous frame when done
	defer vm.resumeFrame(baseFP, baseIP, baseSP)

	// Evaluate the function code then return the result from TOS
	if err := vm.eval(ctx); err != nil {
		return nil, err
	}
	return vm.pop(), nil
}

// Wrap the *compiler.Code in a *code object to make it usable by the VM.
func (vm *VirtualMachine) load(cc *compiler.Code) *code {
	if code, ok := vm.loadedCode[cc]; ok {
		return code
	}
	// Loading is slightly different if this is the "root" (entrypoint) code
	// vs. a child of that. The root code owns the globals array, while the
	// children will reuse the globals from the root.
	rootCompiled := cc.Root()
	if rootCompiled == cc {
		c := loadRootCode(cc, vm.globals)
		vm.loadedCode[cc] = c
		return c
	}
	rootLoaded := vm.load(rootCompiled)
	c := loadChildCode(rootLoaded, cc)
	vm.loadedCode[cc] = c
	return c
}

// Reloads the main code while preserving global variables.
func (vm *VirtualMachine) reload(main *compiler.Code) *code {
	oldWrappedMain, ok := vm.loadedCode[main]
	if !ok {
		panic("main code not loaded")
	}
	vm.loadedCode = map[*compiler.Code]*code{}
	newWrappedMain := vm.load(main)
	copy(newWrappedMain.Globals, oldWrappedMain.Globals)
	return newWrappedMain
}

// Resume the frame at the given frame pointer, restoring the given IP and SP.
func (vm *VirtualMachine) resumeFrame(fp, ip, sp int) *frame {
	// The return value of the previous frame is on the top of the stack
	var frameResult object.Object = nil
	if vm.sp > sp {
		frameResult = vm.pop()
	}
	// Remove any items left on the stack by the previous frame
	for i := vm.sp; i > sp; i-- {
		vm.stack[i] = nil
	}
	vm.sp = sp
	// Push the frame result back onto the stack
	if frameResult != nil {
		vm.push(frameResult)
	}
	// Activate the resumed frame
	vm.fp = fp
	vm.ip = ip
	vm.activeFrame = &vm.frames[fp]
	vm.activeCode = vm.activeFrame.code
	return vm.activeFrame
}

// Activate a frame with the given code. This is typically used to begin
// running the entrypoint for a module or script.
func (vm *VirtualMachine) activateCode(fp, ip int, code *code) *frame {
	vm.fp = fp
	vm.ip = ip
	vm.activeFrame = &vm.frames[fp]
	vm.activeFrame.ActivateCode(code)
	vm.activeCode = code
	return vm.activeFrame
}

// Activate a frame with the given function, to implement a function call.
func (vm *VirtualMachine) activateFunction(fp, ip int, fn *object.Function, locals []object.Object) *frame {
	code := vm.load(fn.Code())
	returnAddr := vm.ip
	returnSp := vm.sp
	vm.fp = fp
	vm.ip = ip
	vm.activeFrame = &vm.frames[fp]
	vm.activeFrame.ActivateFunction(fn, code, returnAddr, returnSp, locals)
	vm.activeCode = code
	return vm.activeFrame
}

// Clone a stopped Virtual Machine. An error is returned if the Virtual Machine
// is running. The returned clone will have the same code, globals, modules,
// stack, and limits as the original. Any Risor objects present as global
// variables will be carried over to the clone. And since multiple clones can be
// created, the caller is responsible for ensuring that the global variables are
// not mutated, or that the mutations are safe for concurrent use. Do not use
// this function if you need isolation between clones, as this is not provided.
func (vm *VirtualMachine) Clone() (*VirtualMachine, error) {
	if vm.running {
		return nil, errors.New("cannot clone while the vm is running")
	}
	if vm.fp != 0 {
		return nil, errors.New("cannot clone while a frame is active")
	}
	modules := map[string]*object.Module{}
	for name, module := range vm.modules {
		modules[name] = module
	}
	inputGlobals := map[string]any{}
	for name, value := range vm.inputGlobals {
		inputGlobals[name] = value
	}
	globals := map[string]object.Object{}
	for name, value := range vm.globals {
		globals[name] = value
	}
	loadedCode := map[*compiler.Code]*code{}
	for code, loaded := range vm.loadedCode {
		loadedCode[code] = loaded.Clone()
	}
	clone := &VirtualMachine{
		sp:           vm.sp,
		ip:           vm.ip,
		fp:           0,
		stack:        vm.stack,
		main:         vm.main,
		importer:     vm.importer,
		modules:      modules,
		inputGlobals: inputGlobals,
		globals:      globals,
		limits:       vm.limits,
		loadedCode:   loadedCode,
	}
	clone.activateCode(0, vm.ip, loadedCode[vm.main])
	return clone, nil
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
