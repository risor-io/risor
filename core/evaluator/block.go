package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalBlockStatement(
	ctx context.Context,
	block *ast.Block,
	s *scope.Scope,
) object.Object {
	var result object.Object = object.Nil
	for _, statement := range block.Statements() {
		result = e.Evaluate(ctx, statement, s)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.Control:
				return result
			}
		}
	}
	return result
}
