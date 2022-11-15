package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

func (e *Evaluator) evalPrefixExpression(
	ctx context.Context,
	node *ast.PrefixExpression,
	s *scope.Scope,
) object.Object {
	right := e.Evaluate(ctx, node.Right, s)
	if isError(right) {
		return right
	}
	operator := node.Operator
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("syntax error: unknown operator: %s", operator)
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return newError("type error: bad operand type for unary -: %s", right.Type())
	}
}
