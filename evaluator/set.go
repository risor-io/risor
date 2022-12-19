package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
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
	if err := set.Add(items...); err != nil {
		return object.NewError(err.Error())
	}
	return set
}
