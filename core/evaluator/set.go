package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalSetLiteral(
	ctx context.Context,
	node *ast.Set,
	s *scope.Scope,
) object.Object {
	set := object.NewSetWithSize(len(node.Items()))
	items := make([]object.Object, 0, len(node.Items()))
	for _, itemNode := range node.Items() {
		item := e.Evaluate(ctx, itemNode, s)
		if object.IsError(item) {
			return item
		}
		items = append(items, item)
	}
	if err := set.Add(items...); err != object.Nil {
		return err
	}
	return set
}
