package vm

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
)

const (
	MaxArgs       = 255
	MaxFrameDepth = 1024
	MaxStackDepth = 1024
)

func Run(code string) (object.Object, error) {
	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}
	builtins := map[string]object.Object{}
	for _, b := range evaluator.GlobalBuiltins() {
		builtins[b.Key()] = b
	}
	c := compiler.New(compiler.Options{
		Builtins: builtins,
		Name:     "main",
	})
	mainScope, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}
	vm := New(mainScope)
	if err := vm.Run(); err != nil {
		return nil, err
	}
	return vm.Pop(), nil
}

type VM struct {
	ip           int // instruction pointer
	sp           int // stack pointer
	fp           int // frame pointer
	stack        [MaxStackDepth]object.Object
	frames       [MaxFrameDepth]Frame
	tmp          [MaxArgs]object.Object
	currentFrame *Frame
	main         *compiler.Scope
	currentScope *compiler.Scope
	globals      []object.Object
	builtins     []object.Object
}

func New(main *compiler.Scope) *VM {
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

func (vm *VM) Run() error {

	ctx := context.Background()

	// Initialize the call frame with the main function
	vm.fp++
	vm.ip++
	vm.currentFrame = &vm.frames[vm.fp]
	vm.currentFrame.Init(nil, 0, vm.currentScope.Symbols.Size())
	vm.currentFrame.scope = vm.main

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
				return fmt.Errorf("attribute %q not found", name)
			}
			vm.Push(value)
		case op.LoadConst:
			vm.Push(vm.currentScope.Constants[vm.fetch()])
		case op.LoadFast:
			vm.Push(vm.currentFrame.locals[vm.fetch()])
		case op.LoadGlobal:
			vm.Push(vm.globals[vm.fetch()])
		case op.LoadFree:
			freeVars := vm.currentFrame.fn.FreeVars()
			vm.Push(freeVars[vm.fetch()].Value())
		case op.LoadBuiltin:
			vm.Push(vm.builtins[vm.fetch()])
		case op.StoreFast:
			vm.currentFrame.locals[vm.fetch()] = vm.Pop()
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
					return errors.New("expected cell")
				}
			}
			fn := vm.currentScope.Constants[constIndex].(*object.CompiledFunction)
			closure := object.NewClosure(fn, fn.Scope(), free)
			vm.Push(closure)
		case op.MakeCell:
			symbolIndex := vm.fetch()
			framesBack := int(vm.fetch())
			frameIndex := vm.fp - framesBack
			if frameIndex < 0 {
				return fmt.Errorf("no frame at depth %d", framesBack)
			}
			frame := &vm.frames[frameIndex]
			locals := frame.Locals()
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
			vm.Push(vm.runCompareOp(opType, a, b))
		case op.BinaryOp:
			opType := op.BinaryOpType(vm.fetch())
			b := vm.Pop()
			a := vm.Pop()
			vm.Push(vm.runBinaryOp(opType, a, b))
		case op.Call:
			argc := int(vm.fetch())
			for argIndex := argc - 1; argIndex >= 0; argIndex-- {
				vm.tmp[argIndex] = vm.Pop()
			}
			obj := vm.Pop()
			switch obj := obj.(type) {
			case *object.Builtin:
				result := obj.Call(ctx, vm.tmp[:argc]...)
				vm.Push(result)
			case *object.CompiledFunction:
				vm.fp++
				frame := &vm.frames[vm.fp]
				scope := obj.Scope().(*compiler.Scope)
				if scope.IsNamed {
					vm.tmp[argc] = obj
					argc++
				}
				frame.InitWithLocals(obj, vm.ip, vm.tmp[:argc])
				vm.currentFrame = frame
				vm.currentScope = scope
				vm.ip = 0
			default:
				return fmt.Errorf("object is not callable: %T", obj)
			}
		case op.ReturnValue:
			returnAddr := vm.frames[vm.fp].returnAddr
			vm.fp--
			vm.currentFrame = &vm.frames[vm.fp]
			vm.currentScope = vm.currentFrame.Scope()
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
		case op.Halt:
			return nil
		default:
			return fmt.Errorf("unknown opcode: %d", opcode)
		}
	}
	return nil
}

func (vm *VM) runCompareOp(opType op.CompareOpType, a, b object.Object) object.Object {
	switch opType {
	case op.Equal:
		return a.Equals(b)
	case op.NotEqual:
		if a.Equals(b) == object.True {
			return object.False
		} else {
			return object.True
		}
	case op.LessThan:
		return object.NewBool(a.(*object.Int).Value() < b.(*object.Int).Value())
	case op.LessThanOrEqual:
		return object.NewBool(a.(*object.Int).Value() <= b.(*object.Int).Value())
	case op.GreaterThan:
		return object.NewBool(a.(*object.Int).Value() > b.(*object.Int).Value())
	case op.GreaterThanOrEqual:
		return object.NewBool(a.(*object.Int).Value() >= b.(*object.Int).Value())
	default:
		panic("unknown compare op")
	}
}

func (vm *VM) runBinaryOp(opType op.BinaryOpType, a, b object.Object) object.Object {
	switch opType {
	case op.Add:
		return object.NewInt(a.(*object.Int).Value() + b.(*object.Int).Value())
	case op.Subtract:
		return object.NewInt(a.(*object.Int).Value() - b.(*object.Int).Value())
	case op.Multiply:
		return object.NewInt(a.(*object.Int).Value() * b.(*object.Int).Value())
	case op.Divide:
		return object.NewInt(a.(*object.Int).Value() / b.(*object.Int).Value())
	case op.Modulo:
		return object.NewInt(a.(*object.Int).Value() % b.(*object.Int).Value())
	case op.And:
		return object.NewInt(a.(*object.Int).Value() & b.(*object.Int).Value())
	case op.Or:
		return object.NewInt(a.(*object.Int).Value() | b.(*object.Int).Value())
	case op.Xor:
		return object.NewInt(a.(*object.Int).Value() ^ b.(*object.Int).Value())
	case op.Power:
		return object.NewInt(int64(math.Pow(float64(a.(*object.Int).Value()), float64(b.(*object.Int).Value()))))
	case op.LShift:
		return object.NewInt(a.(*object.Int).Value() << b.(*object.Int).Value())
	case op.RShift:
		return object.NewInt(a.(*object.Int).Value() >> b.(*object.Int).Value())
	}
	return nil
}

func (vm *VM) TOS() object.Object {
	return vm.stack[vm.sp]
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

func (vm *VM) fetch() uint16 {
	ip := vm.ip
	vm.ip++
	return uint16(vm.currentScope.Instructions[ip])
}
