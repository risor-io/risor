package evaluator

import (
	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalPostfixExpression(
	s *scope.Scope,
	operator string,
	node *ast.Postfix,
) object.Object {
	switch operator {
	case "++":
		val, ok := s.Get(node.Literal())
		if !ok {
			return object.Errorf("name error: %q is not defined", node.Literal())
		}
		switch arg := val.(type) {
		case *object.Int:
			if err := s.Update(node.Literal(), object.NewInt(arg.Value()+1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		case *object.Float:
			if err := s.Update(node.Literal(), object.NewFloat(arg.Value()+1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		default:
			return object.Errorf("type error: cannot increment %s (type %s)", node.Literal(), arg)
		}
	case "--":
		val, ok := s.Get(node.Literal())
		if !ok {
			return object.Errorf("name error: %q is not defined", node.Literal())
		}
		switch arg := val.(type) {
		case *object.Int:
			if err := s.Update(node.Literal(), object.NewInt(arg.Value()-1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		case *object.Float:
			if err := s.Update(node.Literal(), object.NewFloat(arg.Value()-1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		default:
			return object.Errorf("type error: cannot decrement %s (type %s)", node.Literal(), arg)
		}
	default:
		return object.Errorf("syntax error: unknown operator: %s", operator)
	}
}
