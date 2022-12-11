package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalIndexExpression(
	ctx context.Context,
	node *ast.IndexExpression,
	s *scope.Scope,
) object.Object {
	left := e.Evaluate(ctx, node.Left, s)
	if isError(left) {
		return left
	}
	index := e.Evaluate(ctx, node.Index, s)
	if isError(index) {
		return index
	}
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.INTEGER_OBJ:
		return e.evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return e.evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		return e.evalStringIndexExpression(left, index)
	default:
		return newError("type error: %s object is not scriptable", left.Type())
	}
}

func (e *Evaluator) evalArrayIndexExpression(list, index object.Object) object.Object {
	listObject := list.(*object.List)
	idx := index.(*object.Integer).Value
	len := int64(len(listObject.Items))
	max := len - 1
	if idx > max {
		return newError("index error: array index out of range: %d", idx)
	}
	if idx >= 0 {
		return listObject.Items[idx]
	}
	// Handle negative indices, where -1 is the last item in the array
	reversed := idx + len
	if reversed < 0 || reversed > max {
		return newError("index error: array index out of range: %d", idx)
	}
	return listObject.Items[reversed]
}

func (e *Evaluator) evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, err := object.AsString(index)
	if err != nil {
		return err
	}
	value, ok := hashObject.Map[key]
	if !ok {
		return newError("key error: %v", index.Inspect())
	}
	return value
}

func (e *Evaluator) evalStringIndexExpression(input, index object.Object) object.Object {
	str := input.(*object.String).Value
	idx := index.(*object.Integer).Value
	len := int64(len(str))
	max := len - 1
	if idx > max {
		return newError("index error: string index out of range: %d", idx)
	}
	if idx >= 0 {
		chars := []rune(str)
		ret := chars[idx]
		return &object.String{Value: string(ret)}
	}
	// Handle negative indices, where -1 is the last rune in the string
	reversed := idx + len
	if reversed < 0 || reversed > max {
		return newError("index error: string index out of range: %d", idx)
	}
	chars := []rune(str)
	ret := chars[reversed]
	return &object.String{Value: string(ret)}
}
