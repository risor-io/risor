package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
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
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.BreakValue:
				return result
			}
		}
	}
	return result
}
