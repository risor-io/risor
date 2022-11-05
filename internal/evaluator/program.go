package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

func (e *Evaluator) evalProgram(
	ctx context.Context,
	program *ast.Program,
	s *scope.Scope,
) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = e.Evaluate(ctx, statement, s)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}
