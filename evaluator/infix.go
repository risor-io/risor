package evaluator

import (
	"context"
	"math"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalInfixExpression(ctx context.Context, node *ast.InfixExpression, s *scope.Scope) object.Object {
	left := e.Evaluate(ctx, node.Left, s)
	if isError(left) {
		return left
	}
	if node.Operator == "&&" && !isTruthy(left) {
		return left // Short circuit
	} else if node.Operator == "||" && isTruthy(left) {
		return left // Short circuit
	}
	right := e.Evaluate(ctx, node.Right, s)
	if isError(right) {
		return right
	}
	return e.evalInfix(node.Operator, left, right, s)
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object, s *scope.Scope) object.Object {
	// Expressions that are handled the same for all types
	switch operator {
	case "==":
		return left.Equals(right)
	case "!=":
		return object.Not(left.Equals(right).(*object.Bool))
	case "&&":
		if isTruthy(left) && isTruthy(right) {
			return right
		}
		return right
	case "||":
		if isTruthy(left) {
			return left
		}
		return right
	}
	// Everything else
	leftType := left.Type()
	rightType := right.Type()
	switch {
	case leftType == object.INT && rightType == object.INT:
		return evalIntegerInfixExpression(operator, left, right)
	case leftType == object.FLOAT && rightType == object.FLOAT:
		return evalFloatInfixExpression(operator, left, right)
	case leftType == object.FLOAT && rightType == object.INT:
		return evalFloatIntegerInfixExpression(operator, left, right)
	case leftType == object.INT && rightType == object.FLOAT:
		return evalIntegerFloatInfixExpression(operator, left, right)
	case leftType == object.STRING && rightType == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case leftType == object.BOOL && rightType == object.BOOL:
		return evalBooleanInfixExpression(operator, left, right)
	case leftType != rightType:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, leftType, rightType)
	default:
		return newError("syntax error: invalid operation %s for types: %s and %s",
			operator, leftType, rightType)
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
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value
	switch operator {
	case "+":
		return &object.Int{Value: leftVal + rightVal}
	case "+=":
		return &object.Int{Value: leftVal + rightVal}
	case "%":
		return &object.Int{Value: leftVal % rightVal}
	case "**":
		return &object.Int{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "-":
		return &object.Int{Value: leftVal - rightVal}
	case "-=":
		return &object.Int{Value: leftVal - rightVal}
	case "*":
		return &object.Int{Value: leftVal * rightVal}
	case "*=":
		return &object.Int{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("zero division error")
		}
		return &object.Int{Value: leftVal / rightVal}
	case "/=":
		return &object.Int{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
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
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalFloatIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := float64(right.(*object.Int).Value)
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
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalIntegerFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := float64(left.(*object.Int).Value)
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
	default:
		return newError("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String)
	r := right.(*object.String)
	switch operator {
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
