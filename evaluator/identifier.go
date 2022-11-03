package evaluator

import (
	"github.com/myzie/tamarin/ast"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
)

func (e *Evaluator) evalIdentifier(node *ast.Identifier, s *scope.Scope) object.Object {
	if val, ok := s.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("name error: %s is not defined", node.Value)
}
