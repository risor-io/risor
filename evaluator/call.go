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
	if isError(function) {
		return function
	}
	args := e.evalExpressions(ctx, node.Arguments, s)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return e.applyFunction(ctx, s, function, args)
}

// evalObjectCallExpression invokes methods against objects.
func (e *Evaluator) evalObjectCallExpression(
	ctx context.Context,
	call *ast.ObjectCallExpression,
	s *scope.Scope,
) object.Object {
	obj := e.Evaluate(ctx, call.Object, s)
	if isError(obj) {
		return obj
	}
	if method, ok := call.Call.(*ast.CallExpression); ok {
		args := e.evalExpressions(ctx, call.Call.(*ast.CallExpression).Arguments, s)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		funcName := method.Function.String()
		return e.evalObjectCall(ctx, s, obj, funcName, args)
	}
	return newError("Failed to invoke method: %s",
		call.Call.(*ast.CallExpression).Function.String())
}

func (e *Evaluator) evalObjectCall(
	ctx context.Context,
	s *scope.Scope,
	obj object.Object,
	method string,
	args []object.Object,
) object.Object {
	switch obj := obj.(type) {
	case *object.Module:
		moduleScope := obj.Scope.(*scope.Scope)
		moduleFunc, ok := moduleScope.Get(method)
		if !ok {
			return newError("attribute error: module %s has no attribute \"%s\"", obj.Name, method)
		}
		return e.applyFunction(ctx, s, moduleFunc, args)
	}
	if attr, found := obj.GetAttr(method); found {
		return e.applyFunction(ctx, s, attr, args)
	}
	return newError("attribute error: %s has no attribute \"%s\"", obj.Type(), method)
}

func (e *Evaluator) evalGetAttributeExpression(
	ctx context.Context,
	node *ast.GetAttributeExpression,
	s *scope.Scope,
) object.Object {
	obj := e.Evaluate(ctx, node.Object, s)
	if isError(obj) {
		return obj
	}
	attrName := node.Attribute.Token.Literal
	switch obj := obj.(type) {
	case *object.Module:
		s := obj.Scope.(*scope.Scope)
		result, ok := s.Get(attrName)
		if !ok {
			return newError("attribute error: %s object has no attribute \"%s\"",
				obj.Type(), attrName)
		}
		return result
	default:
		if attr, found := obj.GetAttr(attrName); found {
			return attr
		}
		return newError("attribute error: %s object has no attribute \"%s\"",
			obj.Type(), attrName)
	}
}
