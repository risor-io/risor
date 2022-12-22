package evaluator

import (
	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalPostfixExpression(
	s *scope.Scope,
	operator string,
	node *ast.PostfixExpression,
) object.Object {
	switch operator {
	case "++":
		val, ok := s.Get(node.Token.Literal)
		if !ok {
			return object.Errorf("name error: %q is not defined", node.Token.Literal)
		}
		switch arg := val.(type) {
		case *object.Int:
			if err := s.Update(node.Token.Literal, object.NewInt(arg.Value()+1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		case *object.Float:
			if err := s.Update(node.Token.Literal, object.NewFloat(arg.Value()+1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		default:
			return object.Errorf("type error: cannot increment %s (type %s)", node.Token.Literal, arg)
		}
	case "--":
		val, ok := s.Get(node.Token.Literal)
		if !ok {
			return object.Errorf("name error: %q is not defined", node.Token.Literal)
		}
		switch arg := val.(type) {
		case *object.Int:
			if err := s.Update(node.Token.Literal, object.NewInt(arg.Value()-1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		case *object.Float:
			if err := s.Update(node.Token.Literal, object.NewFloat(arg.Value()-1)); err != nil {
				return object.Errorf(err.Error())
			}
			return arg
		default:
			return object.Errorf("type error: cannot decrement %s (type %s)", node.Token.Literal, arg)
		}
	default:
		return object.Errorf("syntax error: unknown operator: %s", operator)
	}
}
