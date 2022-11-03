package evaluator

import (
	"context"
	"math"

	"github.com/myzie/tamarin/internal/ast"
	"github.com/myzie/tamarin/internal/object"
	"github.com/myzie/tamarin/internal/scope"
)

func (e *Evaluator) evalInfixExpression(ctx context.Context, node *ast.InfixExpression, s *scope.Scope) object.Object {
	left := e.Evaluate(ctx, node.Left, s)
	if isError(left) {
		return left
	}
	right := e.Evaluate(ctx, node.Right, s)
	if isError(right) {
		return right
	}
	return e.evalInfix(node.Operator, left, right, s)
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object, s *scope.Scope) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatIntegerInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalIntegerFloatInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "&&":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) && objectToNativeBoolean(right))
	case operator == "||":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) || objectToNativeBoolean(right))
	case operator == "!~":
		return notMatches(left, right)
	case operator == "~=":
		return matches(left, right, s)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	default:
		return newError("syntax error: invalid operation %s for types: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	l := &object.String{Value: string(left.Inspect())}
	r := &object.String{Value: string(right.Inspect())}
	switch operator {
	case "<":
		return evalStringInfixExpression(operator, l, r)
	case "<=":
		return evalStringInfixExpression(operator, l, r)
	case ">":
		return evalStringInfixExpression(operator, l, r)
	case ">=":
		return evalStringInfixExpression(operator, l, r)
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "**":
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("zero division error")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "/=":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "..":
		if rightVal > leftVal {
			len := int(rightVal-leftVal) + 1
			if len > 1000 {
				return newError("exceeded max length: %d", len)
			}
			array := make([]object.Object, len)
			for i := 0; i < len; i++ {
				array[i] = &object.Integer{Value: leftVal}
				leftVal++
			}
			return &object.Array{Elements: array}
		}
		len := int(leftVal-rightVal) + 1
		if len > 1000 {
			return newError("exceeded max length: %d", len)
		}
		array := make([]object.Object, len)
		for i := 0; i < len; i++ {
			array[i] = &object.Integer{Value: leftVal}
			leftVal--
		}
		return &object.Array{Elements: array}
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalFloatIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := float64(right.(*object.Integer).Value)
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalIntegerFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := float64(left.(*object.Integer).Value)
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "+=":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "-=":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "*=":
		return &object.Float{Value: leftVal * rightVal}
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "/=":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String)
	r := right.(*object.String)
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	case ">=":
		return nativeBoolToBooleanObject(l.Value >= r.Value)
	case ">":
		return nativeBoolToBooleanObject(l.Value > r.Value)
	case "<=":
		return nativeBoolToBooleanObject(l.Value <= r.Value)
	case "<":
		return nativeBoolToBooleanObject(l.Value < r.Value)
	case "+":
		return &object.String{Value: l.Value + r.Value}
	case "+=":
		return &object.String{Value: l.Value + r.Value}
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}
