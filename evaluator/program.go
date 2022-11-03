package evaluator

import (
	"context"

	"github.com/myzie/tamarin/ast"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
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
