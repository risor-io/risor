package evaluator

import (
	"github.com/cloudcmds/tamarin/core/ast"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
)

func (e *Evaluator) evalIdentifier(node *ast.Ident, s *scope.Scope) object.Object {
	name := node.String()
	if val, ok := s.Get(name); ok {
		return val
	}
	if builtin, ok := e.builtins[name]; ok {
		return builtin
	}
	return object.Errorf("name error: %q is not defined", name)
}
