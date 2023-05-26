package compiler

import (
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/token"
)

const (
	MaxArgs = 255
)

type Compiler struct {

	// The entrypoint code we are compiling. This remains fixed throughout
	// the compilation process.
	main *object.Code

	// The current code we are compiling into. This changes as we enter
	// and leave functions.
	current *object.Code

	// Set on a compilation error
	failure error
}

// Options used to configure compilation
type Options struct {

	// Builtins that should be available during compilation. The name in the
	// map is the name the builtin will be available as in the code.
	Builtins map[string]object.Object

	// Name to be given to the code we are compiling
	Name string

	// Optional code object to append to
	Code *object.Code
}

func New(opts Options) *Compiler {
	// By default we create a new, empty code object to compile into. However
	// if the caller supplied an existing code object we will append to that.
	// This especially supports the REPL where compilation is incremental.
	var main *object.Code
	if opts.Code != nil {
		main = opts.Code
	} else {
		main = &object.Code{Name: opts.Name, Symbols: object.NewSymbolTable()}
	}
	// Insert builtins into the symbol table, including the objects themselves
	for name, builtin := range opts.Builtins {
		if _, err := main.Symbols.InsertBuiltin(name, builtin); err != nil {
			panic(fmt.Sprintf("failed to insert builtin %s: %s", name, err))
		}
	}
	return &Compiler{main: main, current: main}
}

func (c *Compiler) Instructions() []op.Code {
	return c.main.Instructions
}

func (c *Compiler) Compile(node ast.Node) (*object.Code, error) {
	c.failure = nil
	if err := c.compile(node); err != nil {
		return nil, err
	}
	// Check for failures that happened that aren't propagated up the call
	// stack. Some errors are difficult to propagate without bloating the code.
	if c.failure != nil {
		return nil, c.failure
	}
	return c.main, nil
}

func (c *Compiler) compile(node ast.Node) error {
	code := c.current
	switch node := node.(type) {
	case *ast.Nil:
		c.emit(op.Nil)
	case *ast.Int:
		c.emit(op.LoadConst, c.constant(object.NewInt(node.Value())))
	case *ast.Float:
		c.emit(op.LoadConst, c.constant(object.NewFloat(node.Value())))
	case *ast.String:
		if err := c.compileString(node); err != nil {
			return err
		}
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
		statements := node.Statements()
		count := len(statements)
		if count == 0 {
			// Guarantee that the program evaluates to a value
			c.emit(op.Nil)
		} else {
			for i, stmt := range statements {
				if err := c.compile(stmt); err != nil {
					return err
				}
				if i < count-1 {
					if stmt.IsExpression() {
						c.emit(op.PopTop)
					}
				}
			}
			// Guarantee that the program evaluates to a value
			lastStatement := statements[count-1]
			if !lastStatement.IsExpression() {
				c.emit(op.Nil)
			}
		}
	case *ast.Block:
		code.Symbols = code.Symbols.NewBlock()
		defer func() {
			code.Symbols = code.Symbols.Parent()
		}()
		statements := node.Statements()
		count := len(statements)
		if count == 0 {
			// Guarantee that the block evaluates to a value
			c.emit(op.Nil)
		} else {
			for i, stmt := range statements {
				if err := c.compile(stmt); err != nil {
					return err
				}
				if i < count-1 {
					if stmt.IsExpression() {
						c.emit(op.PopTop)
					}
				}
			}
			// Guarantee that the block evaluates to a value
			lastStatement := statements[count-1]
			if !lastStatement.IsExpression() {
				c.emit(op.Nil)
			}
		}
	case *ast.Var:
		name, expr := node.Value()
		if err := c.compile(expr); err != nil {
			return err
		}
		sym, err := code.Symbols.InsertVariable(name)
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
		sym, found := code.Symbols.Lookup(name)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch sym.Scope {
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
	case *ast.Ternary:
		if err := c.compileTernary(node); err != nil {
			return err
		}
	case *ast.Range:
		if err := c.compileRange(node); err != nil {
			return err
		}
	case *ast.Slice:
		if err := c.compileSlice(node); err != nil {
			return err
		}
	case *ast.Import:
		if err := c.compileImport(node); err != nil {
			return err
		}
	case *ast.Switch:
		if err := c.compileSwitch(node); err != nil {
			return err
		}
	case *ast.MultiVar:
		if err := c.compileMultiVar(node); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unknown ast node type: %T", node))
	}
	return nil
}

func (c *Compiler) currentLoop() *object.Loop {
	code := c.current
	if len(code.Loops) == 0 {
		return nil
	}
	return code.Loops[len(code.Loops)-1]
}

func (c *Compiler) currentPosition() int {
	return len(c.Instructions())
}

func (c *Compiler) compileMultiVar(node *ast.MultiVar) error {
	names, expr := node.Value()
	if len(names) > math.MaxUint16 {
		return fmt.Errorf("too many variables in multi-variable assignment")
	}
	// Compile the RHS value
	if err := c.compile(expr); err != nil {
		return err
	}
	// Emit the Unpack opcode to unpack the tuple-like object onto the stack
	c.emit(op.Unpack, uint16(len(names)))
	// Iterate through the names in reverse order and assign the values
	if node.IsWalrus() {
		for i := len(names) - 1; i >= 0; i-- {
			name := names[i]
			sym, err := c.current.Symbols.InsertVariable(name)
			if err != nil {
				return err
			}
			if c.current.Parent == nil {
				c.emit(op.StoreGlobal, sym.Index)
			} else {
				c.emit(op.StoreFast, sym.Index)
			}
		}
		return nil
	}
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		sym, found := c.current.Symbols.Lookup(name)
		if !found {
			return fmt.Errorf("undefined variable: %s", name)
		}
		switch sym.Scope {
		case object.ScopeGlobal:
			c.emit(op.StoreGlobal, sym.Symbol.Index)
		case object.ScopeLocal:
			c.emit(op.StoreFast, sym.Symbol.Index)
		case object.ScopeFree:
			c.emit(op.StoreFree, sym.Symbol.Index)
		case object.ScopeBuiltin:
			c.emit(op.LoadBuiltin, sym.Symbol.Index)
		}
	}
	return nil
}

