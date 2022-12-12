package evaluator

import (
	"context"
	"fmt"

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
	if method, ok := call.Call.(*ast.CallExpression); ok {
		args := e.evalExpressions(ctx, call.Call.(*ast.CallExpression).Arguments, s)
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
	case *object.List:
		switch method {
		case "map":
			return e.execListMap(ctx, s, obj, args...)
		case "filter":
			return e.execListFilter(ctx, s, obj, args...)
		case "each":
			return e.execListEach(ctx, s, obj, args...)
		}
	case *object.Module:
		moduleScope := obj.Scope.(*scope.Scope)
		moduleFunc, ok := moduleScope.Get(method)
		if !ok {
			return newError("attribute error: module %s has no attribute %s",
				obj.Name, method)
		}
		return e.applyFunction(ctx, s, moduleFunc, args)
	}
	result := obj.InvokeMethod(method, args...)
	if result != nil {
		return result
	}
	return newError("failed to invoke method: %s", method)
}

func (e *Evaluator) execListMap(
	ctx context.Context,
	s *scope.Scope,
	array *object.List,
	args ...object.Object,
) object.Object {
	if len(args) != 1 {
		return newError(fmt.Sprintf("expected one argument to map call; got %d", len(args)))
	}
	mapFunc := args[0]

	var numParameters int
	switch mapFunc := mapFunc.(type) {
	case *object.Function:
		numParameters = len(mapFunc.Parameters)
	case *object.Builtin:
		numParameters = 1 // This is not always correct. Need a lookup table?
	default:
		return newError("type error: %s is not callable", mapFunc.Type())
	}

	if numParameters < 1 || numParameters > 2 {
		return newError("error: function parameter count is incompatible with map call")
	}

	var index object.Int
	mapArgs := make([]object.Object, 2)
	result := make([]object.Object, 0, len(array.Items))
	for i, value := range array.Items {
		index.Value = int64(i)
		mapArgs[0] = &index
		mapArgs[1] = value
		var outputValue object.Object
		if numParameters == 1 {
			outputValue = e.applyFunction(ctx, s, mapFunc, mapArgs[1:])
		} else {
			outputValue = e.applyFunction(ctx, s, mapFunc, mapArgs)
		}
		if isError(outputValue) {
			return outputValue
		}
		result = append(result, outputValue)
	}
	return &object.List{Items: result}
}

func (e *Evaluator) execListEach(
	ctx context.Context,
	s *scope.Scope,
	array *object.List,
	args ...object.Object,
) object.Object {
	if len(args) != 1 {
		return newError(fmt.Sprintf("expected one argument to each call; got %d", len(args)))
	}
	mapFunc := args[0]

	var numParameters int
	switch mapFunc := mapFunc.(type) {
	case *object.Function:
		numParameters = len(mapFunc.Parameters)
	case *object.Builtin:
		numParameters = 1 // This is not always correct. Need a lookup table?
	default:
		return newError("type error: %s is not callable", mapFunc.Type())
	}
	if numParameters != 1 {
		return newError("error: function parameter count is incompatible with each call")
	}
	for _, value := range array.Items {
		callArgs := []object.Object{value}
		outputValue := e.applyFunction(ctx, s, mapFunc, callArgs)
		if isError(outputValue) {
			return outputValue
		}
	}
	return object.Nil
}

func (e *Evaluator) execListFilter(
	ctx context.Context,
	s *scope.Scope,
	array *object.List,
	args ...object.Object,
) object.Object {
	if len(args) != 1 {
		return newError(fmt.Sprintf("expected one argument to filter call; got %d", len(args)))
	}
	mapFunc := args[0]

	var numParameters int
	switch mapFunc := mapFunc.(type) {
	case *object.Function:
		numParameters = len(mapFunc.Parameters)
	case *object.Builtin:
		return newError("type error: built-in function %s may not be used for filter call",
			mapFunc.Inspect())
	default:
		return newError("type error: %s is not callable", mapFunc.Type())
	}
	if numParameters != 1 {
		return newError("type error: expected a function with a single parameter")
	}

	filterArgs := make([]object.Object, 1)
	var result []object.Object
	for _, value := range array.Items {
		filterArgs[0] = value
		decision := e.applyFunction(ctx, s, mapFunc, filterArgs)
		if isError(decision) {
			return decision
		}
		if isTruthy(decision) {
			result = append(result, value)
		}
	}
	return &object.List{Items: result}
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
	case *object.Map:
		return obj.Get(attrName)
	case *object.Module:
		s := obj.Scope.(*scope.Scope)
		result, ok := s.Get(attrName)
		if !ok {
			return newError("attribute error: %s object has no attribute %s",
				obj.Type(), attrName)
		}
		return result
	default:
		return newError("attribute error: %s object has no attribute %s",
			obj.Type(), attrName)
	}
}
