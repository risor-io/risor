package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
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
	case object.True:
		return object.False
	case object.False:
		return object.True
	default:
		return newError("type error: expected boolean to follow ! operator (got %s)", right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Int:
		return &object.Int{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return newError("type error: expected int or float to follow - operator (got %s)", right.Type())
	}
}
