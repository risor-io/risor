package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalProgram(
	ctx context.Context,
	program *ast.Program,
	s *scope.Scope,
) object.Object {
	var result object.Object
	for _, statement := range program.Statements() {
		result = e.Evaluate(ctx, statement, s)
		switch result := result.(type) {
		case *object.Control:
			switch result.Keyword() {
			case "break":
				return object.Errorf("eval error: break statement outside loop")
			case "continue":
				return object.Errorf("eval error: continue statement outside loop")
			case "return":
				return result.Value()
			}
		case *object.Error:
			return result
		}
	}
	return result
}
