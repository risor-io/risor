package evaluator

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/skx/monkey/lexer"
	"github.com/skx/monkey/object"
	"github.com/skx/monkey/parser"
)

func TestEvalArithmeticExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5+5+5+5-10", 10},
		{"2*2*2*2*2", 32},
		{"-50+100+ -50", 0},
		{"5*2+10", 20},
		{"5+2*10", 25},
		{"20 + 2 * -10", 0},
		{"50/2 * 2 +10", 60},
		{"2*(5+10)", 30},
		{"3*3*3+10", 37},
		{"3*(3*3)+10", 37},
		{"(5+10*2+15/3)*2+-10", 50},
		{"1.2", 1.2},
		{"-2.3", -2.3},
		{"1.2+3.2", 4.4},
		{"1+2.3", 3.3},
		{"2.3*1.0", 2.3},
		{"3.2-5.8", -2.6},
		{"2**3", 8},
		{"2.0**3", 8},
		{"2**3.0", 8},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testDecimalObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testDecimalObject(t *testing.T, obj object.Object, expected interface{}) bool {
	switch exp := expected.(type) {
	case int64:
		return testIntegerObject(t, obj, exp)
	case float64:
		return testFloatObject(t, obj, exp)
	default:
		return false
	}
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not Integer. got=%T(%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("obj is not Float. got=%T(%+v)", obj, obj)
		return false
	}
	if math.Abs(result.Value-expected) > 0.00001 {
		t.Errorf("object has wrong value. got=%f, want=%f",
			result.Value, expected)
		return false
	}
	return true
}
func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("obj is not String. got=%T(%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1<1", false},
		{"1>1", false},
		{"1==1", true},
		{"\"a\">=\"A\"", true},
		{"\"a\"<=\"A\"", false},
		{"\"steve\"==\"steve\"", true},
		{"\"steve\"!=\"Steve\"", true},
		{"\"steve\"==\"kemp\"", false},
		{"\"abc123\"==\"abc\" + \"123\"", true},
		{"1!=1", false},
		{"1==2", false},
		{"1.0==1", true},
		{"1.5==1", false},
		{"1!=2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(1<2)==true", true},
		{"(1<2) == false", false},
		{"(1>2) == true", false},
		{"(1>2)==false", true},
		{"(1>=1)==true", true},
		{"(2<=2)==true", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)

	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean. got=%T(%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1<2) {10}", 10},
		{"if (1<2) { 10} else {20}", 10},
		{"if (1>2) {10} else {20}", 20},
		{"if (1>=1) {10} else {100}", 10},
		{"if (1<=1) {10} else {100}", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDecimalObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T(%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2*5;9;", 10},
		{"9; return 2*5; 9;", 10},
		{`if (10>1) { if (10>1) { return 10;} return 1;}`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testDecimalObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5+true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5+true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"3--", "3 is unknown"},
		{"true+false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5;true+false;5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10>1) { true+false;}", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if (10 > 1) {
      if (10>1) {
			return true+false;
			}
			return 1;
}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let a=5;a;", 5},
		{"let a=5*5; a;", 25},
		{"let a=5; let b=a; b;", 5},
		{"let a=5; a--; a;", 4},
		{"let a=5; a++; a;", 6},
		{"let a=5; let b=a; let c=a+b+5; c;", 15},
	}
	for _, tt := range tests {
		testDecimalObject(t, testEval(tt.input), tt.expect)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x+2; };`
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T(%+v)",
			evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := `(x + 2)`
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity=fn(x){x;}; identity(5);", 5},
		{"let identity=fn(x){return x;}; identity(5);", 5},
		{"let double=fn(x){x*2;}; double(5);", 10},
		{"let add = fn(x, y) { x+y;}; add(5,5);", 10},
		{"let add=fn(x,y){x+y;}; add(5+5, add(5,5));", 20},
		{"fn(x){x;}(5)", 5},
	}
	for _, tt := range tests {
		testDecimalObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x+y };
};
let addTwo = newAdder(2);
addTwo(2);
`
	testDecimalObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T(%+v)",
			evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("狐犬")`, 2},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got=INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testDecimalObject(t, evaluated, int64(expected))
		case string:
			if evaluated == NULL {
				t.Errorf("Got NULL output on input of '%s'\n", tt.input)
			} else {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error, got=%T(%+v)",
						evaluated, evaluated)
				}
				if errObj.Message != expected {
					t.Errorf("wrong err messsage. expected=%q, got=%q",
						expected, errObj.Message)
				}
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2*2, 3+3]`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T(%v)",
			evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testDecimalObject(t, result.Elements[0], 1)
	testDecimalObject(t, result.Elements[1], 4)
	testDecimalObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1,2,3][0]",
			1,
		},
		{
			"[1,2,3][1]",
			2,
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"let i =0; [1][i]",
			1,
		},
		{
			"let myArray=[1,2,3];myArray[2];",
			3,
		},
		{
			"let myArray=[1,2,3];myArray[0]+myArray[1]+myArray[2]",
			6,
		},
		{
			"let myArray=[1,2,3];let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1,2,3][3]",
			nil,
		},
		{
			"[1,2,3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDecimalObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"\"Steve\"[0]",
			"S",
		},
		{
			"\"Steve\"[1]",
			"t",
		},
		{
			"\"Steve\"[101]",
			nil,
		},
		{
			"\"Steve\"[-1]",
			nil,
		},
		{
			"\"狐犬\"[0]",
			"狐",
		},
		{
			"\"狐犬\"[1]",
			"犬",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		str, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, str)
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func TestHashLiterals(t *testing.T) {
	input := `let two="two";
	{
		"one":10-9,
		two:1+1,
		"thr" + "ee" : 6/2,
		4 : 4,
		true:5,
		false:6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval did't return Hash. got=%T(%+v)",
			evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testDecimalObject(t, pair.Value, expectedValue)
	}

}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo":5}["foo"]`,
			5,
		},
		{
			`{"foo":5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo":5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5:5}[5]`,
			5,
		},
		{
			`{true:5}[true]`,
			5,
		},
		{
			`{false:5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDecimalObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForLoopExpression(t *testing.T) {
	input := `
let x = 1;
let sum = 0;
let up = 100;
for (x < up){
	let sum = sum + x;
	let x = x + 1;
}
sum
`
	evaluated := testEval(input)
	testDecimalObject(t, evaluated, 4950)
}

func TestTypeBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"type( \"Steve\" );",
			"string",
		},
		{
			"type( 1 );",
			"integer",
		},
		{
			"type( 3.14159 );",
			"float",
		},
		{
			"type( [1,2,3] );",
			"array",
		},
		{
			"type( { \"name\":\"monkey\", true:1, 7:\"sevent\"} );",
			"hash",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		str, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, str)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestTimeout(t *testing.T) {
	input := `
i = 1;
for ( true ) {
  i++;
}
`
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	SetContext(ctx)
	evaluated := Eval(program, env)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Errorf("no error object returned. got=%T(%+v)",
			evaluated, evaluated)
	}

	if !strings.Contains(errObj.Message, "deadline") {
		t.Errorf("got error, but wasn't timeout: %s", errObj.Message)
	}

}
