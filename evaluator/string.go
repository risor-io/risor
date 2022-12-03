package evaluator

import (
	"context"
	"strings"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalStringLiteral(ctx context.Context,
	node *ast.StringLiteral,
	s *scope.Scope,
) object.Object {
	if node.Template == nil {
		return &object.String{Value: node.Value}
	}
	var exprIndex int
	var parts []string
	for _, f := range node.Template.Fragments {
		switch f.IsVariable {
		case true:
			expr := node.TemplateExpressions[exprIndex]
			if expr == nil {
				parts = append(parts, "")
				continue
			}
			exprIndex++
			// Evaluate the variable
			obj := New(Opts{}).Evaluate(ctx, expr, s)
			if isError(obj) {
				return obj
			}
			parts = append(parts, obj.Inspect())
		case false:
			parts = append(parts, f.Value)
		}
	}
	return &object.String{Value: strings.Join(parts, "")}
}
