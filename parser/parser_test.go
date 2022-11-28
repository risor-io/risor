package parser

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cloudcmds/tamarin/ast"
	"github.com/cloudcmds/tamarin/token"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value interface{}
	}{
		{"let x =5;", "x", 5},
		{"let z =1.3;", "z", 1.3},
		{"let y = true;", "y", true},
		{"let foobar=y;", "foobar", "y"},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.ident)
		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.value)
	}
}

func TestDeclareStatements(t *testing.T) {
	input := `
	let x = foo.bar()
	y := foo.bar()
	`
	program, err := Parse(input)
	printMultiError(err)
	require.Nil(t, err)
	require.Len(t, program.Statements, 2)
	stmt1, ok := program.Statements[0].(*ast.LetStatement)
	require.True(t, ok)
	stmt2, ok := program.Statements[1].(*ast.LetStatement)
	require.True(t, ok)
	fmt.Println(stmt1)
	fmt.Println(stmt2)
	// require.Equal(t, let.TokenLiteral(), "x")
}

func TestBadLetConstStatement(t *testing.T) {
	inputs := []struct {
		input string
		err   string
	}{
		{"let", "parse error: unexpected end of file while parsing let statement (expected identifier)"},
		{"const", "parse error: unexpected end of file while parsing const statement (expected identifier)"},
		{"let x;", "parse error: unexpected ; while parsing let statement (expected =)"},
		{"const x;", "parse error: unexpected ; while parsing const statement (expected =)"},
	}
	for _, tt := range inputs {
		_, err := Parse(tt.input)
		require.NotNil(t, err)
		e, ok := err.(ParserError)
		require.True(t, ok)
		require.Equal(t, tt.err, e.Error())
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"const x =5;", "x", 5},
		{"const z =1.3;", "z", 1.3},
		{"const y = true;", "y", true},
		{"const foobar=y;", "foobar", "y"},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt := program.Statements[0]
		if !testConstStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.ConstStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 0b11;
return 0x15;
return 993322;
`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 3)
	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		require.True(t, ok)
		require.Equal(t, returnStatement.TokenLiteral(), "return")
	}
}

func TestIdentifierExpression(t *testing.T) {
	program, err := Parse("foobar;")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	ident, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, ident.Value, "foobar")
	require.Equal(t, ident.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"0", 0},
		{"5", 5},
		{"10", 10},
		{"9876543210", 9876543210},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		integer, ok := stmt.Expression.(*ast.IntegerLiteral)
		require.True(t, ok, "got %T", stmt.Expression)
		require.Equal(t, integer.Value, tt.value)
		require.Equal(t, integer.TokenLiteral(), fmt.Sprintf("%d", tt.value))
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input     string
		boolValue bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		exp, ok := stmt.Expression.(*ast.Boolean)
		require.True(t, ok)
		require.Equal(t, exp.Value, tt.boolValue)
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false", "!", false},
	}
	for _, tt := range prefixTests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok)
		require.Equal(t, exp.Operator, tt.operator)
		testLiteralExpression(t, exp.Right, tt.integerValue)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"0.4+1.3", 0.4, "+", 1.3},
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5<5;", 5, "<", 5},
		{"2**3;", 2, "**", 3},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true!=false", true, "!=", false},
		{"false==false", false, "==", false},
	}
	for _, tt := range infixTests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		{"a+b-c", "((a + b) - c)"},
		{"a*b*c", "((a * b) * c)"},
		{"a*b/c", "((a * b) / c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3+4;-5*5", "(3 + 4)((-5) * 5)"},
		{"5>4==3<4", "((5 > 4) == (3 < 4))"},
		{"5<4!=3>4", "((5 < 4) != (3 > 4))"},
		{"3+4*5==3*1+4*5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3>5==false", "((3 > 5) == false)"},
		{"3<5==true", "((3 < 5) == true)"},
		{"1+(2+3)+4", "((1 + (2 + 3)) + 4)"},
		{"(5+5)*2", "((5 + 5) * 2)"},
		{"2/(5+5)", "(2 / (5 + 5))"},
		{"2**3", "(2 ** 3)"},
		{"-(5+5)", "(-(5 + 5))"},
		{"!(true==true)", "(!(true == true))"},
		{"a + add(b*c)+d", "((a + add((b * c))) + d)"},
		{"a*[1,2,3,4][b*c]*d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a*b[2], b[1], 2 * [1,2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		actual := program.String()
		require.Equal(t, tt.expected, actual)
	}
}

func TestIfExpression(t *testing.T) {
	program, err := Parse("if x < y { x }")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok)
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	require.Len(t, exp.Consequence.Statements, 1)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	require.Nil(t, exp.Alternative)
}

func TestFunctionLiteralParsing(t *testing.T) {
	program, err := Parse("func(x, y=3) { x + y }")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok)
	require.Len(t, function.Parameters, 2)
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	require.Len(t, function.Body.Statements, 1)
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParsing(t *testing.T) {
	program, err := Parse("func f(x, y) { x + y; }")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	function, ok := stmt.Expression.(*ast.FunctionDefineLiteral)
	require.True(t, ok)
	require.Len(t, function.Parameters, 2)
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	require.Len(t, function.Body.Statements, 1)
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input            string
		expectedParameer []string
	}{
		{"func(){}", []string{}},
		{"func(x){}", []string{"x"}},
		{"func(x,y){}", []string{"x", "y"}},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		require.Len(t, function.Parameters, len(tt.expectedParameer))
		for i, ident := range tt.expectedParameer {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	program, err := Parse("add(1, 2*3, 4+5)")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t, ok)
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	require.Len(t, exp.Arguments, 3)
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	program, err := Parse(`"hello world";`)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok)
	require.Equal(t, "hello world", literal.Value)
}

func TestParsingArrayLiteral(t *testing.T) {
	program, err := Parse("[1, 2*2, 3+3]")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	require.True(t, ok)
	require.Len(t, array.Elements, 3)
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1+1]"
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	require.True(t, ok)
	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestParsingHashLiteral(t *testing.T) {
	input := `{"one":1, "two":2, "three":3}`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 3)
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		require.True(t, ok)
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 0)
}

func TestParsingHashLiteralWithExpression(t *testing.T) {
	input := `{"one":0+1, "two":10 - 8, "three": 15/5}`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 3)
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		require.True(t, ok)
		testFunc, ok := tests[literal.String()]
		require.True(t, ok)
		testFunc(value)
	}
}

// Test operators: +=, -=, /=, and *=.
func TestMutators(t *testing.T) {
	inputs := []string{
		"let w = 5; w *= 3;",
		"let x = 15; x += 3;",
		"let y = 10; y /= 2;",
		"let z = 10; y -= 2;",
		"let z = 1; z++;",
		"let z = 1; z--;",
		"let z = 10; let a = 3; y = a;",
	}
	for _, input := range inputs {
		_, err := Parse(input)
		require.Nil(t, err)
	}
}

// Test method-call operation.
func TestObjectMethodCall(t *testing.T) {
	inputs := []string{
		"\"steve\".len()",
		"let x = 15; x.string();",
	}
	for _, input := range inputs {
		_, err := Parse(input)
		require.Nil(t, err)
	}
}

// Test that incomplete blocks / statements are handled.
func TestIncompleThings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`if ( true ) { `, "parse error: unterminated block statement"},
		{`if ( true ) { puts( "OK" ) ; } else { `, "parse error: unterminated block statement"},
		{`let x = `, "parse error: assignment is missing a value"},
		{`const x =`, "parse error: assignment is missing a value"},
		{`func foo( a, b ="steve", `, "parse error: unterminated function parameters"},
		{`func foo() {`, "parse error: unterminated block statement"},
		{`switch (foo) { `, "parse error: unterminated switch statement"},
	}
	for _, tt := range tests {
		_, err := Parse(tt.input)
		require.NotNil(t, err)
		pe, ok := err.(ParserError)
		require.True(t, ok)
		require.Equal(t, tt.expected, pe.Error())
	}
}

func TestSwitch(t *testing.T) {
	input := `switch val {
   case 1:
      x
   default:
      y
	  x = x + 1
}`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	switchExpr, ok := expr.Expression.(*ast.SwitchExpression)
	require.True(t, ok)
	require.Equal(t, "val", switchExpr.Value.String())
	require.Len(t, switchExpr.Choices, 2)
	choice1 := switchExpr.Choices[0]
	require.Len(t, choice1.Expr, 1)
	require.Equal(t, "1", choice1.Expr[0].String())
	choice2 := switchExpr.Choices[1]
	require.Len(t, choice2.Expr, 0)
}

func TestMultiDefault(t *testing.T) {
	input := `
switch val {
case 1:
    print("1")
case 2:
    print("2")
default:
    print("default")
default:
    print("oh no!")
}`
	_, err := Parse(input)
	require.NotNil(t, err)
	parserErr, ok := err.(ParserError)
	require.True(t, ok)
	require.Equal(t, "parse error: switch statement has multiple default blocks", parserErr.Error())
	require.Equal(t, 0, parserErr.StartPosition().Column)
	require.Equal(t, 8, parserErr.StartPosition().Line)
	require.Equal(t, 6, parserErr.EndPosition().Column) // last col in the word "default"
	require.Equal(t, 8, parserErr.EndPosition().Line)
}

func TestPipeExpression(t *testing.T) {
	tests := []struct {
		input          string
		exprType       string
		expectedIdents []string
	}{
		{"let x = foo | bar;", "ident", []string{"foo", "bar"}},
		{`let x = foo() | bar(name="foo") | baz(y=4);`, "call", []string{"foo", "bar", "baz"}},
		{`let x = a() | b();`, "call", []string{"a", "b"}},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		stmt := program.Statements[0].(*ast.LetStatement)
		expr, ok := stmt.Value.(*ast.PipeExpression)
		require.True(t, ok)
		require.Len(t, expr.Arguments, len(tt.expectedIdents))
		if tt.exprType == "ident" {
			for i, ident := range tt.expectedIdents {
				identExpr, ok := expr.Arguments[i].(*ast.Identifier)
				require.True(t, ok)
				require.Equal(t, ident, identExpr.Value)
			}
		} else if tt.exprType == "call" {
			for i, ident := range tt.expectedIdents {
				callExpr, ok := expr.Arguments[i].(*ast.CallExpression)
				require.True(t, ok)
				require.Equal(t, ident, callExpr.Function.String())
			}
		}
	}
}

func TestHashExpression(t *testing.T) {
	input := `{
		"a": "b",

		"c": "d",

	}
	`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	hash, ok := expr.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 2)
}

func TestSetExpression(t *testing.T) {
	input := `{
		"a",
		1, 2,
		"c",
	}
	`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	set, ok := expr.Expression.(*ast.SetLiteral)
	require.True(t, ok)
	require.Len(t, set.Items, 4)
	require.Equal(t, `{a, 1, 2, c}`, set.String())
}

func TestCallExpression(t *testing.T) {
	input := `foo(
		a=1,
		b=2,
	)
	`
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	call, ok := expr.Expression.(*ast.CallExpression)
	require.True(t, ok)
	require.Equal(t, "foo", call.Function.String())
	require.Len(t, call.Arguments, 2)
	arg0 := call.Arguments[0].(*ast.AssignStatement)
	require.Equal(t, "a = 1", arg0.String())
	arg1 := call.Arguments[1].(*ast.AssignStatement)
	require.Equal(t, "b = 2", arg1.String())
}

func TestGetAttribute(t *testing.T) {
	program, err := Parse("foo.bar")
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	getAttr, ok := expr.Expression.(*ast.GetAttributeExpression)
	require.True(t, ok)
	require.Equal(t, "bar", getAttr.Attribute.String())
	getAttrStr := getAttr.String()
	require.Equal(t, "foo.bar", getAttrStr)
}

func TestForLoop(t *testing.T) {
	tests := []struct {
		input   string
		initStr string
		condStr string
		postStr string
	}{
		{
			"for let i = 0; i < 5; i++ { }",
			"let i = 0",
			"(i < 5)",
			"(i++)",
		},
		{
			"for i := 2+2; x < i; x-- { }",
			"i := (2 + 2)",
			"(x < i)",
			"(x--)",
		},
	}
	for _, tt := range tests {
		program, err := Parse(tt.input)
		printMultiError(err)
		require.Nil(t, err)
		require.Len(t, program.Statements, 1)
		s := program.Statements[0]
		exprStatement := s.(*ast.ExpressionStatement)
		expr, ok := exprStatement.Expression.(*ast.ForLoopExpression)
		require.True(t, ok)
		require.Equal(t, tt.condStr, expr.Condition.String())
		if tt.initStr != "" {
			require.Equal(t, tt.initStr, expr.InitStatement.String())
		} else {
			if expr.InitStatement != nil {
				t.Fatalf("expected no init statement. got='%v'", expr.InitStatement.String())
			}
		}
		if tt.postStr != "" {
			require.Equal(t, tt.postStr, expr.PostStatement.String())
		} else {
			if expr.PostStatement != nil {
				t.Fatalf("expected no post statement. got='%v'", expr.PostStatement.String())
			}
		}
	}
}

func TestBreak(t *testing.T) {
	input := "break"
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	_, ok := stmt.(*ast.BreakStatement)
	require.True(t, ok)
}

func TestBacktick(t *testing.T) {
	input := "`" + `\\n\t foo bar /hey there/` + "`"
	program, err := Parse(input)
	require.Nil(t, err)
	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	str, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok)
	require.Equal(t, `\\n\t foo bar /hey there/`, str.Value)
}

func TestUnterminatedBacktickString(t *testing.T) {
	input := "`foo"
	_, err := Parse(input)
	require.NotNil(t, err)
	require.Equal(t, "syntax error: unterminated string literal", err.Error())
	var syntaxErr *SyntaxError
	ok := errors.As(err, &syntaxErr)
	require.True(t, ok)
	require.NotNil(t, syntaxErr.Cause())
	require.Equal(t, "unterminated string literal", syntaxErr.Cause().Error())
	require.Equal(t, NewSyntaxError(ErrorOpts{
		ErrType: "syntax error",
		Cause:   syntaxErr.Cause(),
		StartPosition: token.Position{
			Value: rune('`'),
		},
		EndPosition: token.Position{
			Value:  rune('o'),
			Column: 3, // the last "o" in foo is at index 3
			Char:   3, // the last "o" in foo is at index 3
		},
		File:       "",
		SourceCode: "`foo",
	}), syntaxErr)
}

func TestUnterminatedString(t *testing.T) {
	input := `42
x := "a`
	ctx := context.Background()
	_, err := ParseWithOpts(ctx, Opts{Input: input, File: "main.tm"})
	require.NotNil(t, err)
	fmt.Printf("%+v\n", err.(*SyntaxError).StartPosition())
	require.Equal(t, "syntax error: unterminated string literal", err.Error())
	var syntaxErr *SyntaxError
	ok := errors.As(err, &syntaxErr)
	require.True(t, ok)
	require.NotNil(t, syntaxErr.Cause())
	require.Equal(t, "unterminated string literal", syntaxErr.Cause().Error())
	require.Equal(t, NewSyntaxError(ErrorOpts{
		ErrType: "syntax error",
		Cause:   syntaxErr.Cause(),
		StartPosition: token.Position{
			Value:     rune('"'),
			Column:    5,
			Line:      1,
			LineStart: 3,
			Char:      8,
		},
		EndPosition: token.Position{
			Value:     rune('a'),
			Column:    6,
			Line:      1,
			LineStart: 3,
			Char:      9,
		},
		File:       "main.tm",
		SourceCode: "x := \"a",
	}), syntaxErr)
}
