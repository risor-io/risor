package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalIn(ctx context.Context, node *ast.In, s *scope.Scope) object.Object {
	left := e.Evaluate(ctx, node.Left(), s)
	if object.IsError(left) {
		return left
	}
	right := e.Evaluate(ctx, node.Right(), s)
	if object.IsError(right) {
		return right
	}
	container, ok := right.(object.Container)
	if !ok {
		return object.Errorf("eval error: right hand side of 'in' operator must be a container")
	}
	return container.Contains(left)
}