func (c *Compiler) compileSwitch(node *ast.Switch) error {
	// Compile the switch expression
	if err := c.compile(node.Value()); err != nil {
		return err
	}

	choices := node.Choices()

	// Emit jump positions for each case
	var caseJumpPositions []int
	defaultJumpPos := -1

	for i, choice := range choices {
		if choice.IsDefault() {
			defaultJumpPos = i
			continue
		}
		for _, expr := range choice.Expressions() {
			// Duplicate the switch value for each case comparison
			c.emit(op.Copy, 0)
			// Compile the case expression
			if err := c.compile(expr); err != nil {
				return err
			}
			// Emit the CompareOp equal comparison
			c.emit(op.CompareOp, uint16(op.Equal))
			// Emit PopJumpForwardIfTrue and store its position
			jumpPos := c.emit(op.PopJumpForwardIfTrue, 9999)
			caseJumpPositions = append(caseJumpPositions, jumpPos)
		}
	}

	jumpDefaultPos := c.emit(op.JumpForward, 9999)

	// Update case jump positions and compile case blocks
	var offset int
	var endBlockPosits []int
	for i, choice := range choices {
		if i == defaultJumpPos {
			continue
		}
		for range choice.Expressions() {
			delta, err := c.calculateDelta(caseJumpPositions[offset])
			if err != nil {
				return err
			}
			c.changeOperand(caseJumpPositions[offset], delta)
			offset++
		}
		if choice.Block() == nil {
			// Empty case block
			c.emit(op.Nil)
		} else {
			if err := c.compile(choice.Block()); err != nil {
				return err
			}
		}
		endBlockPosits = append(endBlockPosits, c.emit(op.JumpForward, 9999))
	}

	delta, err := c.calculateDelta(jumpDefaultPos)
	if err != nil {
		return err
	}
	c.changeOperand(jumpDefaultPos, delta)

	// Compile the default case block if it exists
	if defaultJumpPos != -1 {
		if err := c.compile(choices[defaultJumpPos].Block()); err != nil {
			return err
		}
	} else {
		c.emit(op.Nil)
	}

	// Update end block jump positions
	for _, pos := range endBlockPosits {
		delta, err := c.calculateDelta(pos)
		if err != nil {
			return err
		}
		c.changeOperand(pos, delta)
	}

	c.emit(op.Swap, 1)

	// Remove the duplicated switch value from the stack
	c.emit(op.PopTop)
	return nil
}

