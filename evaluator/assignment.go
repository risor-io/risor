package evaluator

import (
	"context"

	"github.com/myzie/tamarin/ast"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
)

func (e *Evaluator) evalLetStatement(
	ctx context.Context,
	node *ast.LetStatement,
	s *scope.Scope,
) object.Object {
	value := e.Evaluate(ctx, node.Value, s)
	if isError(value) {
		return value
	}
	if err := s.Declare(node.Name.Value, value, false); err != nil {
		return newError(err.Error())
	}
	return value
}

func (e *Evaluator) evalConstStatement(
	ctx context.Context,
	node *ast.ConstStatement,
	s *scope.Scope,
) object.Object {
	value := e.Evaluate(ctx, node.Value, s)
	if isError(value) {
		return value
	}
	if err := s.Declare(node.Name.Value, value, true); err != nil {
		return newError(err.Error())
	}
	return value
}

func (e *Evaluator) evalAssignStatement(
	ctx context.Context,
	a *ast.AssignStatement,
	s *scope.Scope,
) (val object.Object) {
	evaluated := e.Evaluate(ctx, a.Value, s)
	if isError(evaluated) {
		return evaluated
	}
	switch a.Operator {
	case "+=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return newError("name error: %s is not defined", a.Name.String())
		}
		res := e.evalInfix("+=", current, evaluated, s)
		if isError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return newError(err.Error())
		}
		return res

	case "-=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return newError("name error: %s is not defined", a.Name.String())
		}
		res := e.evalInfix("-=", current, evaluated, s)
		if isError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return newError(err.Error())
		}
		return res

	case "*=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return newError("name error: %s is not defined", a.Name.String())
		}
		res := e.evalInfix("*=", current, evaluated, s)
		if isError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return newError(err.Error())
		}
		return res

	case "/=":
		current, ok := s.Get(a.Name.String())
		if !ok {
			return newError("name error: %s is not defined", a.Name.String())
		}
		res := e.evalInfix("/=", current, evaluated, s)
		if isError(res) {
			return res
		}
		if err := s.Update(a.Name.String(), res); err != nil {
			return newError(err.Error())
		}
		return res

	case ":=":
		if err := s.Declare(a.Name.String(), evaluated, false); err != nil {
			return newError(err.Error())
		}

	case "=":
		if err := s.Update(a.Name.String(), evaluated); err != nil {
			return newError(err.Error())
		}
	}
	return evaluated
}
