package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalIndexExpression(ctx context.Context, node *ast.IndexExpression, s *scope.Scope) object.Object {
	left := e.Evaluate(ctx, node.Left, s)
	if object.IsError(left) {
		return left
	}
	container, ok := left.(object.Container)
	if !ok {
		return object.Errorf("type error: %s object is not scriptable", left.Type())
	}
	// Retrieve an item with a single index
	if node.Index != nil {
		index := e.Evaluate(ctx, node.Index, s)
		if object.IsError(index) {
			return index
		}
		item, err := container.GetItem(index)
		if err != nil {
			return err
		}
		return item
	}
	// Retrieve a slice of items with a range of indices
	var startIndex, stopIndex object.Object
	if node.FromIndex != nil {
		startIndex = e.Evaluate(ctx, node.FromIndex, s)
		if object.IsError(startIndex) {
			return startIndex
		}
	}
	if node.ToIndex != nil {
		stopIndex = e.Evaluate(ctx, node.ToIndex, s)
		if object.IsError(stopIndex) {
			return stopIndex
		}
	}
	items, err := container.GetSlice(object.Slice{Start: startIndex, Stop: stopIndex})
	if err != nil {
		return err
	}
	return items
}
