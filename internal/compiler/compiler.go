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
	Scopes []*Scope
	// Constants []object.Object
	// Symbols   *SymbolTable
}

type Scope struct {
	Name         string
	Instructions []op.Code
	children     []*Scope
	parent       *Scope
	Symbols      *SymbolTable
	Constants    []object.Object
	loops        []*Loop
}

type Compiler struct {
	scopes       []*Scope
	currentScope *Scope
}

type Options struct {
	Builtins []*object.Builtin
}

type Loop struct {
	ContinuePos []int
	BreakPos    []int
}

func New(opts Options) *Compiler {
	symbols := NewSymbolTable()
	for _, b := range opts.Builtins {
		symbols.Insert(b.Name(), SymbolAttrs{
			IsBuiltin: true,
			// Type:      string(b.Type()),
		})
	}
	mainScope := &Scope{
		Name:    "main",
		Symbols: symbols,
	}
	return &Compiler{
		scopes:       []*Scope{mainScope},
		currentScope: mainScope,
	}
}

// func (c *Compiler) Symbols() *SymbolTable {
// 	return c.Symbols
// }

func (c *Compiler) CurrentScope() *Scope {
	return c.currentScope
}

func (c *Compiler) Instructions() []op.Code {
	return c.CurrentScope().Instructions
}

// func (c *Compiler) Constants() []object.Object {
// 	return c.constants
// }

func (c *Compiler) Compile(node ast.Node) (*Bytecode, error) {
	if err := c.compile(node); err != nil {
		return nil, err
	}
	return &Bytecode{Scopes: c.scopes}, nil
}

func (c *Compiler) compile(node ast.Node) error {
	scope := c.CurrentScope()
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
		scope.Symbols = scope.Symbols.NewChild()
		defer func() {
			scope.Symbols = scope.Symbols.Parent()
		}()
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
		symbol, err := scope.Symbols.Insert(name, SymbolAttrs{})
		if err != nil {
			return err
		}
		c.emit(node, op.StoreFast, symbol.Index)
	case *ast.Assign:
		if err := c.compileAssign(node); err != nil {
			return err
		}
	case *ast.Ident:
		name := node.Literal()
		symbol, found := scope.Symbols.Lookup(name)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch symbol.Scope {
		case ScopeGlobal:
			c.emit(node, op.LoadGlobal, symbol.Index)
		case ScopeLocal:
			c.emit(node, op.LoadFast, symbol.Index)
		}
	case *ast.For:
		if err := c.compileFor(node); err != nil {
			return err
		}
	case *ast.Control:
		if err := c.compileControl(node); err != nil {
			return err
		}
	case *ast.Call:
		if err := c.compileCall(node); err != nil {
			return err
		}
	case *ast.Func:
		if err := c.compileFunc(node); err != nil {
			return err
		}
	default:
		fmt.Println("DEFAULT", node, reflect.TypeOf(node))
	}
	// panic(fmt.Sprintf("unknown ast node type: %T", node))
	return nil
}

func (c *Compiler) currentLoop() *Loop {
	scope := c.CurrentScope()
	if len(scope.loops) == 0 {
		return nil
	}
	return scope.loops[len(scope.loops)-1]
}

func (c *Compiler) compileFunc(node *ast.Func) error {

	// scope.Symbols = scope.Symbols.NewChild()
	// defer func() {
	// 	scope.Symbols = scope.Symbols.Parent()
	// }()

	var name string
	ident := node.Name()
	if ident != nil {
		name = ident.Literal()
	} else {
		name = "anonymous"
	}

	funcScope := &Scope{
		Name:    name,
		parent:  c.CurrentScope(),
		Symbols: c.currentScope.Symbols.NewChild(),
	}
	c.currentScope.children = append(c.currentScope.children, funcScope)
	c.scopes = append(c.scopes, funcScope)
	c.currentScope = funcScope
	defer func() {
		c.currentScope = c.currentScope.parent
	}()

	for _, arg := range node.Parameters() {
		funcScope.Symbols.Insert(arg.Literal(), SymbolAttrs{})
	}
	if err := c.compile(node.Body()); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) compileCall(node *ast.Call) error {
	args := node.Arguments()
	if err := c.compile(node.Function()); err != nil {
		return err
	}
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	c.emit(node, op.Call, len(args))
	return nil
}

