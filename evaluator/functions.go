package evaluator

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/cloudcmds/tamarin/stack"
)

func (e *Evaluator) evalFunctionLiteral(ctx context.Context, node *ast.Func, s *scope.Scope) object.Object {
	if node.Name() != nil {
		name := node.Name().String()
		fn := object.NewFunction(name, node.Parameters(), node.Body(), node.Defaults(), s)
		if err := s.Declare(name, fn, true); err != nil {
			return object.NewError(err)
		}
		return object.Nil
	}
	return object.NewFunction("", node.Parameters(), node.Body(), node.Defaults(), s)
}

func (e *Evaluator) evalFunctionDefinition(ctx context.Context, node *ast.Func, s *scope.Scope) object.Object {
	name := node.Literal()
	fn := object.NewFunction(name, node.Parameters(), node.Body(), node.Defaults(), s)
	if err := s.Declare(name, fn, true); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func (e *Evaluator) applyFunction(ctx context.Context, s *scope.Scope, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// Use the function's scope, not the current execution scope! This is
		// what enables closures to work as expected!
		nestedScope, err := e.newFunctionScope(ctx, fn.Scope().(*scope.Scope), fn, args)
		if err != nil {
			return object.NewError(err)
		}
		funcBody := fn.Body()
		e.stack.Push(stack.NewFrame(stack.FrameOpts{
			Name:  fn.Name(),
			Scope: nestedScope,
		}))
		defer e.stack.Pop()
		result := e.Evaluate(ctx, funcBody, nestedScope)
		return e.upwrapReturnValue(result)
	case *object.Builtin:
		e.stack.Push(stack.NewFrame(stack.FrameOpts{
			Name:  fn.Key(),
			Scope: s,
		}))
		defer e.stack.Pop()
		if priorityBuiltin, found := e.builtins[fn.Key()]; found {
			// This is a priority builtin, possibly an override, so
			// we should use this one
			return priorityBuiltin.Call(ctx, args...)
		}
		// This is a non-priority builtin
		return fn.Call(ctx, args...)
	default:
		return object.Errorf("type error: %s is not callable", fn.Type())
	}
}

func (e *Evaluator) newFunctionScope(ctx context.Context, s *scope.Scope, fn *object.Function, args []object.Object) (*scope.Scope, error) {
	declared := map[string]bool{}
	nestedScope := s.NewChild(scope.Opts{Name: "function"})
	for key, val := range fn.Defaults() {
		evaluatedValue := e.Evaluate(ctx, val, s)
		if object.IsError(evaluatedValue) {
			return nil, fmt.Errorf("failed to evaluate parameter: %s", key)
		}
		if err := nestedScope.Declare(key, evaluatedValue, false); err != nil {
			return nil, err
		}
		declared[key] = true
	}
	if len(fn.Defaults()) == 0 && len(args) != len(fn.Parameters()) {
		return nil, fmt.Errorf("type error: function expected %d arguments (%d given)",
			len(fn.Parameters()), len(args))
	}
	for paramIdx, param := range fn.Parameters() {
		if paramIdx < len(args) {
			name := param.String()
			if declared[name] {
				if err := nestedScope.Update(name, args[paramIdx]); err != nil {
					return nil, err
				}
			} else {
				if err := nestedScope.Declare(name, args[paramIdx], false); err != nil {
					return nil, err
				}
			}
		} else {
			break
		}
	}
	return nestedScope, nil
}
