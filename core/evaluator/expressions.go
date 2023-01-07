package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalExpressions(
	ctx context.Context,
	exps []ast.Expression,
	s *scope.Scope,
) []object.Object {
	values := make([]object.Object, len(exps))
	for i, exp := range exps {
		value := e.Evaluate(ctx, exp, s)
		if object.IsError(value) {
			return []object.Object{value}
		}
		values[i] = value
	}
	return values
}

func (e *Evaluator) evalExpressionsIgnoreErrors(
	ctx context.Context,
	exps []ast.Expression,
	s *scope.Scope,
) []object.Object {
	values := make([]object.Object, len(exps))
	for i, exp := range exps {
		values[i] = e.Evaluate(ctx, exp, s)
	}
	return values
}
