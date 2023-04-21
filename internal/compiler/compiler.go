package compiler

import (
	"fmt"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
)

type Compiler struct {
	main     *object.Code
	current  *object.Code
	startPos int
}

type Options struct {
	Builtins map[string]object.Object
	Name     string
	Code     *object.Code
}

func New(opts Options) *Compiler {
	var main *object.Code
	if opts.Code != nil {
		main = opts.Code
	} else {
		main = &object.Code{Name: opts.Name, Symbols: object.NewSymbolTable()}
	}
	for name, builtin := range opts.Builtins {
		if _, err := main.Symbols.InsertBuiltin(name, builtin); err != nil {
			panic(fmt.Sprintf("failed to insert builtin %s: %s", name, err))
		}
	}
	return &Compiler{main: main, current: main, startPos: len(main.Instructions)}
}

func (c *Compiler) CurrentScope() *object.Code {
	return c.current
}

func (c *Compiler) Instructions() []op.Code {
	return c.main.Instructions
}

func (c *Compiler) NewInstructions() []op.Code {
	if c.startPos > len(c.main.Instructions) {
		return nil
	}
	return c.main.Instructions[c.startPos:]
}

func (c *Compiler) Compile(node ast.Node) (*object.Code, error) {
	if err := c.compile(node); err != nil {
		return nil, err
	}
	return c.main, nil
}

