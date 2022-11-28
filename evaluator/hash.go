package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

func (e *Evaluator) evalHashLiteral(
	ctx context.Context,
	node *ast.HashLiteral,
	s *scope.Scope,
) object.Object {
	hash := object.NewHash(make(map[string]interface{}, len(node.Pairs)))
	for keyNode, valueNode := range node.Pairs {
		key := e.Evaluate(ctx, keyNode, s)
		if isError(key) {
			return key
		}
		keyStr, err := object.AsString(key)
		if err != nil {
			return err
		}
		value := e.Evaluate(ctx, valueNode, s)
		if isError(value) {
			return value
		}
		hash.Set(keyStr, value)
	}
	return hash
}
