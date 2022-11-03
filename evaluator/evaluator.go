// Package evaluator contains the core of our interpreter, which walks
// the AST produced by the parser and evaluates the user-submitted program.
package evaluator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/myzie/tamarin/ast"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
)

type Opts struct {
	Importer Importer
}

type Evaluator struct {
	importer Importer
}

// New returns a new Evaluator
func New(opts Opts) *Evaluator {
	return &Evaluator{importer: opts.Importer}
}

// Evaluate an AST node. The context can be used to cancel a running evaluation.
func (e *Evaluator) Evaluate(ctx context.Context, node ast.Node, s *scope.Scope) object.Object {

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
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ImportStatement:
		return e.evalImportStatement(ctx, node, s)

	// Assignment
	case *ast.LetStatement:
		return e.evalLetStatement(ctx, node, s)
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

	// Literals
	case *ast.NullLiteral:
		return object.NULL
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.RegexpLiteral:
		return &object.Regexp{Value: node.Value, Flags: node.Flags}
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(ctx, node, s)
	case *ast.HashLiteral:
		return e.evalHashLiteral(ctx, node, s)
	case *ast.SetLiteral:
		return e.evalSetLiteral(ctx, node, s)
	}

	panic(fmt.Sprintf("unknown ast node type: %s", reflect.TypeOf(node)))
}
