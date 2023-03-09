package compiler

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

type Bytecode struct {
	Instructions []op.Code
	Constants    []object.Object
	Symbols      *SymbolTable
}

type Compiler struct {
	symbols      *SymbolTable
	constants    []object.Object
	instructions []op.Code
}

type Options struct {
	Builtins []*object.Builtin
}

func New(opts Options) *Compiler {

	symbols := NewSymbolTable()

	for _, b := range opts.Builtins {
		symbols.Insert(b.Name(), SymbolAttrs{
			IsBuiltin: true,
			// Type:      string(b.Type()),
		})
	}

	return &Compiler{
		symbols: symbols,
	}
}

func (c *Compiler) Symbols() *SymbolTable {
	return c.symbols
}

func (c *Compiler) Instructions() []op.Code {
	return c.instructions
}

func (c *Compiler) Constants() []object.Object {
	return c.constants
}

func (c *Compiler) Compile(node ast.Node) (*Bytecode, error) {
	if err := c.compile(node); err != nil {
		return nil, err
	}
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
		Symbols:      c.symbols,
	}, nil
}

func (c *Compiler) compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Nil:
		c.emit(node, op.Nil)
	case *ast.Int:
		c.emit(node, op.LoadConst, c.constant(object.NewInt(node.Value())))
	case *ast.Float:
		c.emit(node, op.LoadConst, c.constant(object.NewFloat(node.Value())))
	case *ast.String:
		c.emit(node, op.LoadConst, c.constant(object.NewString(node.Value())))
	case *ast.Bool:
		if node.Value() {
			c.emit(node, op.True)
		} else {
			c.emit(node, op.False)
		}
	case *ast.If:
		if err := c.compileIf(node); err != nil {
			return err
		}
	case *ast.Infix:
		if err := c.compileInfix(node); err != nil {
			return err
		}
	case *ast.Program:
		for _, stmt := range node.Statements() {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
	case *ast.Block:
		// TODO: implement behavior for block specific variables
		for _, stmt := range node.Statements() {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
	case *ast.Var:
		name, expr := node.Value()
		if err := c.compile(expr); err != nil {
			return err
		}
		symbol, err := c.symbols.Insert(name, SymbolAttrs{})
		if err != nil {
			return err
		}
		c.emit(node, op.StoreFast, symbol.Index)
	case *ast.Assign:
		name := node.Name()
		expr := node.Value()
		if err := c.compile(expr); err != nil {
			return err
		}
		symbol, found := c.symbols.Lookup(name)
		fmt.Println("ASSIGN", name, symbol, found, symbol.Scope)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch symbol.Scope {
		case ScopeGlobal:
			c.emit(node, op.StoreGlobal, symbol.Index)
		case ScopeLocal:
			c.emit(node, op.StoreFast, symbol.Index)
			fmt.Println("EMIT", node, op.StoreFast, symbol.Index)
		}
	case *ast.Ident:
		name := node.Literal()
		symbol, found := c.symbols.Lookup(name)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch symbol.Scope {
		// case ScopeBuiltin:
		// 	c.emit(node, op.LoadBuiltin, symbol.Index)
		// case ScopeFree:
		// 	c.emit(node, op.LoadFree, symbol.Index)
		case ScopeGlobal:
			c.emit(node, op.LoadGlobal, symbol.Index)
		case ScopeLocal:
			c.emit(node, op.LoadFast, symbol.Index)
		}
	default:
		fmt.Println("DEFAULT", node, reflect.TypeOf(node))
	}
	// panic(fmt.Sprintf("unknown ast node type: %T", node))
	return nil
}

