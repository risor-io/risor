package evaluator

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/ast"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
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
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return e.evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return e.evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		return e.evalStringIndexExpression(left, index)
	default:
		return newError("type error: %s object is not scriptable", left.Type())
	}
}

func (e *Evaluator) evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	len := int64(len(arrayObject.Elements))
	max := len - 1
	if idx > max {
		return newError("index error: array index out of range: %d", idx)
	}
	if idx >= 0 {
		return arrayObject.Elements[idx]
	}
	// Handle negative indices, where -1 is the last item in the array
	reversed := idx + len
	if reversed < 0 || reversed > max {
		return newError("index error: array index out of range: %d", idx)
	}
	return arrayObject.Elements[reversed]
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