func (c *Compiler) compileControl(node *ast.Control) error {
	literal := node.Literal()
	if literal == "return" {
		if c.currentScope.parent == nil {
			return fmt.Errorf("return outside of function")
		}
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		c.emit(node, op.ReturnValue)
		return nil
	}
	loop := c.currentLoop()
	if loop == nil {
		if literal == "break" {
			return fmt.Errorf("break outside of loop")
		}
		return fmt.Errorf("continue outside of loop")
	}
	if literal == "break" {
		controlPos := c.emit(node, op.JumpForward, 9999)
		loop.BreakPos = append(loop.BreakPos, controlPos)
	} else {
		controlPos := c.emit(node, op.JumpBackward, 9999)
		loop.ContinuePos = append(loop.ContinuePos, controlPos)
	}
	return nil
}

func (c *Compiler) compileAssign(node *ast.Assign) error {
	name := node.Name()
	symbol, found := c.currentScope.Symbols.Lookup(name)
	if !found {
		return fmt.Errorf("undefined variable: %s", name)
	}
	if node.Operator() == "=" {
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		switch symbol.Scope {
		case ScopeGlobal:
			c.emit(node, op.StoreGlobal, symbol.Index)
		case ScopeLocal:
			c.emit(node, op.StoreFast, symbol.Index)
		}
		return nil
	}
	// Push LHS as TOS
	switch symbol.Scope {
	case ScopeGlobal:
		c.emit(node, op.LoadGlobal, symbol.Index)
	case ScopeLocal:
		c.emit(node, op.LoadFast, symbol.Index)
	}
	// Push RHS as TOS
	if err := c.compile(node.Value()); err != nil {
		return err
	}
	// Result becomes TOS
	switch node.Operator() {
	case "+=":
		c.emit(node, op.BinaryOp, int(op.Add))
	case "-=":
		c.emit(node, op.BinaryOp, int(op.Subtract))
	case "*=":
		c.emit(node, op.BinaryOp, int(op.Multiply))
	case "/=":
		c.emit(node, op.BinaryOp, int(op.Divide))
	}
	// Store TOS in LHS
	switch symbol.Scope {
	case ScopeGlobal:
		c.emit(node, op.StoreGlobal, symbol.Index)
	case ScopeLocal:
		c.emit(node, op.StoreFast, symbol.Index)
	}
	return nil
}

func (c *Compiler) compileFor(node *ast.For) error {
	if node.IsSimpleLoop() {
		return c.compileSimpleFor(node)
	}
	return nil
}

func (c *Compiler) startLoop() *Loop {
	loop := &Loop{}
	c.currentScope.loops = append(c.currentScope.loops, loop)
	return loop
}

func (c *Compiler) endLoop() {
	scope := c.currentScope
	scope.loops = scope.loops[:len(scope.loops)-1]
}

func (c *Compiler) compileSimpleFor(node *ast.For) error {
	scope := c.currentScope
	scope.Symbols = scope.Symbols.NewChild()
	loop := c.startLoop()
	defer func() {
		c.endLoop()
		scope.Symbols = scope.Symbols.Parent()
	}()
	startPos := len(c.Instructions())
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	c.emit(node, op.JumpBackward, c.calculateDelta(startPos))
	nopPos := c.emit(node, op.Nop)
	for _, pos := range loop.BreakPos {
		delta := nopPos - pos
		c.changeOperand2(pos, delta)
	}
	for _, pos := range loop.ContinuePos {
		delta := pos - startPos
		c.changeOperand2(pos, delta)
	}
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
	return len(c.CurrentScope().Instructions) - pos
}

func (c *Compiler) changeOperand2(pos, operand int) {
	instrs := c.CurrentScope().Instructions
	converted := make([]byte, 2)
	binary.LittleEndian.PutUint16(converted, uint16(operand))
	instrs[pos+1] = op.Code(converted[0])
	instrs[pos+2] = op.Code(converted[1])
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
	scope := c.currentScope
	scope.Constants = append(scope.Constants, obj)
	return len(scope.Constants) - 1
}

func (c *Compiler) instruction(b []op.Code) int {
	scope := c.CurrentScope()
	pos := len(scope.Instructions)
	scope.Instructions = append(scope.Instructions, b...)
	return pos
}

func (c *Compiler) emit(node ast.Node, opcode op.Code, operands ...int) int {
	info := op.GetInfo(opcode)
	fmt.Println("EMIT", opcode, info.Name, operands)
	inst := MakeInstruction(opcode, operands...)
	return c.instruction(inst)
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
