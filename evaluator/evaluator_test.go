package evaluator

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	modStrings "github.com/cloudcmds/tamarin/modules/strings"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/stretchr/testify/require"
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
		{"2.0**3", 8.0},
		{"2**3.0", 8.0},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testDecimalObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	program, _ := parser.Parse(input)
	e := New(Opts{})
	return e.Evaluate(context.Background(), program, scope.New(scope.Opts{}))
}

func testDecimalObject(t *testing.T, obj object.Object, expected interface{}) bool {
	t.Helper()
	switch exp := expected.(type) {
	case int:
		return testIntegerObject(t, obj, int64(exp))
	case int64:
		return testIntegerObject(t, obj, exp)
	case float64:
		return testFloatObject(t, obj, exp)
	default:
		return false
	}
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	t.Helper()
	result, ok := obj.(*object.Int)
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
	t.Helper()
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
	t.Helper()
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
	t.Helper()
	result, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("object is not bool. got=%T(%+v)", obj, obj)
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
		{"!!true", true},
		{"!!false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
	result := testEval("!5")
	resultErr, ok := result.(*object.Error)
	require.True(t, ok)
	require.Equal(t, "type error: expected boolean to follow ! operator (got int)", resultErr.Message)
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
			testNilObject(t, evaluated)
		}
	}
}

func testNilObject(t *testing.T, obj object.Object) bool {
	t.Helper()
	if obj != object.Nil {
		t.Errorf("object is not nil. got=%T(%+v)", obj, obj)
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
		{"5+true;", "type error: unsupported operand types for +: int and bool"},
		{"5+true; 5;", "type error: unsupported operand types for +: int and bool"},
		{"-true", "type error: expected int or float to follow - operator (got bool)"},
		{"3--", "name error: \"3\" is not defined"},
		{"true+false", "type error: unsupported operand types for +: bool and bool"},
		{"5;true+false;5", "type error: unsupported operand types for +: bool and bool"},
		{"if (10>1) { true+false;}", "type error: unsupported operand types for +: bool and bool"},
		{`if (10 > 1) {
      if (10>1) {
			return true+false;
			}
			return 1;
}`, "type error: unsupported operand types for +: bool and bool"},
		{"foobar", "name error: \"foobar\" is not defined"},
		{`"Hello" - "World"`, "type error: unsupported operand types for -: string and string"},
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
		{"var a=5;a;", 5},
		{"var a=5*5; a;", 25},
		{"var a=5; var b=a; b;", 5},
		{"var a=5; a--; a;", 4},
		{"var a=5; a++; a;", 6},
		{"var a=5; var b=a; var c=a+b+5; c;", 15},
	}
	for _, tt := range tests {
		testDecimalObject(t, testEval(tt.input), tt.expect)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `func(x) { x+2 }`
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
		{"var identity=func(x){x;}; identity(5);", 5},
		{"var identity=func(x){return x;}; identity(5);", 5},
		{"var double=func(x){x*2;}; double(5);", 10},
		{"var add = func(x, y) { x+y;}; add(5,5);", 10},
		{"var add=func(x,y){x+y;}; add(5+5, add(5,5));", 20},
		{"func(x){x;}(5)", 5},
	}
	for _, tt := range tests {
		testDecimalObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
var newAdder = func(x) {
	func(y) { x + y };
};
var addTwo = newAdder(3);
addTwo(2);
`
	testDecimalObject(t, testEval(input), 5)
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
		{`len(1)`, "type error: len() argument is unsupported (int given)"},
		{`len("one", "two")`, "type error: len() takes exactly 1 argument (2 given)"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testDecimalObject(t, evaluated, int64(expected))
		case string:
			if evaluated == object.Nil {
				t.Errorf("Got nil output on input of '%s'\n", tt.input)
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
	result, ok := evaluated.(*object.List)
	if !ok {
		t.Fatalf("object is not Array, got=%T(%v)",
			evaluated, evaluated)
	}
	if len(result.Items) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Items))
	}
	testDecimalObject(t, result.Items[0], 1)
	testDecimalObject(t, result.Items[1], 4)
	testDecimalObject(t, result.Items[2], 6)
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
			"var i =0; [1][i]",
			1,
		},
		{
			"var myArray=[1,2,3];myArray[2];",
			3,
		},
		{
			"var myArray=[1,2,3];myArray[0]+myArray[1]+myArray[2]",
			6,
		},
		{
			"var myArray=[1,2,3];var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1,2,3][3]",
			"index error: array index out of range: 3",
		},
		{
			"[1,2,3][-1]",
			3,
		},
		{
			"[1,2,3][-2]",
			2,
		},
		{
			"[1,2,3][-3]",
			1,
		},
		{
			"[1,2,3][-4]",
			"index error: array index out of range: -4",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDecimalObject(t, evaluated, int64(integer))
		} else {
			err, ok := evaluated.(*object.Error)
			require.True(t, ok)
			require.Equal(t, tt.expected.(string), err.Message)
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
			"\"Steve\"[-1]",
			"e",
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
			testNilObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `var two="two";
	{
		"one":10-9,
		two:1+1,
		"thr" + "ee" : 6/2,
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval did't return Hash. got=%T(%+v)",
			evaluated, evaluated)
	}
	require.Len(t, result.Items, 3)

	one, ok := result.Get("one").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(1), one.Value)

	two, ok := result.Get("two").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(2), two.Value)

	thr, ok := result.Get("three").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(3), thr.Value)
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
			`var key = "foo"; {"foo":5}[key]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDecimalObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestForLoopSimple(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"for x := 0; x < 0; x++ { x }", nil},
		{"for x := 0; x < 1; x++ { x }", int64(0)},
		{"for x := 0; x < 2; x++ { x }", int64(1)},
		{"for x:=0;x<2;x++{x}", int64(1)},
		{"for x := 0; x < 10; x ++ { x }", int64(9)},
		{"for x := 0; x < 10; x += 2 { x }", int64(8)},
		{"for x := 0; x < 10; x += 2 { 999 }", int64(999)},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		value := object.ToGoType(evaluated)
		require.Equal(t, tt.expected, value)
	}
}

func TestForLoopAdditional(t *testing.T) {
	input := `
