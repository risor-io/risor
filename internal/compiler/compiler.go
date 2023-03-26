package compiler

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/internal/symbol"
	"github.com/cloudcmds/tamarin/object"
)

type Scope struct {
	Name         string
	IsNamed      bool
	Parent       *Scope
	Children     []*Scope
	Symbols      *symbol.Table
	Instructions []op.Code
	Constants    []object.Object
	Loops        []*Loop
	Names        []string
}

func (s *Scope) AddName(name string) int {
	s.Names = append(s.Names, name)
	return len(s.Names) - 1
}

type Compiler struct {
	scopes       []*Scope
	currentScope *Scope
}

type Options struct {
	GlobalSymbols *symbol.Table
	Name          string
}

type Loop struct {
	ContinuePos []int
	BreakPos    []int
}

func NewGlobalSymbols() *symbol.Table {
	table := symbol.NewTable()
	for _, b := range evaluator.GlobalBuiltins() {
		table.Insert(b.Name(), symbol.Attrs{Value: b})
	}
	return table
}

func New(opts Options) *Compiler {
	var symbols *symbol.Table
	if opts.GlobalSymbols != nil {
		symbols = opts.GlobalSymbols
	} else {
		symbols = symbol.NewTable()
	}
	mainScope := &Scope{
		Name:    opts.Name,
		Symbols: symbols,
	}
	return &Compiler{
		scopes:       []*Scope{mainScope},
		currentScope: mainScope,
	}
}

func (c *Compiler) CurrentScope() *Scope {
	return c.currentScope
}

func (c *Compiler) Instructions() []op.Code {
	return c.CurrentScope().Instructions
}

