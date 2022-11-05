package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

// evalTernaryExpression handles a ternary-expression. If the condition
// is true we return the contents of evaluating the true-branch, otherwise
// the false-branch.
func (e *Evaluator) evalTernaryExpression(
	ctx context.Context,
	te *ast.TernaryExpression,
	s *scope.Scope,
) object.Object {
	condition := e.Evaluate(ctx, te.Condition, s)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return e.Evaluate(ctx, te.IfTrue, s)
	}
	return e.Evaluate(ctx, te.IfFalse, s)
}