var sum = 0;
var up = 100;
for x := 1; x < up; x += 1 {
	sum = sum + x;
}
sum
`
	evaluated := testEval(input)
	testDecimalObject(t, evaluated, 4950)
}

func TestForLoopBreak(t *testing.T) {
	input := `
var sum = 0;
for x := 0; x < 100; x++ {
	if x == 2 {
		break
	}
	sum += x
}
sum
`
	evaluated := testEval(input)
	testDecimalObject(t, evaluated, 1)
}

func TestSimpleLoop(t *testing.T) {
	input := `
x := 0
for {
	if x == 100 {
		break
	}
	x += 1
	x
}
`
	evaluated := testEval(input)
	testDecimalObject(t, evaluated, 100)
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
			"int",
		},
		{
			"type( 3.14159 );",
			"float",
		},
		{
			"type( [1,2,3] );",
			"list",
		},
		{
			"type( { \"name\":\"monkey\" } );",
			"map",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, str)
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	input := "for i := 0; i < 999999999999; i++ { i }"
	program, err := parser.Parse(input)
	require.Nil(t, err)

	s := scope.New(scope.Opts{})
	e := &Evaluator{}
	evaluated := e.Evaluate(ctx, program, s)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Errorf("no error object returned. got=%T(%+v)",
			evaluated, evaluated)
	}
	if !strings.Contains(errObj.Message, "deadline") {
		t.Errorf("got error, but wasn't timeout: %s", errObj.Message)
	}
}

func TestSet(t *testing.T) {
	e := &Evaluator{}
	input := `{1, 2, 3}`
	ctx := context.Background()

	program, err := parser.Parse(input)
	require.Nil(t, err)

	s := scope.New(scope.Opts{})
	evaluated := e.Evaluate(ctx, program, s)

	set, ok := evaluated.(*object.Set)
	require.True(t, ok)
	require.Len(t, set.Items, 3)

	hk1 := (&object.Int{Value: 1}).HashKey()
	hk2 := (&object.Int{Value: 2}).HashKey()
	hk3 := (&object.Int{Value: 3}).HashKey()

	require.Equal(t, int64(1), set.Items[hk1].(*object.Int).Value)
	require.Equal(t, int64(2), set.Items[hk2].(*object.Int).Value)
	require.Equal(t, int64(3), set.Items[hk3].(*object.Int).Value)
}

func TestIndexErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1,2,3][99]", "index error: array index out of range: 99"},
		{`{"foo":1}["bar"]`, "key error: \"bar\""},
		{`"foo"[4]`, "index error: string index out of range: 4"},
	}
	for _, tt := range tests {
		resultErr, ok := testEval(tt.input).(*object.Error)
		require.True(t, ok)
		require.Equal(t, tt.expected, resultErr.Message)
	}
}

func TestStringInterpolation(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		input    string
		expected string
	}{
		{"`{10+3}`", "{10+3}"},
		{"'{10+3}'", "13"},
		{`"{10+3}"`, "{10+3}"},
		{`'hey, {}{strings.to_upper(name) + \"!\"}'`, "hey, JOE!"},
		{`'length: {len(name)}'`, "length: 3"},
		{`'{{1,2,  3}} is a set'`, "{1,2,  3} is a set"},
		{`'{"hey"}'`, "hey"},
		{`'a\'b'`, `a'b`},
		{`'a\'b'`, `a'b`},
		{`"a\'b"`, `a'b`},
		{`'a"b'`, `a"b`},
		{`'a\"b'`, `a"b`},
	}
	for _, tt := range tests {
		program, err := parser.Parse(tt.input)
		require.Nil(t, err)
		s := scope.New(scope.Opts{})
		mod, err := modStrings.Module(s)
		require.Nil(t, err)
		s.Declare("name", &object.String{Value: "Joe"}, true)
		s.Declare("strings", mod, true)
		e := &Evaluator{}
		obj := e.Evaluate(ctx, program, s)
		str, ok := obj.(*object.String)
		require.True(t, ok)
		require.Equal(t, tt.expected, str.Value)
	}
}

func TestMethodCallOnError(t *testing.T) {
	ctx := context.Background()
	input := `flub(1).whatever()`
	program, err := parser.Parse(input)
	require.Nil(t, err)
	s := scope.New(scope.Opts{})
	e := &Evaluator{}
	obj := e.Evaluate(ctx, program, s)
	fmt.Println(obj)
	errObj, ok := obj.(*object.Error)
	require.True(t, ok)
	fmt.Println(errObj)
	require.Equal(t, "name error: \"flub\" is not defined", errObj.Message)
}
