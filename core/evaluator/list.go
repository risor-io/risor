package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalListLiteral(
	ctx context.Context,
	node *ast.List,
	s *scope.Scope,
) object.Object {
	elements := e.evalExpressions(ctx, node.Items(), s)
	if len(elements) == 1 && object.IsError(elements[0]) {
		return elements[0]
	}
	return object.NewList(elements)
}
