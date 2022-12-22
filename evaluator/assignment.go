package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalVarStatement(
	ctx context.Context,
	node *ast.VarStatement,
	s *scope.Scope,
) object.Object {
	value := e.Evaluate(ctx, node.Value, s)
	if object.IsError(value) {
		return value
	}
	if err := s.Declare(node.Name.Value, value, false); err != nil {
		return object.Errorf(err.Error())
	}
	return value
}

func (e *Evaluator) evalConstStatement(
	ctx context.Context,
	node *ast.ConstStatement,
	s *scope.Scope,
) object.Object {
	value := e.Evaluate(ctx, node.Value, s)
	if object.IsError(value) {
		return value
	}
	if err := s.Declare(node.Name.Value, value, true); err != nil {
		return object.Errorf(err.Error())
	}
	return value
}

func (e *Evaluator) evalAssignStatement(
	ctx context.Context,
	a *ast.AssignStatement,
	s *scope.Scope,
) (val object.Object) {
	evaluated := e.Evaluate(ctx, a.Value, s)
	if object.IsError(evaluated) {
		return evaluated
	}
	if a.Index != nil {
		return e.evalSetItemStatement(ctx, a, evaluated, s)
	}
	switch a.Operator {
	case "+=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return object.Errorf("name error: %q is not defined", a.Name.String())
		}
		res := e.evalInfix("+=", current, evaluated, s)
		if object.IsError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return object.Errorf(err.Error())
		}
		return res

	case "-=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return object.Errorf("name error: %q is not defined", a.Name.String())
		}
		res := e.evalInfix("-=", current, evaluated, s)
		if object.IsError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return object.Errorf(err.Error())
		}
		return res

	case "*=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return object.Errorf("name error: %q is not defined", a.Name.String())
		}
		res := e.evalInfix("*=", current, evaluated, s)
		if object.IsError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return object.Errorf(err.Error())
		}
		return res

	case "/=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return object.Errorf("name error: %q is not defined", a.Name.String())
		}
		res := e.evalInfix("/=", current, evaluated, s)
		if object.IsError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return object.Errorf(err.Error())
		}
		return res

	case ":=":
		if err := s.Declare(a.Name.String(), evaluated, false); err != nil {
			return object.Errorf(err.Error())
		}

	case "=":
		if err := s.Update(a.Name.String(), evaluated); err != nil {
			return object.Errorf(err.Error())
		}
	}
	return evaluated
}

func (e *Evaluator) evalSetItemStatement(
	ctx context.Context,
	a *ast.AssignStatement,
	value object.Object,
	s *scope.Scope,
) (val object.Object) {
	obj := e.Evaluate(ctx, a.Index.Left, s)
	if object.IsError(obj) {
		return obj
	}
	container, ok := obj.(object.Container)
	if !ok {
		return object.Errorf("type error: %s is not a container", obj.Type())
	}
	index := e.Evaluate(ctx, a.Index.Index, s)
	if object.IsError(index) {
		return index
	}
	switch a.Operator {
	case "=":
		if err := container.SetItem(index, value); err != nil {
			return err
		}
	default:
		return object.Errorf("eval error: invalid set item operator: %q", a.Operator)
	}
	return object.Nil
}
