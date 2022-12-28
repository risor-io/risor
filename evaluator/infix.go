package evaluator

import (
	"context"
	"math"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

func (e *Evaluator) evalInfixExpression(ctx context.Context, node *ast.Infix, s *scope.Scope) object.Object {
	left := e.Evaluate(ctx, node.Left(), s)
	if object.IsError(left) {
		return left
	}
	if node.Operator() == "&&" && !left.IsTruthy() {
		return left // Short circuit
	} else if node.Operator() == "||" && left.IsTruthy() {
		return left // Short circuit
	}
	right := e.Evaluate(ctx, node.Right(), s)
	if object.IsError(right) {
		return right
	}
	return e.evalInfix(node.Operator(), left, right, s)
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object, s *scope.Scope) object.Object {
	// Expressions that are handled the same for all types
	switch operator {
	case "==":
		return left.Equals(right)
	case "!=":
		return object.Not(left.Equals(right).(*object.Bool))
	case "&&":
		if left.IsTruthy() && right.IsTruthy() {
			return right
		}
		return right
	case "||":
		if left.IsTruthy() {
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
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, leftType, rightType)
	default:
		return object.Errorf("syntax error: invalid operation %s for types: %s and %s",
			operator, leftType, rightType)
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	l := object.NewString(string(left.Inspect()))
	r := object.NewString(string(right.Inspect()))
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
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Int).Value()
	rightVal := right.(*object.Int).Value()
	switch operator {
	case "+":
		return object.NewInt(leftVal + rightVal)
	case "+=":
		return object.NewInt(leftVal + rightVal)
	case "%":
		if rightVal == 0 {
			return object.Errorf("eval error: int modulo by zero")
		}
		return object.NewInt(leftVal % rightVal)
	case "**":
		return object.NewInt(int64(math.Pow(float64(leftVal), float64(rightVal))))
	case "-":
		return object.NewInt(leftVal - rightVal)
	case "-=":
		return object.NewInt(leftVal - rightVal)
	case "*":
		return object.NewInt(leftVal * rightVal)
	case "*=":
		return object.NewInt(leftVal * rightVal)
	case "/":
		if rightVal == 0 {
			return object.Errorf("eval error: int divided by zero")
		}
		return object.NewInt(leftVal / rightVal)
	case "/=":
		if rightVal == 0 {
			return object.Errorf("eval error: int divided by zero")
		}
		return object.NewInt(leftVal / rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value()
	rightVal := right.(*object.Float).Value()
	switch operator {
	case "+":
		return object.NewFloat(leftVal + rightVal)
	case "+=":
		return object.NewFloat(leftVal + rightVal)
	case "-":
		return object.NewFloat(leftVal - rightVal)
	case "-=":
		return object.NewFloat(leftVal - rightVal)
	case "*":
		return object.NewFloat(leftVal * rightVal)
	case "*=":
		return object.NewFloat(leftVal * rightVal)
	case "**":
		return object.NewFloat(math.Pow(leftVal, rightVal))
	case "/":
		if rightVal == 0 {
			return object.Errorf("eval error: float divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "/=":
		if rightVal == 0 {
			return object.Errorf("eval error: float divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalFloatIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value()
	rightVal := float64(right.(*object.Int).Value())
	switch operator {
	case "+":
		return object.NewFloat(leftVal + rightVal)
	case "+=":
		return object.NewFloat(leftVal + rightVal)
	case "-":
		return object.NewFloat(leftVal - rightVal)
	case "-=":
		return object.NewFloat(leftVal - rightVal)
	case "*":
		return object.NewFloat(leftVal * rightVal)
	case "*=":
		return object.NewFloat(leftVal * rightVal)
	case "**":
		return object.NewFloat(math.Pow(leftVal, rightVal))
	case "/":
		if rightVal == 0 {
			return object.Errorf("eval error: float divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "/=":
		if rightVal == 0 {
			return object.Errorf("eval error: float divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalIntegerFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := float64(left.(*object.Int).Value())
	rightVal := right.(*object.Float).Value()
	switch operator {
	case "+":
		return object.NewFloat(leftVal + rightVal)
	case "+=":
		return object.NewFloat(leftVal + rightVal)
	case "-":
		return object.NewFloat(leftVal - rightVal)
	case "-=":
		return object.NewFloat(leftVal - rightVal)
	case "*":
		return object.NewFloat(leftVal * rightVal)
	case "*=":
		return object.NewFloat(leftVal * rightVal)
	case "**":
		return object.NewFloat(math.Pow(leftVal, rightVal))
	case "/":
		if rightVal == 0 {
			return object.Errorf("eval error: int divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "/=":
		if rightVal == 0 {
			return object.Errorf("eval error: int divided by zero")
		}
		return object.NewFloat(leftVal / rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String).Value()
	r := right.(*object.String).Value()
	switch operator {
	case ">=":
		return nativeBoolToBooleanObject(l >= r)
	case ">":
		return nativeBoolToBooleanObject(l > r)
	case "<=":
		return nativeBoolToBooleanObject(l <= r)
	case "<":
		return nativeBoolToBooleanObject(l < r)
	case "+":
		return object.NewString(l + r)
	case "+=":
		return object.NewString(l + r)
	default:
		return object.Errorf("type error: unsupported operand types for %s: %s and %s",
			operator, left.Type(), right.Type())
	}
}
