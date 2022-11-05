package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

// evalIfExpression handles an `if` expression, running the block
// if the condition matches, and running any optional else block
// otherwise.
func (e *Evaluator) evalIfExpression(
	ctx context.Context,
	ie *ast.IfExpression,
	s *scope.Scope,
) object.Object {
	condition := e.Evaluate(ctx, ie.Condition, s)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return e.Evaluate(ctx, ie.Consequence, s)
	} else if ie.Alternative != nil {
		return e.Evaluate(ctx, ie.Alternative, s)
	} else {
		return object.NULL
	}
}

func (e *Evaluator) evalForLoopExpression(
	ctx context.Context,
	fle *ast.ForLoopExpression,
	s *scope.Scope,
) object.Object {
	rt := &object.Boolean{Value: true}
	nestedScope := s.NewChild(scope.Opts{Name: "for"})
	if fle.InitStatement != nil {
		if res := e.Evaluate(ctx, fle.InitStatement, nestedScope); isError(res) {
			return res
		}
	}
	for {
		condition := e.Evaluate(ctx, fle.Condition, nestedScope)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			rt := e.Evaluate(ctx, fle.Consequence, nestedScope)
			if !isError(rt) && (rt.Type() == object.RETURN_VALUE_OBJ ||
				rt.Type() == object.ERROR_OBJ) {
				return rt
			}
		} else {
			break
		}
		if fle.PostStatement != nil {
			if res := e.Evaluate(ctx, fle.PostStatement, nestedScope); isError(res) {
				return res
			}
		}
	}
	return rt
}

func (e *Evaluator) evalSwitchStatement(
	ctx context.Context,
	se *ast.SwitchExpression,
	s *scope.Scope,
) object.Object {
	// Get the value
	obj := e.Evaluate(ctx, se.Value, s)
	// Try all the choices
	for _, opt := range se.Choices {
		// skipping the default-case, which we'll handle later
		if opt.Default {
			continue
		}
		// Look at any expression we've got in this case.
		for _, val := range opt.Expr {
			// Get the value of the case
			out := e.Evaluate(ctx, val, s)
			// Is it a literal match?
			if obj.Type() == out.Type() && (obj.Inspect() == out.Inspect()) {
				// Evaluate the block and return the value
				blockOut := e.evalBlockStatement(ctx, opt.Block, s)
				return blockOut
			}
			// Is it a regexp-match?
			if out.Type() == object.REGEXP_OBJ {
				m := matches(obj, out, s)
				if m == object.TRUE {
					// Evaluate the block and return the value
					out := e.evalBlockStatement(ctx, opt.Block, s)
					return out
				}
			}
		}
	}
	// No match? Handle default if present.
	for _, opt := range se.Choices {
		// skip default
		if opt.Default {
			out := e.evalBlockStatement(ctx, opt.Block, s)
			return out
		}
	}
	return nil
}

func prependObject(slice []object.Object, obj object.Object) []object.Object {
	slice = append(slice, nil)
	copy(slice[1:], slice)
	slice[0] = obj
	return slice
}

func (e *Evaluator) evalPipeExpression(
	ctx context.Context,
	pe *ast.PipeExpression,
	s *scope.Scope,
) object.Object {
	var nextArg object.Object
	for i, expr := range pe.Arguments {
		switch expression := expr.(type) {
		case *ast.CallExpression:
			// Can't use evalCallExpression because we need to customize argument handling
			function := e.Evaluate(ctx, expression.Function, s)
			if isError(function) {
				return function
			}
			// Resolve the call arguments
			var args []object.Object
			if len(expression.Arguments) > 0 {
				args = e.evalExpressions(ctx, expression.Arguments, s)
				if len(args) == 1 && isError(args[0]) {
					return args[0]
				}
			}
			// Prepend any arguments present from the previous pipeline stage and then run the call
			args = prependObject(args, nextArg)
			res := e.applyFunction(ctx, s, function, args)
			if isError(res) {
				return res
			}
			// Save the output as arguments for the next stage
			nextArg = res
		case *ast.ObjectCallExpression:
			// Resolve the object
			obj := e.Evaluate(ctx, expression.Object, s)
			if isError(obj) {
				return obj
			}
			// Resolve the call arguments
			callExpr := expression.Call.(*ast.CallExpression)
			var args []object.Object
			if len(callExpr.Arguments) > 0 {
				args = e.evalExpressions(ctx, callExpr.Arguments, s)
				if len(args) == 1 && isError(args[0]) {
					return args[0]
				}
			}
			// Prepend any arguments present from the previous pipeline stage and then run the call
			args = prependObject(args, nextArg)
			method, ok := callExpr.Function.(*ast.Identifier)
			if !ok {
				return newError("invalid function in pipe expression: %v", callExpr.Function)
			}
			res := e.evalObjectCall(ctx, s, obj, method.Value, args)
			if isError(res) {
				return res
			}
			// Save the output as arguments for the next stage
			nextArg = res
		default:
			// Evaluate the expression. We expect it to evaluate to a function, or, if its the
			// first stage in the pipeline, to the first argument to be passed to the next stage.
			obj := e.Evaluate(ctx, expression, s)
			if isError(obj) {
				return obj
			}
			switch obj := obj.(type) {
			case *object.Function, *object.Builtin:
				var args []object.Object
				if nextArg != nil {
					args = []object.Object{nextArg}
				}
				res := e.applyFunction(ctx, s, obj, args)
				if isError(res) {
					return res
				}
				// Save the output as arguments for the next stage
				nextArg = res
			default:
				if i == 0 {
					// Save the output as arguments for the next stage
					nextArg = obj
				} else {
					return newError("type error: unexpected %s object in pipe expression", obj.Type())
				}
			}
		}
	}
	if nextArg != nil {
		return nextArg
	}
	return object.NULL
}

func (e *Evaluator) evalReturnStatement(
	ctx context.Context,
	node *ast.ReturnStatement,
	s *scope.Scope,
) object.Object {
	value := e.Evaluate(ctx, node.ReturnValue, s)
	if isError(value) {
		return value
	}
	return &object.ReturnValue{Value: value}
}

func (e *Evaluator) upwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
