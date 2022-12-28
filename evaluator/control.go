package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// evalIfExpression handles an `if` expression, running the block
// if the condition matches, and running any optional else block
// otherwise.
func (e *Evaluator) evalIfExpression(ctx context.Context, ie *ast.If, s *scope.Scope) object.Object {
	condition := e.Evaluate(ctx, ie.Condition(), s)
	if object.IsError(condition) {
		return condition
	}
	if object.IsTruthy(condition) {
		return e.Evaluate(ctx, ie.Consequence(), s)
	} else if ie.Alternative() != nil {
		return e.Evaluate(ctx, ie.Alternative(), s)
	}
	return object.Nil
}

func (e *Evaluator) evalForLoopExpression(ctx context.Context, fle *ast.For, s *scope.Scope) object.Object {

	forScope := s.NewChild(scope.Opts{Name: "for"})
	loopScope := forScope.NewChild(scope.Opts{Name: "for-loop"})

	// Evaluate the initialization statement if there is one
	init := fle.Init()
	if init != nil {
		if res := e.Evaluate(ctx, init, forScope); object.IsError(res) {
			return res
		}
	}

	// The for loop evaluates to this value. It is set to the last value
	// evaluated in the for loop block.
	var latestValue object.Object = object.Nil

	// This is a simple for loop, like "for { ... }". It will run until
	// an error occurs or a break or return statement is encountered.
	if fle.IsSimpleLoop() {

	simpleLoop:
		for {
			loopScope.Clear()
			result := e.Evaluate(ctx, fle.Consequence(), loopScope)
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.BreakValue:
				break simpleLoop
			}
			latestValue = result
		}
		return latestValue
	}

	// This is a standard for loop that runs until a specified condition is met
loop:
	for {
		loopScope.Clear()
		// Evaluate the condition
		condition := e.Evaluate(ctx, fle.Condition(), forScope)
		if object.IsError(condition) {
			return condition
		}
		if object.IsTruthy(condition) {
			// Evaluate the block
			rt := e.Evaluate(ctx, fle.Consequence(), loopScope)
			switch rt := rt.(type) {
			case *object.Error:
				return rt
			case *object.ReturnValue:
				return rt
			case *object.BreakValue:
				break loop
			}
			latestValue = rt
		} else {
			break
		}
		// Evaluate the post statement (usually used to increment a counter)
		if fle.Post() != nil {
			if res := e.Evaluate(ctx, fle.Post(), forScope); object.IsError(res) {
				return res
			}
		}
	}
	return latestValue
}

func (e *Evaluator) evalSwitchStatement(ctx context.Context, se *ast.Switch, s *scope.Scope) object.Object {
	value := e.Evaluate(ctx, se.Value(), s)
	if object.IsError(value) {
		return value
	}
	for _, opt := range se.Choices() {
		if opt.IsDefault() {
			continue
		}
		for _, val := range opt.Expressions() {
			out := e.Evaluate(ctx, val, s)
			if object.IsError(out) {
				return out
			}
			if object.Equals(value, out) {
				return e.evalBlockStatement(ctx, opt.Block(), s)
			}
		}
	}
	// No match found, so run the default block if there is one
	for _, opt := range se.Choices() {
		if opt.IsDefault() {
			return e.evalBlockStatement(ctx, opt.Block(), s)
		}
	}
	return object.Nil
}

func prependObject(slice []object.Object, obj object.Object) []object.Object {
	slice = append(slice, nil)
	copy(slice[1:], slice)
	slice[0] = obj
	return slice
}

func (e *Evaluator) evalPipeExpression(ctx context.Context, pe *ast.Pipe, s *scope.Scope) object.Object {
	exprs := pe.Expressions()
	if len(exprs) < 2 {
		return object.Errorf("eval error: invalid pipe expression (got only %d arguments)", len(exprs))
	}
	// Evaluate the expression preceding the first pipe operator
	nextArg := e.Evaluate(ctx, exprs[0], s)
	if object.IsError(nextArg) {
		return nextArg
	}
	// Evaluate the rest of the pipe expression
	for i, expr := range exprs[1:] {
		switch expression := expr.(type) {
		case *ast.Call:
			// Can't use evalCallExpression because we need to customize argument handling
			function := e.Evaluate(ctx, expression.Function(), s)
			if object.IsError(function) {
				return function
			}
			// Resolve the call arguments
			var args []object.Object
			if len(expression.Arguments()) > 0 {
				args = e.evalExpressions(ctx, expression.Arguments(), s)
				if len(args) == 1 && object.IsError(args[0]) {
					return args[0]
				}
			}
			// Prepend any arguments present from the previous pipeline stage and then run the call
			if nextArg != nil {
				args = prependObject(args, nextArg)
			}
			res := e.applyFunction(ctx, s, function, args)
			if object.IsError(res) {
				return res
			}
			// Save the output as arguments for the next stage
			nextArg = res
		case *ast.ObjectCall:
			// Resolve the object
			obj := e.Evaluate(ctx, expression.Object(), s)
			if object.IsError(obj) {
				return obj
			}
			// Resolve the call arguments
			callExpr := expression.Call().(*ast.Call)
			var args []object.Object
			if len(callExpr.Arguments()) > 0 {
				args = e.evalExpressions(ctx, callExpr.Arguments(), s)
				if len(args) == 1 && object.IsError(args[0]) {
					return args[0]
				}
			}
			// Prepend any arguments present from the previous pipeline stage and then run the call
			if nextArg != nil {
				args = prependObject(args, nextArg)
			}
			method, ok := callExpr.Function().(*ast.Ident)
			if !ok {
				return object.Errorf("invalid function in pipe expression: %v", callExpr.Function)
			}
			res := e.evalObjectCall(ctx, s, obj, method.Literal(), args)
			if object.IsError(res) {
				return res
			}
			// Save the output as arguments for the next stage
			nextArg = res
		default:
			// Evaluate the expression. We expect it to evaluate to a function, or, if its the
			// first stage in the pipeline, to the first argument to be passed to the next stage.
			obj := e.Evaluate(ctx, expression, s)
			if object.IsError(obj) {
				return obj
			}
			switch obj := obj.(type) {
			case *object.Function, *object.Builtin:
				var args []object.Object
				if nextArg != nil {
					args = []object.Object{nextArg}
				}
				res := e.applyFunction(ctx, s, obj, args)
				if object.IsError(res) {
					return res
				}
				// Save the output as arguments for the next stage
				nextArg = res
			default:
				if i == 0 {
					// Save the output as arguments for the next stage
					nextArg = obj
				} else {
					return object.Errorf("type error: unexpected %s object in pipe expression", obj.Type())
				}
			}
		}
	}
	if nextArg != nil {
		return nextArg
	}
	return object.Nil
}

func (e *Evaluator) evalReturnStatement(ctx context.Context, node *ast.Return, s *scope.Scope) object.Object {
	if node.Value() == nil {
		return object.Nil // Should we adjust the parser output in this case?
	}
	value := e.Evaluate(ctx, node.Value(), s)
	if object.IsError(value) {
		return value
	}
	return object.NewReturnValue(value)
}

func (e *Evaluator) upwrapReturnValue(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value()
	}
	return obj
}
