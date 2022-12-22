package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalMapLiteral(ctx context.Context, node *ast.MapLiteral, s *scope.Scope) object.Object {
	items := make(map[string]object.Object, len(node.Pairs))
	for keyNode, valueNode := range node.Pairs {
		value := e.Evaluate(ctx, valueNode, s)
		if object.IsError(value) {
			return value
		}
		var key string
		if keyIdent, ok := keyNode.(*ast.Identifier); ok {
			// Key is an identifier (no quotes), e.g. { foo: 5 }
			key = keyIdent.String()
		} else {
			// Key is an expression, e.g. { "foo": 5 }
			keyObj := e.Evaluate(ctx, keyNode, s)
			if object.IsError(keyObj) {
				return keyObj
			}
			var err *object.Error
			key, err = object.AsString(keyObj)
			if err != nil {
				return err
			}
		}
		items[key] = value
	}
	return object.NewMap(items)
}
