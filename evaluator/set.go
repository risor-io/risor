package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

func (e *Evaluator) evalSetLiteral(
	ctx context.Context,
	node *ast.SetLiteral,
	s *scope.Scope,
) object.Object {
	set := object.NewSetWithSize(len(node.Items))
	items := make([]object.Object, 0, len(node.Items))
	for _, itemNode := range node.Items {
		item := e.Evaluate(ctx, itemNode, s)
		if isError(item) {
			return item
		}
		items = append(items, item)
	}
	set.Add(items...)
	return set
}
