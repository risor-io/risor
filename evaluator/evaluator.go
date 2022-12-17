// Package evaluator contains the core of our interpreter, which walks
// the AST produced by the parser and evaluates user code.
package evaluator

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// Opts configures Tamarin code evaluation.
type Opts struct {
	// Importer is used to import Tamarin code modules. If nil, module imports
	// are not supported and an import will result in an error that stops code
	// execution.
	Importer Importer

	// If set to true, the default builtins will not be registered.
	DisableDefaultBuiltins bool

	// Supplies extra and/or override builtins for evaluation.
	Builtins []*object.Builtin
}

// Evaluator is used to execute Tamarin AST nodes.
type Evaluator struct {
	importer Importer
	builtins map[string]*object.Builtin
}

// New returns a new Evaluator
func New(opts Opts) *Evaluator {
	e := &Evaluator{
		importer: opts.Importer,
		builtins: map[string]*object.Builtin{},
	}
	// Conditionally register default global builtins
	if !opts.DisableDefaultBuiltins {
		for _, b := range GlobalBuiltins() {
			e.builtins[b.Key()] = b
		}
	}
	// Add override builtins
	for _, b := range opts.Builtins {
		e.builtins[b.Key()] = b
	}
	return e
}

// Returns a function that implements object.CallFunc
func (e *Evaluator) getCallFunc() object.CallFunc {
	return func(ctx context.Context, s interface{}, fn object.Object, args []object.Object) object.Object {
		var scopeObj *scope.Scope
		if s != nil {
			scopeObj = s.(*scope.Scope)
		}
		return e.applyFunction(ctx, scopeObj, fn, args)
	}
}

// Evaluate an AST node. The context can be used to cancel a running evaluation.
// If evaluation encounters an error, a Tamarin error object is returned.
func (e *Evaluator) Evaluate(ctx context.Context, node ast.Node, s *scope.Scope) object.Object {

	// Add an object.CallFunc to the context so that objects can call Tamarin
	// functions if needed
	ctx = object.WithCallFunc(ctx, e.getCallFunc())

	// Check for context timeout
	select {
	case <-ctx.Done():
		return &object.Error{Message: ctx.Err().Error()}
	default:
	}

	// Evaluate the AST node based on its type
	switch node := node.(type) {

	// High level types
	case *ast.Program:
		return e.evalProgram(ctx, node, s)
	case *ast.ExpressionStatement:
		return e.Evaluate(ctx, node.Expression, s)
	case *ast.BlockStatement:
		return e.evalBlockStatement(ctx, node, s)

	// Operator expressions
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(ctx, node, s)
	case *ast.PostfixExpression:
		return e.evalPostfixExpression(s, node.Operator, node)
	case *ast.InfixExpression:
		return e.evalInfixExpression(ctx, node, s)
	case *ast.TernaryExpression:
		return e.evalTernaryExpression(ctx, node, s)

	// Miscellaneous
	case *ast.Identifier:
		return e.evalIdentifier(node, s)
	case *ast.IndexExpression:
		return e.evalIndexExpression(ctx, node, s)
	case *ast.Bool:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ImportStatement:
		return e.evalImportStatement(ctx, node, s)

	// Assignment
	case *ast.VarStatement:
		return e.evalVarStatement(ctx, node, s)
	case *ast.ConstStatement:
		return e.evalConstStatement(ctx, node, s)
	case *ast.AssignStatement:
		return e.evalAssignStatement(ctx, node, s)

	// Functions
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(ctx, node, s)
	case *ast.FunctionDefineLiteral:
		return e.evalFunctionDefinition(ctx, node, s)

	// Calls
	case *ast.ObjectCallExpression:
		return e.evalObjectCallExpression(ctx, node, s)
	case *ast.CallExpression:
		return e.evalCallExpression(ctx, node, s)
	case *ast.GetAttributeExpression:
		return e.evalGetAttributeExpression(ctx, node, s)

	// Control
	case *ast.IfExpression:
		return e.evalIfExpression(ctx, node, s)
	case *ast.ForLoopExpression:
		return e.evalForLoopExpression(ctx, node, s)
	case *ast.SwitchExpression:
		return e.evalSwitchStatement(ctx, node, s)
	case *ast.PipeExpression:
		return e.evalPipeExpression(ctx, node, s)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(ctx, node, s)
	case *ast.BreakStatement:
		return &object.BreakValue{}

	// Literals
	case *ast.NilLiteral:
		return object.Nil
	case *ast.IntegerLiteral:
		return &object.Int{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return e.evalStringLiteral(ctx, node, s)
	case *ast.RegexpLiteral:
		return &object.Regexp{Value: node.Value, Flags: node.Flags}
	case *ast.ListLiteral:
		return e.evalListLiteral(ctx, node, s)
	case *ast.HashLiteral:
		return e.evalHashLiteral(ctx, node, s)
	case *ast.SetLiteral:
		return e.evalSetLiteral(ctx, node, s)
	}

	panic(fmt.Sprintf("unknown ast node type: %T", node))
}
