package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/myzie/tamarin/internal/ast"
	"github.com/myzie/tamarin/internal/lexer"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x =5;", "x", 5},
		{"let z =1.3;", "z", 1.3},
		{"let y = true;", "y", true},
		{"let foobar=y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// Test that errors are returned when incomplete let/const expressions are seen
func TestBadLetConstStatement(t *testing.T) {
	input := []string{"let", "const", "let x;", "const x;"}

	for _, str := range input {
		l := lexer.New(str)
		p := New(l)
		_ = p.ParseProgram()

		errors := p.errors
		if len(errors) < 1 {
			t.Errorf("UNexpected error-count!")
		}

		if len(p.Errors()) != len(errors) {
			t.Errorf("Mismatch of errors + error-messages!")
		}
	}
}

// TestConstStatements tests the "const" token.
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
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

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got %q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	return true
}

func testConstStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "const" {
		t.Errorf("s.TokenLiteral not 'const'. got %q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.ConstStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
return 0b11;
return 0x15;
return 993322;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		fmt.Println(reflect.TypeOf(program.Statements[0]), program.Statements[0].TokenLiteral())
		t.Fatalf("program does not contain 3 statements, got=%d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statemtnets[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	integer, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if integer.Value != 5 {
		t.Errorf("integer.Value not %d. got=%d", 5, integer.Value)
	}
	if integer.TokenLiteral() != "5" {
		t.Errorf("integer.TokenLiteral not %s. got=%s", "5", integer.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	boolTests := []struct {
		input     string
		boolValue bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, tt := range boolTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statement. got=%d", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp is not ast.Boolean. got=%v", exp)
		}
		if exp.Value != tt.boolValue {
			t.Fatalf("exp.Value is not %t, got=%t", exp.Value, tt.boolValue)
		}
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statement doest not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

// skip float literal test
func testFloatLiteral(t *testing.T, exp ast.Expression, v float64) bool {
	float, ok := exp.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("exp not *ast.FloatLiteral. got=%T", exp)
		return false
	}
	if float.Value != v {
		t.Errorf("float.Value not %f. got=%f", v, float.Value)
		return false
	}
	return true
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t, got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case float32:
		return testFloatLiteral(t, exp, float64(v))
	case float64:
		return testFloatLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if(x<y){x}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statement was not nil. got=%+v", exp.Alternative)
	}
}

func TestForLoopExpression(t *testing.T) {
	input := `for(x<y) { let x=x+1; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.ForLoopExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForLoopExpression. got=%T",
			stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statement[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}
	if !testLetStatement(t, consequence, "x") {
		t.Fatalf("exp.Consequence is not LetStatement")
	}

}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x,y=3){x+y}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statement. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 Statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParsing(t *testing.T) {
	input := `func f(x,y){x+y;}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statement. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionDefineLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionDefineLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("stmt.Expression is not ast.FunctionDefineLiteral. got=%T",
			stmt.Expression)
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 Statements. got=%d\n",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParameer) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParameer), len(function.Parameters))
		}
		for i, ident := range tt.expectedParameer {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2*3, 4+5)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doest not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.errors
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiteral(t *testing.T) {
	input := `[1, 2*2, 3+3]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1+1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteral(t *testing.T) {
	input := `{"one":1, "two":2, "three":3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong legnth. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp isn not ast.HashLiteral. got=%T",
			stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralWithExpression(t *testing.T) {
	input := `{"one":0+1, "two":10 - 8, "three": 15/5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T",
			stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d",
			len(hash.Pairs))
	}
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
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

// Test operators: +=, -=, /=, and *=.
func TestMutators(t *testing.T) {
	input := []string{"let w = 5; w *= 3;",
		"let x = 15; x += 3;",
		"let y = 10; y /= 2;",
		"let z = 10; y -= 2;",
		"let z = 1; z++;",
		"let z = 1; z--;",
		"let z = 10; let a = 3; y = a;"}

	for _, txt := range input {
		l := lexer.New(txt)
		p := New(l)
		_ = p.ParseProgram()
		checkParserErrors(t, p)
	}
}

// Test method-call operation.
func TestObjectMethodCall(t *testing.T) {
	input := []string{"\"steve\".len()",
		"let x = 15; x.string();"}

	for _, txt := range input {
		l := lexer.New(txt)
		p := New(l)
		_ = p.ParseProgram()
		checkParserErrors(t, p)
	}
}

// Test that incomplete blocks / statements are handled.
func TestIncompleThings(t *testing.T) {
	input := []string{
		`if ( true ) { `,
		`if ( true ) { puts( "OK" ) ; } else { `,
		`let x = `,
		`const x =`,
		`func foo( a, b ="steve", `,
		`func foo() {`,
		`switch (foo) { `,
		`foo | bar`,
	}

	for _, str := range input {
		l := lexer.New(str)
		p := New(l)
		_ = p.ParseProgram()

		if len(p.errors) < 1 {
			t.Errorf("unexpected error-count, got %d  expected %d", len(p.errors), 1)
		}

		if !strings.Contains(p.errors[0], "unterminated") {
			t.Errorf("Unexpected error-message %s\n", p.errors[0])
		}
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
	parser := New(lexer.New(input))
	p := parser.ParseProgram()
	require.Len(t, parser.errors, 0)
	require.Len(t, p.Statements, 1)
	expr, ok := p.Statements[0].(*ast.ExpressionStatement)
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
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	if len(p.errors) < 1 {
		t.Errorf("unexpected error-count, got %d expected %d", len(p.errors), 1)
	}
	if !strings.Contains(p.errors[0], "only have one default block") {
		t.Errorf("Unexpected error-message %s\n", p.errors[0])
	}
}

func TestForLoop(t *testing.T) {
	fmt.Println("FOR")
	tests := []struct {
		input   string
		initStr string
		condStr string
		postStr string
	}{
		{
			"for (let i = 0; i < 5; i++) { }",
			"let i = 0;",
			"(i < 5)",
			"(i++)",
		},
		{
			"for (i < 5) { }",
			"",
			"(i < 5)",
			"",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		s := program.Statements[0]
		exprStatement := s.(*ast.ExpressionStatement)
		expr, ok := exprStatement.Expression.(*ast.ForLoopExpression)
		if !ok {
			t.Fatalf("Expected a for loop expression; got=%v", s)
		}
		fmt.Println(expr)
		if expr.Condition.String() != tt.condStr {
			t.Fatalf("incorrect condition. got='%v' want='%v'", expr.Condition.String(), tt.condStr)
		}
		if tt.initStr != "" {
			if expr.InitStatement.String() != tt.initStr {
				t.Fatalf("incorrect condition. got='%v' want='%v'", expr.InitStatement.String(), tt.initStr)
			}
		} else {
			if expr.InitStatement != nil {
				t.Fatalf("expected no init statement. got='%v'", expr.InitStatement.String())
			}
		}
		if tt.postStr != "" {
			if expr.PostStatement.String() != tt.postStr {
				t.Fatalf("incorrect condition. got='%v' want='%v'", expr.PostStatement.String(), tt.postStr)
			}
		} else {
			if expr.PostStatement != nil {
				t.Fatalf("expected no post statement. got='%v'", expr.PostStatement.String())
			}
		}
	}
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
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
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
	test := `{
		"a": "b",

		"c": "d",

	}
	`
	l := lexer.New(test)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	fmt.Println(program)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	hash, ok := expr.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 2)
}

func TestSetExpression(t *testing.T) {
	test := `{
		"a",
		1, 2,
		"c",
	}
	`
	l := lexer.New(test)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	fmt.Println(program)
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
	test := `foo(
		a=1,
		b=2,
	)
	`
	l := lexer.New(test)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	fmt.Println(program)
	require.Len(t, program.Statements, 1)
	stmt := program.Statements[0]
	expr, ok := stmt.(*ast.ExpressionStatement)
	require.True(t, ok)
	call, ok := expr.Expression.(*ast.CallExpression)
	require.True(t, ok)
	require.Equal(t, "foo", call.Function.String())
	require.Len(t, call.Arguments, 2)
	arg0 := call.Arguments[0].(*ast.AssignStatement)
	require.Equal(t, "a=1", arg0.String())
	arg1 := call.Arguments[1].(*ast.AssignStatement)
	require.Equal(t, "b=2", arg1.String())
}

func TestGetAttribute(t *testing.T) {
	test := "foo.bar"
	l := lexer.New(test)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
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
