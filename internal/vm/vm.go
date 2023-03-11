package vm

import (
	"errors"
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
)

const MaxFrameCount = 1024

func Run(code string) (object.Object, error) {
	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}
	c := compiler.New(compiler.Options{})
	bytecode, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}
	vm := New(bytecode)
	if err := vm.Run(); err != nil {
		return nil, err
	}
	return vm.Pop(), nil
}

type VM struct {
	ip           int
	sp           int
	stack        *Stack[object.Object]
	frameStack   *Stack[*Frame]
	scopes       []*compiler.Scope
	currentScope *compiler.Scope
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		stack:        NewStack[object.Object](1024),
		frameStack:   NewStack[*Frame](1024),
		sp:           -1,
		scopes:       bytecode.Scopes,
		currentScope: bytecode.Scopes[0],
	}
}

func (vm *VM) Run() error {
	// for i, b := range vm.code {
	// 	fmt.Printf("%d %d\n", i, b)
	// }
	scope := vm.currentScope
	symbolCount := vm.currentScope.Symbols.Size()
	vm.frameStack.Push(NewFrame(nil, make([]object.Object, symbolCount), 0))
	for vm.ip < len(scope.Instructions) {
		opcode := scope.Instructions[vm.ip]
		opinfo := op.GetInfo(opcode)
		fmt.Println("IP:", vm.ip, "OPCODE:", opcode, "INFO:", opinfo)
		vm.ip++
		switch opcode {
		case op.Nop:
		case op.LoadConst:
			constIndex := vm.fetch2()
			vm.stack.Push(scope.Constants[constIndex])
		case op.StoreFast:
			obj := vm.Pop()
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
		case op.CompareOp:
			opType := op.CompareOpType(vm.fetch())
			b := vm.Pop()
			a := vm.Pop()
			vm.stack.Push(vm.runCompareOp(opType, a, b))
		case op.BinaryOp:
			opType := op.BinaryOpType(vm.fetch())
			b := vm.Pop()
			a := vm.Pop()
			vm.stack.Push(vm.runBinaryOp(opType, a, b))
		case op.Call:
			if vm.frameStack.Size() >= MaxFrameCount {
				return errors.New("stack overflow")
			}
			fn, ok := vm.Pop().(*object.Function)
			if !ok {
				return fmt.Errorf("not a function: %T", fn)
			}
			fnCode := fn.GetCodeStart()
			argc := vm.fetch()
			args := make([]object.Object, argc)
			for i := 0; i < argc; i++ {
				args[len(args)-1-i] = vm.Pop()
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
		case op.PopJumpForwardIfTrue:
			tos := vm.Pop()
			delta := vm.fetch2() - 3
			if tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpForwardIfFalse:
			tos := vm.Pop()
			delta := vm.fetch2() - 3
			if !tos.IsTruthy() {
				vm.ip += delta
			}
		case op.PopJumpBackwardIfTrue:
			tos := vm.Pop()
			delta := vm.fetch2() - 3
			if tos.IsTruthy() {
				vm.ip -= delta
			}
		case op.PopJumpBackwardIfFalse:
			tos := vm.Pop()
			delta := vm.fetch2() - 3
			if !tos.IsTruthy() {
				vm.ip -= delta
			}
		case op.JumpForward:
			base := vm.ip - 1
			delta := vm.fetch2()
			vm.ip = base + delta
		case op.JumpBackward:
			base := vm.ip - 1
			delta := vm.fetch2()
			vm.ip = base - delta
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

func (vm *VM) TOS() (object.Object, bool) {
	return vm.stack.Top()
}

func (vm *VM) Frame() (*Frame, bool) {
	return vm.frameStack.Top()
}

func (vm *VM) Pop() object.Object {
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
	return int(vm.currentScope.Instructions[ip])
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