func (c *Compiler) compileImport(node *ast.Import) error {
	name := node.Module().String()
	c.emit(op.LoadConst, c.constant(object.NewString(name)))
	c.emit(op.Import)
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

func (c *Compiler) compileSlice(node *ast.Slice) error {
	if err := c.compile(node.Left()); err != nil {
		return err
	}
	to := node.ToIndex()
	if to == nil {
		c.emit(op.Copy, 0)
		c.emit(op.Length)
	} else {
		if err := c.compile(to); err != nil {
			return err
		}
	}
	from := node.FromIndex()
	if from == nil {
		c.emit(op.LoadConst, c.constant(object.NewInt(0)))
	} else {
		if err := c.compile(from); err != nil {
			return err
		}
	}
	c.emit(op.Slice)
	return nil
}

func (c *Compiler) compileRange(node *ast.Range) error {
	if err := c.compile(node.Container()); err != nil {
		return err
	}
	c.emit(op.Range)
	return nil
}

func (c *Compiler) compileTernary(node *ast.Ternary) error {
	// evaluate the condition and then conditionally jump to the false case
	if err := c.compile(node.Condition()); err != nil {
		return err
	}
	jumpIfFalsePos := c.emit(op.PopJumpForwardIfFalse, 9999)

	// true case execution, then jump over false case
	if err := c.compile(node.IfTrue()); err != nil {
		return err
	}
	trueCaseEndPos := c.emit(op.JumpForward, 9999)

	// set the jump amount to reach the beginning of the false case
	falseCaseDelta, err := c.calculateDelta(jumpIfFalsePos)
	if err != nil {
		return err
	}
	c.changeOperand(jumpIfFalsePos, falseCaseDelta)

	// false case execution
	if err := c.compile(node.IfFalse()); err != nil {
		return err
	}

	// set the jump amount for the end of the true case
	endDelta, err := c.calculateDelta(trueCaseEndPos)
	if err != nil {
		return err
	}
	c.changeOperand(trueCaseEndPos, endDelta)
	return nil
}

func (c *Compiler) compileString(node *ast.String) error {

	// Is the string a template or a simple string?
	tmpl := node.Template()

	// Simple strings are just emitted as a constant
	if tmpl == nil {
		c.emit(op.LoadConst, c.constant(object.NewString(node.Value())))
		return nil
	}

	if len(tmpl.Fragments) > math.MaxUint16 {
		return fmt.Errorf("string template exceeded max fragment size")
	}

	var expressionIndex int
	expressions := node.TemplateExpressions()

	// Emit code that pushes each fragment of the string onto the stack
	for _, f := range tmpl.Fragments {
		switch f.IsVariable {
		case true:
			expr := expressions[expressionIndex]
			expressionIndex++
			// Nil expression should be treated as empty string
			if expr == nil {
				c.emit(op.LoadConst, c.constant(object.NewString("")))
				continue
			}
			// Transform the expression into a *ast.Func
			astFn := ast.NewFunc(
				token.Token{},
				nil, // no name
				nil, // no params
				nil, // no defaults
				ast.NewBlock(token.Token{}, []ast.Node{expr}),
			)
			// Emit code to push the compiled function as TOS
			if _, err := c.compileReturnFunc(astFn); err != nil {
				return err
			}
			// Emit code to call the function to build the fragment
			c.emit(op.Call, 0)
		case false:
			// Push the fragment as a constant as TOS
			c.emit(op.LoadConst, c.constant(object.NewString(f.Value)))
		}
	}
	// Emit a BuildString to concatenate all the fragments
	c.emit(op.BuildString, uint16(len(tmpl.Fragments)))
	return nil
}

func (c *Compiler) compilePipe(node *ast.Pipe) error {
	if c.current.PipeActive {
		return fmt.Errorf("invalid nested pipe")
	}
	exprs := node.Expressions()
	if len(exprs) < 2 {
		return fmt.Errorf("pipe operator requires at least two expressions")
	}
	// Compile the first expression (filling TOS with the initial pipe value)
	if err := c.compile(exprs[0]); err != nil {
		return err
	}
	// Set the pipe active flag for the remainder of the pipe
	c.current.PipeActive = true
	defer func() {
		c.current.PipeActive = false
	}()
	// Iterate over the remaining expressions. Each should eval to a function.
	// TODO: may need to compile to a partial as well.
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
	// Push the named variable onto the stack
	switch sym.Scope {
	case object.ScopeGlobal:
		c.emit(op.LoadGlobal, sym.Symbol.Index)
	case object.ScopeLocal:
		c.emit(op.LoadFast, sym.Symbol.Index)
	case object.ScopeFree:
		c.emit(op.LoadFree, sym.Symbol.Index)
	case object.ScopeBuiltin:
		return fmt.Errorf("cannot assign to builtin: %s", name)
	}
	// Push the integer amount to the stack (1 or -1)
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
	switch sym.Scope {
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

func (c *Compiler) compileCall(node *ast.Call) error {
	args := node.Arguments()
	argc := len(args)
	if argc > MaxArgs {
		return fmt.Errorf("max arguments limit of %d exceeded (got %d)", MaxArgs, argc)
	}
	if err := c.compile(node.Function()); err != nil {
		return err
	}
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	if c.current.PipeActive {
		c.emit(op.Partial, uint16(argc))
	} else {
		c.emit(op.Call, uint16(argc))
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
	argc := len(args)
	if argc > MaxArgs {
		return fmt.Errorf("max arguments limit of %d exceeded (got %d)", MaxArgs, argc)
	}
	for _, arg := range args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	if c.current.PipeActive {
		c.emit(op.Partial, uint16(len(args)))
	} else {
		c.emit(op.Call, uint16(len(args)))
	}
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
	items := node.Items()
	count := len(items)
	if count > math.MaxUint16 {
		return fmt.Errorf("list literal exceeds max size")
	}
	for _, expr := range items {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	c.emit(op.BuildList, uint16(count))
	return nil
}

func (c *Compiler) compileMap(node *ast.Map) error {
	items := node.Items()
	count := len(items)
	for k, v := range items {
		switch k := k.(type) {
		case *ast.String:
			if err := c.compile(k); err != nil {
				return err
			}
		case *ast.Ident:
			c.emit(op.LoadConst, c.constant(object.NewString(k.String())))
		default:
			return fmt.Errorf("invalid map key type: %v", k)
		}
		if err := c.compile(v); err != nil {
			return err
		}
	}
	c.emit(op.BuildMap, uint16(count))
	return nil
}

func (c *Compiler) compileSet(node *ast.Set) error {
	items := node.Items()
	count := len(items)
	for _, expr := range items {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	c.emit(op.BuildSet, uint16(count))
	return nil
}

func (c *Compiler) compileFunc(node *ast.Func) error {
	_, err := c.compileReturnFunc(node)
	return err
}

func (c *Compiler) compileReturnFunc(node *ast.Func) (*object.Function, error) {

	// Python cell variables:
	// https://stackoverflow.com/questions/23757143/what-is-a-cell-in-the-context-of-an-interpreter-or-compiler

	if len(node.Parameters()) > 255 {
		return nil, fmt.Errorf("function exceeded parameter limit of 255")
	}

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
		Parent:  c.current,
		Symbols: c.current.Symbols.NewChild(),
		Source:  node.Body().String(),
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
			return nil, fmt.Errorf("unsupported default value: %s", expr)
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

	// Compile the function body
	body := node.Body()
	if err := c.compile(body); err != nil {
		return nil, err
	}
	if !body.EndsWithReturn() {
		c.emit(op.ReturnValue)
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
	// code. Otherwise, we just leave it on the stack.
	if functionName != "" {
		funcSymbol, err := c.current.Symbols.InsertVariable(functionName)
		if err != nil {
			return nil, err
		}
		if c.current.Parent == nil {
			c.emit(op.StoreGlobal, funcSymbol.Index)
		} else {
			c.emit(op.StoreFast, funcSymbol.Index)
		}
	}
	return fn, nil
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
		c.emit(op.ReturnValue)
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
		loop.BreakPos = append(loop.BreakPos, position)
	} else {
		position := c.emit(op.JumpBackward, 9999)
		loop.ContinuePos = append(loop.ContinuePos, position)
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
		switch sym.Scope {
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
	switch sym.Scope {
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
	switch sym.Scope {
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

func (c *Compiler) compileForRange(forNode *ast.For, name string, rangeNode *ast.Range) error {

	if err := c.compile(rangeNode.Container()); err != nil {
		return err
	}
	c.emit(op.GetIter)

	code := c.current
	code.Symbols = code.Symbols.NewBlock()
	c.startLoop()
	defer func() {
		c.endLoop()
		code.Symbols = code.Symbols.Parent()
	}()

	iterPos := c.emit(op.ForIter, 0)

	// assign the current value of the iterator to the loop variable
	sym, err := code.Symbols.InsertVariable(name)
	if err != nil {
		return err
	}
	if code.Symbols.IsGlobal() {
		c.emit(op.StoreGlobal, sym.Index)
	} else {
		c.emit(op.StoreFast, sym.Index)
	}

	// compile the body of the loop
	if err := c.compile(forNode.Consequence()); err != nil {
		return err
	}
	c.emit(op.PopTop)

	// jump back to the start of the loop
	delta, err := c.calculateDelta(iterPos)
	if err != nil {
		return err
	}
	c.emit(op.JumpBackward, delta)

	delta, err = c.calculateDelta(iterPos)
	if err != nil {
		return err
	}
	c.changeOperand(iterPos, delta)

	return nil
}

func (c *Compiler) compileFor(node *ast.For) error {
	if node.IsSimpleLoop() {
		return c.compileSimpleFor(node)
	}
	if node.Init() == nil && node.Post() == nil {
		varExpr, ok := node.Condition().(*ast.Var)
		if ok {
			name, rhs := varExpr.Value()
			if rangeNode, ok := rhs.(*ast.Range); ok {
				return c.compileForRange(node, name, rangeNode)
			}
		}
	}

	code := c.current
	code.Symbols = code.Symbols.NewBlock()
	loop := c.startLoop()
	defer func() {
		c.endLoop()
		code.Symbols = code.Symbols.Parent()
	}()

	// Compile the init statement if present
	if node.Init() != nil {
		if err := c.compile(node.Init()); err != nil {
			return err
		}
	}

	// Mark the position of the loop start
	loopStart := c.currentPosition()

	// Compile the condition expression if present
	var conditionJumpPos int
	if node.Condition() != nil {
		if err := c.compile(node.Condition()); err != nil {
			return err
		}
		// Emit a jump to execute if the condition fails
		conditionJumpPos = c.emit(op.PopJumpForwardIfFalse, 9999)
	}

	// Compile the loop body
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	c.emit(op.PopTop)

	// Compile the post statement if present
	if node.Post() != nil {
		post := node.Post()
		if err := c.compile(post); err != nil {
			return err
		}
		// If the post statement is an expression, pop the value so its ignored
		if post.IsExpression() {
			c.emit(op.PopTop)
		}
	}

	// Jump back to the loop start
	delta, err := c.calculateDelta(loopStart)
	if err != nil {
		return err
	}
	c.emit(op.JumpBackward, delta)

	// Update the condition jump position
	if conditionJumpPos != 0 {
		delta, err = c.calculateDelta(conditionJumpPos)
		if err != nil {
			return err
		}
		c.changeOperand(conditionJumpPos, delta)
	}

	// Update breaks to jump to this point
	for _, pos := range loop.BreakPos {
		delta, err = c.calculateDelta(conditionJumpPos)
		if err != nil {
			return err
		}
		c.changeOperand(pos, uint16(delta))
	}

	// Update continues to jump to the loop start
	for _, pos := range loop.ContinuePos {
		delta := pos - loopStart
		if delta > math.MaxUint16 {
			return fmt.Errorf("loop size exceeded limits")
		}
		c.changeOperand(pos, uint16(delta))
	}

	return nil
}

func (c *Compiler) startLoop() *object.Loop {
	loop := &object.Loop{}
	c.current.Loops = append(c.current.Loops, loop)
	return loop
}

func (c *Compiler) endLoop() {
	code := c.current
	code.Loops = code.Loops[:len(code.Loops)-1]
}

func (c *Compiler) compileSimpleFor(node *ast.For) error {
	code := c.current
	code.Symbols = code.Symbols.NewBlock()
	loop := c.startLoop()
	defer func() {
		c.endLoop()
		code.Symbols = code.Symbols.Parent()
	}()
	startPos := len(c.Instructions())
	if err := c.compile(node.Consequence()); err != nil {
		return err
	}
	delta, err := c.calculateDelta(startPos)
	if err != nil {
		return err
	}
	c.emit(op.JumpBackward, delta)
	nopPos := c.emit(op.Nop)
	for _, pos := range loop.BreakPos {
		delta := nopPos - pos
		if delta > math.MaxUint16 {
			return fmt.Errorf("loop size exceeded limits")
		}
		c.changeOperand(pos, uint16(delta))
	}
	for _, pos := range loop.ContinuePos {
		delta := pos - startPos
		if delta > math.MaxUint16 {
			return fmt.Errorf("loop size exceeded limits")
		}
		c.changeOperand(pos, uint16(delta))
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
		delta, err := c.calculateDelta(jumpIfFalsePos)
		if err != nil {
			return err
		}
		c.changeOperand(jumpIfFalsePos, delta)
		if err := c.compile(alternative); err != nil {
			return err
		}
		delta, err = c.calculateDelta(jumpForwardPos)
		if err != nil {
			return err
		}
		c.changeOperand(jumpForwardPos, delta)
	} else {
		// Jump forward to skip the alternative by default
		jumpForwardPos := c.emit(op.JumpForward, 9999)
		// Update PopJumpForwardIfFalse to point to this alternative,
		// so that the alternative is executed if the condition is false
		delta, err := c.calculateDelta(jumpIfFalsePos)
		if err != nil {
			return err
		}
		c.changeOperand(jumpIfFalsePos, delta)
		// This allows ifs to be used as expressions. If the if check fails and
		// there is no alternative, the result of the if expression is nil.
		c.emit(op.Nil)
		delta, err = c.calculateDelta(jumpForwardPos)
		if err != nil {
			return err
		}
		c.changeOperand(jumpForwardPos, delta)
	}
	return nil
}

func (c *Compiler) calculateDelta(pos int) (uint16, error) {
	instrCount := len(c.current.Instructions)
	delta := instrCount - pos
	if delta > math.MaxUint16 {
		return 0, fmt.Errorf("jump destination too far away")
	}
	return uint16(delta), nil
}

func (c *Compiler) changeOperand(instructionIndex int, operand uint16) {
	c.current.Instructions[instructionIndex+1] = op.Code(operand)
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
	code := c.current
	if len(code.Constants) >= math.MaxUint16 {
		c.failure = fmt.Errorf("number of constants exceeded limits")
		return 0
	}
	code.Constants = append(code.Constants, obj)
	return uint16(len(code.Constants) - 1)
}

func (c *Compiler) emit(opcode op.Code, operands ...uint16) int {
	inst := MakeInstruction(opcode, operands...)
	code := c.current
	pos := len(code.Instructions)
	code.Instructions = append(code.Instructions, inst...)
	return pos
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
