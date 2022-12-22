package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalCallExpression(
	ctx context.Context,
	node *ast.CallExpression,
	s *scope.Scope,
) object.Object {
	function := e.Evaluate(ctx, node.Function, s)
	if object.IsError(function) {
		return function
	}
	if builtin, ok := function.(*object.Builtin); ok {
		if builtin.IsErrorHandler() {
			return e.applyFunction(ctx, s, function,
				e.evalExpressionsIgnoreErrors(ctx, node.Arguments, s))
		}
	}
	args := e.evalExpressions(ctx, node.Arguments, s)
	if len(args) == 1 && object.IsError(args[0]) {
		return args[0]
	}
	return e.applyFunction(ctx, s, function, args)
}

func (e *Evaluator) evalObjectCallExpression(
	ctx context.Context,
	call *ast.ObjectCallExpression,
	s *scope.Scope,
) object.Object {
	obj := e.Evaluate(ctx, call.Object, s)
	if object.IsError(obj) {
		return obj
	}
	if method, ok := call.Call.(*ast.CallExpression); ok {
		args := e.evalExpressions(ctx, call.Call.(*ast.CallExpression).Arguments, s)
		if len(args) == 1 && object.IsError(args[0]) {
			return args[0]
		}
		funcName := method.Function.String()
		return e.evalObjectCall(ctx, s, obj, funcName, args)
	}
	return object.Errorf("Failed to invoke method: %s",
		call.Call.(*ast.CallExpression).Function.String())
}

func (e *Evaluator) evalObjectCall(
	ctx context.Context,
	s *scope.Scope,
	obj object.Object,
	method string,
	args []object.Object,
) object.Object {
	if attr, found := obj.GetAttr(method); found {
		if object.IsError(attr) {
			return attr
		}
		return e.applyFunction(ctx, s, attr, args)
	}
	return object.Errorf("attribute error: %s has no attribute \"%s\"", obj.Type(), method)
}

func (e *Evaluator) evalGetAttributeExpression(
	ctx context.Context,
	node *ast.GetAttributeExpression,
	s *scope.Scope,
) object.Object {
	obj := e.Evaluate(ctx, node.Object, s)
	if object.IsError(obj) {
		return obj
	}
	attrName := node.Attribute.Token.Literal
	if attr, found := obj.GetAttr(attrName); found {
		if object.IsError(attr) {
			return attr
		}
		return attr
	}
	return object.Errorf("attribute error: %s object has no attribute \"%s\"",
		obj.Type(), attrName)
}