func (c *Compiler) Compile(node ast.Node) (*Scope, error) {
	if err := c.compile(node); err != nil {
		return nil, err
	}
	return c.scopes[0], nil
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
		scope.Symbols = scope.Symbols.NewChild(true)
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
		sym, err := scope.Symbols.Insert(name, symbol.Attrs{})
		if err != nil {
			return err
		}
		if c.currentScope.Parent == nil {
			c.emit(node, op.StoreGlobal, sym.Index)
		} else {
			c.emit(node, op.StoreFast, sym.Index)
		}
	case *ast.Assign:
		if err := c.compileAssign(node); err != nil {
			return err
		}
	case *ast.Ident:
		name := node.Literal()
		sym, found := scope.Symbols.Lookup(name)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch sym.Scope {
		case symbol.ScopeGlobal:
			c.emit(node, op.LoadGlobal, sym.Symbol.Index)
		case symbol.ScopeLocal:
			c.emit(node, op.LoadFast, sym.Symbol.Index)
		case symbol.ScopeFree:
			c.emit(node, op.LoadFree, sym.Symbol.Index)
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
	case *ast.List:
		if err := c.compileList(node); err != nil {
			return err
		}
	case *ast.Map:
		if err := c.compileMap(node); err != nil {
			return err
		}
	case *ast.Set:
		if err := c.compileSet(node); err != nil {
			return err
		}
	case *ast.Index:
		if err := c.compileIndex(node); err != nil {
			return err
		}
	case *ast.GetAttr:
		if err := c.compileGetAttr(node); err != nil {
			return err
		}
	case *ast.ObjectCall:
		if err := c.compileObjectCall(node); err != nil {
			return err
		}
	case *ast.Prefix:
		if err := c.compilePrefix(node); err != nil {
			return err
		}
	case *ast.In:
		if err := c.compileIn(node); err != nil {
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
	if len(scope.Loops) == 0 {
		return nil
	}
	return scope.Loops[len(scope.Loops)-1]
}

func (c *Compiler) compileIn(node *ast.In) error {
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	c.emit(node, op.ContainsOp, 0)
	return nil
}

func (c *Compiler) compilePrefix(node *ast.Prefix) error {
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	switch node.Operator() {
	case "!":
		c.emit(node, op.UnaryNot)
	case "-":
		c.emit(node, op.UnaryNegative)
	}
	return nil
}

func (c *Compiler) compileObjectCall(node *ast.ObjectCall) error {
	if err := c.compile(node.Object()); err != nil {
		return err
	}
	expr := node.Call()
	method, ok := expr.(*ast.Call)
	if !ok {
		return fmt.Errorf("invalid call expression")
	}
	name := method.Function().String()
	c.emit(node, op.LoadAttr, c.currentScope.AddName(name))
	args := method.Arguments()
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	c.emit(node, op.Call, len(args))
	return nil
}

func (c *Compiler) compileGetAttr(node *ast.GetAttr) error {
	if err := c.compile(node.Object()); err != nil {
		return err
	}
	idx := c.currentScope.AddName(node.Name())
	c.emit(node, op.LoadAttr, idx)
	return nil
}

func (c *Compiler) compileIndex(node *ast.Index) error {
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	if err := c.compile(node.Index()); err != nil {
		return err
	}
	c.emit(node, op.BinarySubscr)
	return nil
}

func (c *Compiler) compileList(node *ast.List) error {
	for _, expr := range node.Items() {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	c.emit(node, op.BuildList, len(node.Items()))
	return nil
}

func (c *Compiler) compileMap(node *ast.Map) error {
	for k, v := range node.Items() {
		if err := c.compile(k); err != nil {
			return err
		}
		if err := c.compile(v); err != nil {
			return err
		}
	}
	c.emit(node, op.BuildMap, len(node.Items()))
	return nil
}

func (c *Compiler) compileSet(node *ast.Set) error {
	for _, expr := range node.Items() {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	c.emit(node, op.BuildSet, len(node.Items()))
	return nil
}

func (c *Compiler) compileFunc(node *ast.Func) error {

	// Python cell variables:
	// https://stackoverflow.com/questions/23757143/what-is-a-cell-in-the-context-of-an-interpreter-or-compiler

	var name string
	ident := node.Name()
	if ident != nil {
		name = ident.Literal()
	}

	funcScope := &Scope{
		Name:    name,
		IsNamed: ident != nil,
		Parent:  c.CurrentScope(),
		Symbols: c.currentScope.Symbols.NewChild(false),
	}
	c.currentScope.Children = append(c.currentScope.Children, funcScope)
	c.scopes = append(c.scopes, funcScope)
	c.currentScope = funcScope

	paramsIdx := map[string]int{}
	paramsAst := node.Parameters()
	params := make([]string, len(paramsAst))
	for i, param := range paramsAst {
		params[i] = param.Literal()
		paramsIdx[param.Literal()] = i
	}

	defaults := make([]object.Object, len(paramsAst))
	for k := range node.Defaults() {
		idx := paramsIdx[k]
		defaults[idx] = object.NewInt(0) // FIXME
	}

	// Add the function's own name to the symbol table to support recursive calls
	for _, arg := range node.Parameters() {
		funcScope.Symbols.Insert(arg.Literal(), symbol.Attrs{})
	}
	if ident != nil {
		funcScope.Symbols.Insert(name, symbol.Attrs{})
	}
	statements := node.Body().Statements()
	for _, statement := range statements {
		if err := c.compile(statement); err != nil {
			return err
		}
	}
	if len(statements) == 0 {
		c.emit(node, op.Nil)
		c.emit(node, op.ReturnValue, 1)
	} else if _, ok := statements[len(statements)-1].(*ast.Control); !ok {
		c.emit(node, op.ReturnValue, 1)
	}
	c.currentScope = c.currentScope.Parent
	freeSymbols := funcScope.Symbols.Free()
	fn := object.NewCompiledFunction(name, params, defaults, funcScope.Instructions, funcScope)
	if len(freeSymbols) > 0 {
		for _, resolution := range freeSymbols {
			c.emit(nil, op.MakeCell, resolution.Symbol.Index, resolution.Depth-1)
		}
		c.emit(node, op.LoadClosure, c.constant(fn), len(freeSymbols))
	} else {
		c.emit(node, op.LoadConst, c.constant(fn))
	}
	if node.Name() != nil {
		funcSymbol, err := c.currentScope.Symbols.Insert(name, symbol.Attrs{})
		if err != nil {
			return err
		}
		if c.currentScope.Parent == nil {
			c.emit(node, op.StoreGlobal, funcSymbol.Index)
		} else {
			c.emit(node, op.StoreFast, funcSymbol.Index)
		}
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
		if c.currentScope.Parent == nil {
			return fmt.Errorf("return outside of function")
		}
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		c.emit(node, op.ReturnValue, 1)
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
	sym, found := c.currentScope.Symbols.Lookup(name)
	if !found {
		return fmt.Errorf("undefined variable: %s", name)
	}
	if node.Operator() == "=" {
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		switch sym.Scope {
		case symbol.ScopeGlobal:
			c.emit(node, op.StoreGlobal, sym.Symbol.Index)
		case symbol.ScopeLocal:
			c.emit(node, op.StoreFast, sym.Symbol.Index)
		case symbol.ScopeFree:
			c.emit(node, op.StoreFree, sym.Symbol.Index)
		}
		return nil
	}
	// Push LHS as TOS
	switch sym.Scope {
	case symbol.ScopeGlobal:
		c.emit(node, op.LoadGlobal, sym.Symbol.Index)
	case symbol.ScopeLocal:
		c.emit(node, op.LoadFast, sym.Symbol.Index)
	case symbol.ScopeFree:
		c.emit(node, op.LoadFree, sym.Symbol.Index)
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
	switch sym.Scope {
	case symbol.ScopeGlobal:
		c.emit(node, op.StoreGlobal, sym.Symbol.Index)
	case symbol.ScopeLocal:
		c.emit(node, op.StoreFast, sym.Symbol.Index)
	case symbol.ScopeFree:
		c.emit(node, op.StoreFree, sym.Symbol.Index)
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
	c.currentScope.Loops = append(c.currentScope.Loops, loop)
	return loop
}

func (c *Compiler) endLoop() {
	scope := c.currentScope
	scope.Loops = scope.Loops[:len(scope.Loops)-1]
}

func (c *Compiler) compileSimpleFor(node *ast.For) error {
	scope := c.currentScope
	scope.Symbols = scope.Symbols.NewChild(true)
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
	// fmt.Println("constant", obj, len(scope.Constants)-1)
	return len(scope.Constants) - 1
}

func (c *Compiler) instruction(b []op.Code) int {
	scope := c.CurrentScope()
	pos := len(scope.Instructions)
	scope.Instructions = append(scope.Instructions, b...)
	return pos
}

func (c *Compiler) emit(node ast.Node, opcode op.Code, operands ...int) int {
	// info := op.GetInfo(opcode)
	// fmt.Printf("EMIT %2d %-25s %v %p\n", opcode, info.Name, operands, c.currentScope)
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

func ReadOp(instructions []op.Code) (op.Code, []int) {
	opcode := instructions[0]
	opInfo := op.OperandCount[opcode]
	var operands []int
	offset := 0
	for i := 0; i < opInfo.OperandCount; i++ {
		width := opInfo.OperandWidths[i]
		switch width {
		case 1:
			operands = append(operands, int(instructions[offset+1]))
		case 2:
			operands = append(operands, int(binary.LittleEndian.Uint16([]byte{byte(instructions[offset+1]), byte(instructions[offset+2])})))
		}
		offset += width
	}
	return opcode, operands
}
