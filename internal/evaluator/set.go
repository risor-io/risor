package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

func (e *Evaluator) evalSetLiteral(
	ctx context.Context,
	node *ast.SetLiteral,
	s *scope.Scope,
) object.Object {
	items := make(map[object.HashKey]object.Object, len(node.Items))
	for _, itemNode := range node.Items {
		key := e.Evaluate(ctx, itemNode, s)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("type error: %s object is unhashable", key.Type())
		}
		hashed := hashKey.HashKey()
		items[hashed] = key
	}
	return &object.Set{Items: items}
}
