package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalCallExpression(ctx context.Context, node *ast.Call, s *scope.Scope) object.Object {
	function := e.Evaluate(ctx, node.Function(), s)
	if object.IsError(function) {
		return function
	}
	if builtin, ok := function.(*object.Builtin); ok {
		if builtin.IsErrorHandler() {
			return e.applyFunction(ctx, s, function,
				e.evalExpressionsIgnoreErrors(ctx, node.Arguments(), s))
		}
	}
	args := e.evalExpressions(ctx, node.Arguments(), s)
	if len(args) == 1 && object.IsError(args[0]) {
		return args[0]
	}
	return e.applyFunction(ctx, s, function, args)
}

func (e *Evaluator) evalObjectCallExpression(ctx context.Context, call *ast.ObjectCall, s *scope.Scope) object.Object {
	obj := e.Evaluate(ctx, call.Object(), s)
	if object.IsError(obj) {
		return obj
	}
	callExpr := call.Call()
	if method, ok := callExpr.(*ast.Call); ok {
		args := e.evalExpressions(ctx, method.Arguments(), s)
		if len(args) == 1 && object.IsError(args[0]) {
			return args[0]
		}
		funcName := method.Function().String()
		return e.evalObjectCall(ctx, s, obj, funcName, args)
	}
	return object.Errorf("failed to evaluate object call")
}

func (e *Evaluator) evalObjectCall(ctx context.Context, s *scope.Scope, obj object.Object, method string, args []object.Object) object.Object {
	if attr, found := obj.GetAttr(method); found {
		if object.IsError(attr) {
			return attr
		}
		return e.applyFunction(ctx, s, attr, args)
	}
	return object.Errorf("attribute error: %s has no attribute \"%s\"", obj.Type(), method)
}

func (e *Evaluator) evalGetAttributeExpression(ctx context.Context, node *ast.GetAttr, s *scope.Scope) object.Object {
	obj := e.Evaluate(ctx, node.Object(), s)
	if object.IsError(obj) {
		return obj
	}
	name := node.Name()
	if attr, found := obj.GetAttr(name); found {
		return attr
	}
	return object.Errorf("attribute error: %s object has no attribute \"%s\"", obj.Type(), name)
}
