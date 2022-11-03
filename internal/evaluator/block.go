package evaluator

import (
	"context"

	"github.com/myzie/tamarin/internal/ast"
	"github.com/myzie/tamarin/internal/object"
	"github.com/myzie/tamarin/internal/scope"
)

func (e *Evaluator) evalBlockStatement(
	ctx context.Context,
	block *ast.BlockStatement,
	s *scope.Scope,
) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = e.Evaluate(ctx, statement, s)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}