func (c *Compiler) compileIf(node *ast.If) error {
	if err := c.compile(node.Condition()); err != nil {
		return err
	}
	jumpIfFalsePos := c.emit(node, op.PopJumpForwardIfFalse, 9999)
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	alternative := node.Alternative()
	if alternative != nil {
		// Jump forward to skip the alternative by default
		jumpForwardPos := c.emit(node, op.JumpForward, 9999)

		// Update PopJumpForwardIfFalse to point to this alternative,
		// so that the alternative is executed if the condition is false
		delta := c.calculateDelta(jumpIfFalsePos)
		c.changeOperand2(jumpIfFalsePos, delta)

		if err := c.compile(alternative); err != nil {
			return err
		}
		c.changeOperand2(jumpForwardPos, c.calculateDelta(jumpForwardPos))
	} else {
		delta := c.calculateDelta(jumpIfFalsePos)
		c.changeOperand2(jumpIfFalsePos, delta)
	}
	return nil
}

func (c *Compiler) calculateDelta(pos int) int {
	return len(c.instructions) - pos
}

func (c *Compiler) changeOperand2(pos, operand int) {
	converted := make([]byte, 2)
	binary.LittleEndian.PutUint16(converted, uint16(operand))
	c.instructions[pos+1] = op.Code(converted[0])
	c.instructions[pos+2] = op.Code(converted[1])
}

func (c *Compiler) compileInfix(node *ast.Infix) error {
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	switch node.Operator() {
	case "+":
		c.emit(node, op.BinaryOp, int(op.Add))
	case "-":
		c.emit(node, op.BinaryOp, int(op.Subtract))
	case "*":
		c.emit(node, op.BinaryOp, int(op.Multiply))
	case "/":
		c.emit(node, op.BinaryOp, int(op.Divide))
	case "%":
		c.emit(node, op.BinaryOp, int(op.Modulo))
	case "**":
		c.emit(node, op.BinaryOp, int(op.Power))
	case "<<":
		c.emit(node, op.BinaryOp, int(op.LShift))
	case ">>":
		c.emit(node, op.BinaryOp, int(op.RShift))
	case ">":
		c.emit(node, op.CompareOp, int(op.GreaterThan))
	case ">=":
		c.emit(node, op.CompareOp, int(op.GreaterThanOrEqual))
	case "<":
		c.emit(node, op.CompareOp, int(op.LessThan))
	case "<=":
		c.emit(node, op.CompareOp, int(op.LessThanOrEqual))
	case "==":
		c.emit(node, op.CompareOp, int(op.Equal))
	case "!=":
		c.emit(node, op.CompareOp, int(op.NotEqual))
	default:
		return fmt.Errorf("unknown operator: %s", node.Operator())
	}
	return nil
}

func (c *Compiler) constant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) instruction(b []op.Code) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, b...)
	return pos
}

func (c *Compiler) emit(node ast.Node, opcode op.Code, operands ...int) int {
	info := op.GetInfo(opcode)
	fmt.Println("EMIT", opcode, info.Name)
	inst := MakeInstruction(opcode, operands...)
	pos := c.instruction(inst)
	return pos
}

func MakeInstruction(opcode op.Code, operands ...int) []op.Code {
	opInfo := op.OperandCount[opcode]

	totalLen := 1
	for _, w := range opInfo.OperandWidths {
		totalLen += w
	}

	instruction := make([]byte, totalLen)
	instruction[0] = byte(opcode)

	offset := 1
	for i, o := range operands {
		width := opInfo.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			n := uint16(o)
			instruction[offset] = byte(n)
			instruction[offset+1] = byte(n >> 8)
		}
		offset += width
	}

	result := make([]op.Code, 0, len(instruction))
	for _, value := range instruction {
		result = append(result, op.Code(value))
	}
	return result
}

func ReadInstruction(bytes []byte) (op.Code, []int, []byte) {
	opcode := op.Code(bytes[0])
	opInfo := op.OperandCount[opcode]
	totalWidth := 0
	var operands []int
	for i := 0; i < opInfo.OperandCount; i++ {
		width := opInfo.OperandWidths[i]
		totalWidth += width
		switch width {
		case 1:
			operands = append(operands, int(bytes[1]))
		case 2:
			operands = append(operands, int(binary.LittleEndian.Uint16(bytes[1:3])))
		}
	}
	return opcode, operands, bytes[1+totalWidth:]
}
