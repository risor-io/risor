package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalVarStatement(ctx context.Context, node *ast.Var, s *scope.Scope) object.Object {
	ident, expr := node.Value()
	value := e.Evaluate(ctx, expr, s)
	if object.IsError(value) {
		return value
	}
	if err := s.Declare(ident, value, false); err != nil {
		return object.NewError(err)
	}
	return value
}

func (e *Evaluator) evalMultiVarStatement(ctx context.Context, node *ast.MultiVar, s *scope.Scope) object.Object {
	idents, expr := node.Value()
	value := e.Evaluate(ctx, expr, s)
	if object.IsError(value) {
		return value
	}
	switch value := value.(type) {
	case *object.List:
		items := value.Value()
		if len(idents) != len(items) {
			return object.Errorf("eval error: invalid multi variable assignment (list size: %d; identifiers: %d)",
				len(items), len(idents))
		}
		for i, ident := range idents {
			if err := s.Declare(ident, items[i], false); err != nil {
				return object.NewError(err)
			}
		}
	default:
		return object.Errorf("eval error: invalid multi variable assignment")
	}
	return value
}

func (e *Evaluator) evalConstStatement(ctx context.Context, node *ast.Const, s *scope.Scope) object.Object {
	ident, expr := node.Value()
	value := e.Evaluate(ctx, expr, s)
	if object.IsError(value) {
		return value
	}
	if err := s.Declare(ident, value, true); err != nil {
		return object.NewError(err)
	}
	return value
}

func (e *Evaluator) evalAssignStatement(ctx context.Context, a *ast.Assign, s *scope.Scope) object.Object {
	value := e.Evaluate(ctx, a.Value(), s)
	if object.IsError(value) {
		return value
	}
	if a.Index() != nil {
		return e.evalSetItemStatement(ctx, a, value, s)
	}
	name := a.Name()
	switch a.Operator() {
	case "+=":
		current, ok := s.Get(name)
		if !ok {
			return object.Errorf("name error: %q is not defined", name)
		}
		r := e.evalInfix("+=", current, value, s)
		if object.IsError(r) {
			return r
		}
		if err := s.Update(name, r); err != nil {
			return object.NewError(err)
		}
		return r

	case "-=":
		current, ok := s.Get(name)
		if !ok {
			return object.Errorf("name error: %q is not defined", name)
		}
		r := e.evalInfix("-=", current, value, s)
		if object.IsError(r) {
			return r
		}
		if err := s.Update(name, r); err != nil {
			return object.NewError(err)
		}
		return r

	case "*=":
		current, ok := s.Get(name)
		if !ok {
			return object.Errorf("name error: %q is not defined", name)
		}
		r := e.evalInfix("*=", current, value, s)
		if object.IsError(r) {
			return r
		}
		if err := s.Update(name, r); err != nil {
			return object.NewError(err)
		}
		return r

	case "/=":
		current, ok := s.Get(name)
		if !ok {
			return object.Errorf("name error: %q is not defined", name)
		}
		r := e.evalInfix("/=", current, value, s)
		if object.IsError(r) {
			return r
		}
		if err := s.Update(name, r); err != nil {
			return object.NewError(err)
		}
		return r

	case ":=":
		if err := s.Declare(name, value, false); err != nil {
			return object.NewError(err)
		}

	case "=":
		if err := s.Update(name, value); err != nil {
			return object.NewError(err)
		}
	}
	return value
}

func (e *Evaluator) evalSetItemStatement(ctx context.Context, a *ast.Assign, value object.Object, s *scope.Scope) (val object.Object) {
	index := a.Index()
	obj := e.Evaluate(ctx, index.Left(), s)
	if object.IsError(obj) {
		return obj
	}
	container, ok := obj.(object.Container)
	if !ok {
		return object.Errorf("type error: %s is not a container", obj.Type())
	}
	indexObj := e.Evaluate(ctx, index.Index(), s)
	if object.IsError(indexObj) {
		return indexObj
	}
	switch a.Operator() {
	case "=":
		if err := container.SetItem(indexObj, value); err != nil {
			return err
		}
	default:
		return object.Errorf("eval error: invalid set item operator: %q", a.Operator)
	}
	return object.Nil
}
