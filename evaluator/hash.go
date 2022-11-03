package evaluator

import (
	"context"

	"github.com/myzie/tamarin/ast"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
)

func (e *Evaluator) evalHashLiteral(
	ctx context.Context,
	node *ast.HashLiteral,
	s *scope.Scope,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := e.Evaluate(ctx, keyNode, s)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("type error: %s object is unhashable", key.Type())
		}
		value := e.Evaluate(ctx, valueNode, s)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}
