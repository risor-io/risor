package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

type Compiler struct {
	symbols      *SymbolTable
	constants    []object.Object
	instructions []byte
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

func (c *Compiler) Compile(node ast.Node) error {
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
	case *ast.Infix:
		if err := c.compileInfix(node); err != nil {
			return err
		}
	case *ast.Program:
		for _, stmt := range node.Statements() {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}
	case *ast.Var:
		name, expr := node.Value()
		if err := c.Compile(expr); err != nil {
			return err
		}
		symbol, err := c.symbols.Insert(name, SymbolAttrs{})
		if err != nil {
			return err
		}
		c.emit(node, op.StoreFast, symbol.Index)
	}
	// panic(fmt.Sprintf("unknown ast node type: %T", node))
	return nil
}

func (c *Compiler) compileInfix(node *ast.Infix) error {
	if err := c.Compile(node.Left()); err != nil {
		return err
	}
	if err := c.Compile(node.Right()); err != nil {
		return err
	}
	node.
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
	default:
		return fmt.Errorf("unknown operator: %s", node.Operator())
	}
	return nil
}

func (c *Compiler) constant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) instruction(b []byte) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, b...)
	return pos
}

func (c *Compiler) emit(node ast.Node, opcode op.Code, operands ...int) int {
	inst := MakeInstruction(opcode, operands...)
	pos := c.instruction(inst)
	return pos
}

func MakeInstruction(opcode op.Code, operands ...int) []byte {
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
			instruction[offset] = byte(n >> 8)
			instruction[offset+1] = byte(n)
		}
		offset += width
	}
	return instruction
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
			operands = append(operands, int(binary.BigEndian.Uint16(bytes[1:3])))
		}
	}
	return opcode, operands, bytes[1+totalWidth:]
}
