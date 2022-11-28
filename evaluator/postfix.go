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
			return newError("name error: %s is not defined", node.Token.Literal)
		}
		switch arg := val.(type) {
		case *object.Integer:
			if err := s.Update(node.Token.Literal, &object.Integer{Value: arg.Value + 1}); err != nil {
				return newError(err.Error())
			}
			return arg
		default:
			return newError("type error: cannot increment %s (type %s)", node.Token.Literal, arg)
		}
	case "--":
		val, ok := s.Get(node.Token.Literal)
		if !ok {
			return newError("name error: %s is not defined", node.Token.Literal)
		}
		switch arg := val.(type) {
		case *object.Integer:
			if err := s.Update(node.Token.Literal, &object.Integer{Value: arg.Value - 1}); err != nil {
				return newError(err.Error())
			}
			return arg
		default:
			return newError("type error: cannot decrement %s (type %s)", node.Token.Literal, arg)
		}
	default:
		return newError("syntax error: unknown operator: %s", operator)
	}
}