func (c *Compiler) compile(node ast.Node) error {
	scope := c.CurrentScope()
	switch node := node.(type) {
	case *ast.Nil:
		c.emit(op.Nil)
	case *ast.Int:
		c.emit(op.LoadConst, c.constant(object.NewInt(node.Value())))
	case *ast.Float:
		c.emit(op.LoadConst, c.constant(object.NewFloat(node.Value())))
	case *ast.String:
		c.emit(op.LoadConst, c.constant(object.NewString(node.Value())))
	case *ast.Bool:
		if node.Value() {
			c.emit(op.True)
		} else {
			c.emit(op.False)
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
		scope.Symbols = scope.Symbols.NewBlock()
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
		sym, err := scope.Symbols.InsertVariable(name)
		if err != nil {
			return err
		}
		if c.current.Parent == nil {
			c.emit(op.StoreGlobal, sym.Index)
		} else {
			c.emit(op.StoreFast, sym.Index)
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
		switch sym.Code {
		case object.ScopeGlobal:
			c.emit(op.LoadGlobal, sym.Symbol.Index)
		case object.ScopeLocal:
			c.emit(op.LoadFast, sym.Symbol.Index)
		case object.ScopeFree:
			c.emit(op.LoadFree, sym.Symbol.Index)
		case object.ScopeBuiltin:
			c.emit(op.LoadBuiltin, sym.Symbol.Index)
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
	case *ast.Const:
		if err := c.compileConst(node); err != nil {
			return err
		}
	case *ast.Postfix:
		if err := c.compilePostfix(node); err != nil {
			return err
		}
	case *ast.Pipe:
		if err := c.compilePipe(node); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unknown ast node type: %T", node))
	}
	return nil
}

func (c *Compiler) currentLoop() *object.Loop {
	scope := c.CurrentScope()
	if len(scope.Loops) == 0 {
		return nil
	}
	return scope.Loops[len(scope.Loops)-1]
}

func (c *Compiler) compilePipe(node *ast.Pipe) error {
	exprs := node.Expressions()
	if len(exprs) < 2 {
		return fmt.Errorf("pipe operator requires at least two expressions")
	}
	// Compile the first expression (filling TOS with the initial pipe value)
	if err := c.compile(exprs[0]); err != nil {
		return err
	}
	// Iterate over the remaining expressions. Each should eval to a function.
	for i := 1; i < len(exprs); i++ {
		// Compile the current expression, pushing a function as TOS
		if err := c.compile(exprs[i]); err != nil {
			return err
		}
		// Swap the function (TOS) with the argument below it on the stack
		// and then call the function with one argument
		c.emit(op.Swap, 1)
		c.emit(op.Call, 1)
	}
	return nil
}

func (c *Compiler) compilePostfix(node *ast.Postfix) error {
	name := node.Literal()
	sym, found := c.current.Symbols.Lookup(name)
	if !found {
		return fmt.Errorf("undefined variable: %s", name)
	}
	// Push variable as TOS
	switch sym.Code {
	case object.ScopeGlobal:
		c.emit(op.LoadGlobal, sym.Symbol.Index)
	case object.ScopeLocal:
		c.emit(op.LoadFast, sym.Symbol.Index)
	case object.ScopeFree:
		c.emit(op.LoadFree, sym.Symbol.Index)
	case object.ScopeBuiltin:
		return fmt.Errorf("invalid operation on builtin: %s", name)
	}
	// Push integer 1 or -1 as TOS
	operator := node.Operator()
	if operator == "++" {
		c.emit(op.LoadConst, c.constant(object.NewInt(1)))
	} else if operator == "--" {
		c.emit(op.LoadConst, c.constant(object.NewInt(-1)))
	} else {
		return fmt.Errorf("unknown operator: %q", operator)
	}
	// Run increment or decrement as an Add BinaryOp
	c.emit(op.BinaryOp, uint16(op.Add))
	// Store TOS in LHS
	switch sym.Code {
	case object.ScopeGlobal:
		c.emit(op.StoreGlobal, sym.Symbol.Index)
	case object.ScopeLocal:
		c.emit(op.StoreFast, sym.Symbol.Index)
	case object.ScopeFree:
		c.emit(op.StoreFree, sym.Symbol.Index)
	}
	return nil
}

func (c *Compiler) compileConst(node *ast.Const) error {
	name, expr := node.Value()
	if err := c.compile(expr); err != nil {
		return err
	}
	sym, err := c.current.Symbols.InsertVariable(name)
	if err != nil {
		return err
	}
	if c.current.Parent == nil {
		c.emit(op.StoreGlobal, sym.Index)
	} else {
		c.emit(op.StoreFast, sym.Index)
	}
	return nil
}

func (c *Compiler) compileIn(node *ast.In) error {
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	c.emit(op.ContainsOp, 0)
	return nil
}

func (c *Compiler) compilePrefix(node *ast.Prefix) error {
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	switch node.Operator() {
	case "!":
		c.emit(op.UnaryNot)
	case "-":
		c.emit(op.UnaryNegative)
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
	c.emit(op.LoadAttr, c.current.AddName(name))
	args := method.Arguments()
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	c.emit(op.Call, uint16(len(args)))
	return nil
}

func (c *Compiler) compileGetAttr(node *ast.GetAttr) error {
	if err := c.compile(node.Object()); err != nil {
		return err
	}
	idx := c.current.AddName(node.Name())
	c.emit(op.LoadAttr, idx)
	return nil
}

func (c *Compiler) compileIndex(node *ast.Index) error {
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	if err := c.compile(node.Index()); err != nil {
		return err
	}
	c.emit(op.BinarySubscr)
	return nil
}

func (c *Compiler) compileList(node *ast.List) error {
	for _, expr := range node.Items() {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	// TODO: error on too many items
	c.emit(op.BuildList, uint16(len(node.Items())))
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
	// TODO: error on too many items
	c.emit(op.BuildMap, uint16(len(node.Items())))
	return nil
}

func (c *Compiler) compileSet(node *ast.Set) error {
	for _, expr := range node.Items() {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	// TODO: error on too many items
	c.emit(op.BuildSet, uint16(len(node.Items())))
	return nil
}

func (c *Compiler) compileFunc(node *ast.Func) error {

	// Python cell variables:
	// https://stackoverflow.com/questions/23757143/what-is-a-cell-in-the-context-of-an-interpreter-or-compiler

	// The function has an optional name. If it is named, the name will be
	// stored in the function's own symbol table to support recursive calls.
	var functionName string
	if ident := node.Name(); ident != nil {
		functionName = ident.Literal()
	}

	// This new code object will store the compiled code for this function
	code := &object.Code{
		Name:    functionName,
		IsNamed: functionName != "",
		Parent:  c.CurrentScope(),
		Symbols: c.current.Symbols.NewChild(),
	}

	// Setting current here means subsequent calls to compile will add to this
	// code object instead of the parent.
	c.current = code

	// Make it quick to look up the index of a parameter
	paramsIdx := map[string]int{}
	params := node.ParameterNames()
	for i, name := range params {
		paramsIdx[name] = i
	}

	// Build an array of default values for parameters, supporting only
	// the basic types of int, string, bool, float, and nil.
	defaults := make([]object.Object, len(params))
	for name, expr := range node.Defaults() {
		var value object.Object
		switch expr := expr.(type) {
		case *ast.Int:
			value = object.NewInt(expr.Value())
		case *ast.String:
			value = object.NewString(expr.Value())
		case *ast.Bool:
			value = object.NewBool(expr.Value())
		case *ast.Float:
			value = object.NewFloat(expr.Value())
		case *ast.Nil:
			value = object.Nil
		default:
			return fmt.Errorf("unsupported default value: %s", expr)
		}
		defaults[paramsIdx[name]] = value
	}

	// After the function's name, we'll add the parameter names to the symbols.
	for _, arg := range node.Parameters() {
		code.Symbols.InsertVariable(arg.Literal())
	}
	// Add the function's own name to its symbol table. This supports recursive
	// calls to the function. Later when we create the function object, we'll
	// add the object value to the table.
	if code.IsNamed {
		code.Symbols.InsertVariable(functionName)
	}

	// Compile the function code
	statements := node.Body().Statements()
	for _, statement := range statements {
		if err := c.compile(statement); err != nil {
			return err
		}
	}
	if len(statements) == 0 {
		c.emit(op.Nil)
		c.emit(op.ReturnValue, 1)
	} else if _, ok := statements[len(statements)-1].(*ast.Control); !ok {
		c.emit(op.ReturnValue, 1)
	}

	// We're done compiling the function, so switch back to compiling the parent
	c.current = c.current.Parent

	// Create the function object that contains the compiled code
	fn := object.NewFunction(object.FunctionOpts{
		Name:           functionName,
		ParameterNames: params,
		Defaults:       defaults,
		Code:           code,
	})
	if code.IsNamed {
		code.Symbols.SetValue(functionName, fn)
	}

	// Emit the code to load the function object onto the stack. If there are
	// free variables, we use LoadClosure, otherwise we use LoadConst.
	freeSymbols := code.Symbols.Free()
	if len(freeSymbols) > 0 {
		for _, resolution := range freeSymbols {
			c.emit(op.MakeCell, resolution.Symbol.Index, uint16(resolution.Depth-1))
		}
		c.emit(op.LoadClosure, c.constant(fn), uint16(len(freeSymbols)))
	} else {
		c.emit(op.LoadConst, c.constant(fn))
	}

	// If the function was named, we store it as a named variable in the current
	// scope. Otherwise, we just leave it on the stack.
	if functionName != "" {
		funcSymbol, err := c.current.Symbols.InsertVariable(functionName)
		if err != nil {
			return err
		}
		if c.current.Parent == nil {
			c.emit(op.StoreGlobal, funcSymbol.Index)
		} else {
			c.emit(op.StoreFast, funcSymbol.Index)
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
	c.emit(op.Call, uint16(len(args)))
	return nil
}

func (c *Compiler) compileControl(node *ast.Control) error {
	literal := node.Literal()
	if literal == "return" {
		if c.current.Parent == nil {
			return fmt.Errorf("return outside of function")
		}
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		c.emit(op.ReturnValue, 1)
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
		position := c.emit(op.JumpForward, 9999)
		loop.BreakPos = append(loop.BreakPos, uint16(position))
	} else {
		position := c.emit(op.JumpBackward, 9999)
		loop.ContinuePos = append(loop.ContinuePos, uint16(position))
	}
	return nil
}

func (c *Compiler) compileAssign(node *ast.Assign) error {
	name := node.Name()
	sym, found := c.current.Symbols.Lookup(name)
	if !found {
		return fmt.Errorf("undefined variable: %s", name)
	}
	if node.Operator() == "=" {
		if err := c.compile(node.Value()); err != nil {
			return err
		}
		switch sym.Code {
		case object.ScopeGlobal:
			c.emit(op.StoreGlobal, sym.Symbol.Index)
		case object.ScopeLocal:
			c.emit(op.StoreFast, sym.Symbol.Index)
		case object.ScopeFree:
			c.emit(op.StoreFree, sym.Symbol.Index)
		case object.ScopeBuiltin:
			c.emit(op.LoadBuiltin, sym.Symbol.Index)
		}
		return nil
	}
	// Push LHS as TOS
	switch sym.Code {
	case object.ScopeGlobal:
		c.emit(op.LoadGlobal, sym.Symbol.Index)
	case object.ScopeLocal:
		c.emit(op.LoadFast, sym.Symbol.Index)
	case object.ScopeFree:
		c.emit(op.LoadFree, sym.Symbol.Index)
	case object.ScopeBuiltin:
		c.emit(op.LoadBuiltin, sym.Symbol.Index)
	}
	// Push RHS as TOS
	if err := c.compile(node.Value()); err != nil {
		return err
	}
	// Result becomes TOS
	switch node.Operator() {
	case "+=":
		c.emit(op.BinaryOp, uint16(op.Add))
	case "-=":
		c.emit(op.BinaryOp, uint16(op.Subtract))
	case "*=":
		c.emit(op.BinaryOp, uint16(op.Multiply))
	case "/=":
		c.emit(op.BinaryOp, uint16(op.Divide))
	}
	// Store TOS in LHS
	switch sym.Code {
	case object.ScopeGlobal:
		c.emit(op.StoreGlobal, sym.Symbol.Index)
	case object.ScopeLocal:
		c.emit(op.StoreFast, sym.Symbol.Index)
	case object.ScopeFree:
		c.emit(op.StoreFree, sym.Symbol.Index)
	case object.ScopeBuiltin:
		c.emit(op.LoadBuiltin, sym.Symbol.Index)
	}
	return nil
}

func (c *Compiler) compileFor(node *ast.For) error {
	if node.IsSimpleLoop() {
		return c.compileSimpleFor(node)
	}
	return nil
}

func (c *Compiler) startLoop() *object.Loop {
	loop := &object.Loop{}
	c.current.Loops = append(c.current.Loops, loop)
	return loop
}

func (c *Compiler) endLoop() {
	scope := c.current
	scope.Loops = scope.Loops[:len(scope.Loops)-1]
}

func (c *Compiler) compileSimpleFor(node *ast.For) error {
	scope := c.current
	scope.Symbols = scope.Symbols.NewBlock()
	loop := c.startLoop()
	defer func() {
		c.endLoop()
		scope.Symbols = scope.Symbols.Parent()
	}()
	startPos := uint16(len(c.Instructions()))
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	c.emit(op.JumpBackward, c.calculateDelta(startPos))
	nopPos := c.emit(op.Nop)
	for _, pos := range loop.BreakPos {
		delta := uint16(nopPos) - pos
		c.changeOperand(pos, delta)
	}
	for _, pos := range loop.ContinuePos {
		delta := pos - uint16(startPos)
		c.changeOperand(pos, delta)
	}
	return nil
}

func (c *Compiler) compileIf(node *ast.If) error {
	if err := c.compile(node.Condition()); err != nil {
		return err
	}
	jumpIfFalsePos := c.emit(op.PopJumpForwardIfFalse, 9999)
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	alternative := node.Alternative()
	if alternative != nil {
		// Jump forward to skip the alternative by default
		jumpForwardPos := c.emit(op.JumpForward, 9999)
		// Update PopJumpForwardIfFalse to point to this alternative,
		// so that the alternative is executed if the condition is false
		delta := c.calculateDelta(jumpIfFalsePos)
		c.changeOperand(jumpIfFalsePos, delta)
		if err := c.compile(alternative); err != nil {
			return err
		}
		c.changeOperand(jumpForwardPos, c.calculateDelta(jumpForwardPos))
	} else {
		delta := c.calculateDelta(jumpIfFalsePos)
		c.changeOperand(jumpIfFalsePos, delta)
	}
	return nil
}

func (c *Compiler) calculateDelta(pos uint16) uint16 {
	// TODO: error on overflow
	return uint16(len(c.CurrentScope().Instructions)) - pos
}

func (c *Compiler) changeOperand(pos, operand uint16) {
	c.CurrentScope().Instructions[pos+1] = op.Code(operand)
}

func (c *Compiler) compileInfix(node *ast.Infix) error {
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	if err := c.compile(node.Right()); err != nil {
		return err
	}
	switch node.Operator() {
	case "&&":
		c.emit(op.BinaryOp, uint16(op.And))
	case "||":
		c.emit(op.BinaryOp, uint16(op.Or))
	case "+":
		c.emit(op.BinaryOp, uint16(op.Add))
	case "-":
		c.emit(op.BinaryOp, uint16(op.Subtract))
	case "*":
		c.emit(op.BinaryOp, uint16(op.Multiply))
	case "/":
		c.emit(op.BinaryOp, uint16(op.Divide))
	case "%":
		c.emit(op.BinaryOp, uint16(op.Modulo))
	case "**":
		c.emit(op.BinaryOp, uint16(op.Power))
	case "<<":
		c.emit(op.BinaryOp, uint16(op.LShift))
	case ">>":
		c.emit(op.BinaryOp, uint16(op.RShift))
	case ">":
		c.emit(op.CompareOp, uint16(op.GreaterThan))
	case ">=":
		c.emit(op.CompareOp, uint16(op.GreaterThanOrEqual))
	case "<":
		c.emit(op.CompareOp, uint16(op.LessThan))
	case "<=":
		c.emit(op.CompareOp, uint16(op.LessThanOrEqual))
	case "==":
		c.emit(op.CompareOp, uint16(op.Equal))
	case "!=":
		c.emit(op.CompareOp, uint16(op.NotEqual))
	default:
		return fmt.Errorf("unknown operator: %s", node.Operator())
	}
	return nil
}

func (c *Compiler) constant(obj object.Object) uint16 {
	scope := c.current
	scope.Constants = append(scope.Constants, obj)
	// TODO: error if > 65535
	return uint16(len(scope.Constants) - 1)
}

func (c *Compiler) emit(opcode op.Code, operands ...uint16) uint16 {
	// info := op.GetInfo(opcode)
	// fmt.Printf("EMIT %2d %-25s %v %p\n", opcode, info.Name, operands, c.current)
	inst := MakeInstruction(opcode, operands...)
	// return c.instruction(inst)
	scope := c.CurrentScope()
	pos := len(scope.Instructions)
	scope.Instructions = append(scope.Instructions, inst...)
	return uint16(pos)
}

func MakeInstruction(opcode op.Code, operands ...uint16) []op.Code {
	opInfo := op.OperandCount[opcode]
	if len(operands) != opInfo.OperandCount {
		panic("wrong operand count")
	}
	instruction := make([]op.Code, 1+opInfo.OperandCount)
	instruction[0] = opcode
	offset := 1
	for _, o := range operands {
		instruction[offset] = op.Code(o)
		offset++
	}
	return instruction
}

// func ReadInstruction(data []uint16) (op.Code, []int, []byte) {
// 	opcode := op.Code(bytes[0])
// 	opInfo := op.OperandCount[opcode]
// 	totalWidth := 0
// 	var operands []int
// 	for i := 0; i < opInfo.OperandCount; i++ {
// 		width := opInfo.OperandWidths[i]
// 		totalWidth += width
// 		switch width {
// 		case 1:
// 			operands = append(operands, int(bytes[1]))
// 		case 2:
// 			operands = append(operands, int(binary.LittleEndian.Uint16(bytes[1:3])))
// 		}
// 	}
// 	return opcode, operands, bytes[1+totalWidth:]
// }

// func ReadOp(instructions []op.Code) (op.Code, []int) {
// 	opcode := instructions[0]
// 	opInfo := op.OperandCount[opcode]
// 	var operands []int
// 	offset := 0
// 	for i := 0; i < opInfo.OperandCount; i++ {
// 		width := opInfo.OperandWidths[i]
// 		switch width {
// 		case 1:
// 			operands = append(operands, int(instructions[offset+1]))
// 		case 2:
// 			operands = append(operands, int(binary.LittleEndian.Uint16([]byte{byte(instructions[offset+1]), byte(instructions[offset+2])})))
// 		}
// 		offset += width
// 	}
// 	return opcode, operands
// }
