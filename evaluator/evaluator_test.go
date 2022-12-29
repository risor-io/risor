package evaluator

import (
	"context"
	"errors"
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
	program, err := parser.Parse(input)
	if err != nil {
		panic(err)
	}
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
	if result.Value() != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value(), expected)
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
	if math.Abs(result.Value()-expected) > 0.00001 {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value(), expected)
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
	if result.Value() != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value(), expected)
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
	if result.Value() != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value(), expected)
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
		{"!!!true", false},
		{"!!!false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
	result := testEval("!5")
	resultErr, ok := result.(*object.Error)
	require.True(t, ok)
	require.Equal(t, "type error: expected boolean to follow ! operator (got int)",
		resultErr.Message().Value())
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
		require.Equal(t, tt.expected, testEval(tt.input).Interface())
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
		evaluated, ok := testEval(tt.input).Interface().(error)
		require.True(t, ok)
		require.Equal(t, errors.New(tt.expectedMessage), evaluated)
	}
}

func TestVarStatements(t *testing.T) {
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
	params := fn.Parameters()
	if len(params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters())
	}
	if params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", params[0])
	}
	expectedBody := `(x + 2)`
	if fn.Body().String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body())
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
	if str.Value() != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value())
	}
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, int64(0)},
		{`len("four")`, int64(4)},
		{`len("狐犬")`, int64(2)},
		{`len("hello world")`, int64(11)},
		{`len(1)`, errors.New("type error: len() argument is unsupported (int given)")},
		{`len("one", "two")`, errors.New("type error: len() takes exactly 1 argument (2 given)")},
		{`len({1,2})`, int64(2)},
		{`len({one: 1})`, int64(1)},
		{`float("1.3")`, float64(1.3)},
		{`float(-11)`, float64(-11.0)},
		{`float("oops")`, errors.New("value error: invalid literal for float(): \"oops\"")},
		{`int(3)`, int64(3)},
		{`int("oops")`, errors.New("value error: invalid literal for int(): \"oops\"")},
		{`int(float(2.0))`, int64(2)},
		{`call(keys, {"one": 1, "two": 2})`, []any{"one", "two"}},
		{`reversed(sorted([3, 99, 1, 2]))`, []any{int64(99), int64(3), int64(2), int64(1)}},
		{`bool()`, false},
		{`bool([])`, false},
		{`bool([1])`, true},
		{`all(list())`, true},
		{`all([1, 2, 3])`, true},
		{`all([0, 2, 3])`, false},
		{`all(set())`, true},
		{`all({1})`, true},
		{`all({1,0})`, false},
		{`any([1, 2, 3])`, true},
		{`any([0, 2, 3])`, true},
		{`any([0])`, false},
		{`any(set())`, false},
		{`any({1})`, true},
		{`any({1,0})`, true},
		{`any({false,0})`, false},
		{`assert(false)`, errors.New("assertion failed")},
		{`assert(false, "sadface")`, errors.New("sadface")},
		{`assert(true)`, nil},
		{`type(ok("yay"))`, "result"},
		{`type(err("sadface"))`, "result"},
		{`string(3.3)`, "3.3"},
		{`sprintf("a%sc", "b")`, "abc"},
		{`sprintf("%02d%t %s!?", 3, false, 'what')`, "03false what!?"},
		{`sprintf("m")`, "m"},
		{`sprintf()`, errors.New("type error: sprintf() takes 1 or more arguments (0 given)")},
		{`m := {"one": 1, two: 2}; delete(m, "two")`, map[string]any{"one": int64(1)}},
		{`type(getattr([1], "append"))`, "builtin"},
		{`func x(){}; type(x)`, "function"},
		{`set([1,2]) | len`, int64(2)},
		{`set("abc") | len`, int64(3)},
		{`set({one:1}) | list`, []any{"one"}},
		{`set(1)`, errors.New("type error: set() argument is unsupported (int given)")},
		{`ord(chr(12345))`, int64(12345)},
		{`keys([9,8,7])`, []any{int64(0), int64(1), int64(2)}},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, testEval(tt.input).Interface())
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`[1, 2*2, 3+3]`, "[1, 4, 6]"},
		{`len([])`, "0"},
		{`len([1, 2])`, "2"},
		{`[1].append(2).append(3)`, "[1, 2, 3]"},
		{`[1].append(2).append(3)[1]`, "2"},
		{`l := [1,2]; append := l.append; call(append, 3)`, "[1, 2, 3]"},
		{`l := [1,2]; l`, "[1, 2]"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, testEval(tt.input).Inspect())
	}
}

func TestListIndex(t *testing.T) {
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
			"var myList=[1,2,3];myList[2];",
			3,
		},
		{
			"var myList=[1,2,3];myList[0]+myList[1]+myList[2]",
			6,
		},
		{
			"var myList=[1,2,3];var i = myList[0]; myList[i]",
			2,
		},
		{
			"[1,2,3][3]",
			"index error: index out of range: 3",
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
			"index error: index out of range: -4",
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
			require.Equal(t, tt.expected.(string), err.Message().Value())
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

func TestMapLiterals(t *testing.T) {
	input := `var two="ignored";
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
	require.Len(t, result.Value(), 3)

	one, ok := result.Get("one").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(1), one.Value())

	two, ok := result.Get("two").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(2), two.Value())

	thr, ok := result.Get("three").(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(3), thr.Value())
}

func TestMapIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo":5}["foo"]`,
			5,
		},
		{
			`var key = "foo"; {"foo":10}[key]`,
			10,
		},
		{
			`{foo:42}["foo"]`,
			42,
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

func TestForLoop(t *testing.T) {
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
		value := testEval(tt.input).Interface()
		require.Equal(t, tt.expected, value)
	}
}

func TestForLoopScope(t *testing.T) {
	input := `
sum := 0.0
for x := 0; x < 10; x++ {
	myvar := x
	sum += myvar
}
myvar
`
	evaluated, ok := testEval(input).(*object.Error)
	require.True(t, ok)
	require.Equal(t, `name error: "myvar" is not defined`, evaluated.Message().Value())
	// require.True(t, false)
}

func TestForLoopVariant(t *testing.T) {
	input := `
sum := 0
for sum < 10 {
	sum += 1
}
sum
`
	evaluated := testEval(input)
	require.Equal(t, int64(10), evaluated.Interface())
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
	e := New(Opts{})
	evaluated := e.Evaluate(ctx, program, s)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Errorf("no error object returned. got=%T(%+v)",
			evaluated, evaluated)
	}
	msg := errObj.Message().Value()
	if !strings.Contains(msg, "deadline") {
		t.Errorf("got error, but wasn't timeout: %s", msg)
	}
}

func TestSet(t *testing.T) {
	e := New(Opts{})
	input := `{1, 2, 3}`
	ctx := context.Background()

	program, err := parser.Parse(input)
	require.Nil(t, err)

	s := scope.New(scope.Opts{})
	evaluated := e.Evaluate(ctx, program, s)

	set, ok := evaluated.(*object.Set)
	require.True(t, ok)
	require.Len(t, set.Value(), 3)

	hk1 := (object.NewInt(1)).HashKey()
	hk2 := (object.NewInt(2)).HashKey()
	hk3 := (object.NewInt(3)).HashKey()

	items := set.Value()
	require.Equal(t, int64(1), items[hk1].(*object.Int).Value())
	require.Equal(t, int64(2), items[hk2].(*object.Int).Value())
	require.Equal(t, int64(3), items[hk3].(*object.Int).Value())
}

func TestIndexErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1,2,3][99]", "index error: index out of range: 99"},
		{`{"foo":1}["bar"]`, "key error: \"bar\""},
		{`"foo"[4]`, "index error: index out of range: 4"},
	}
	for _, tt := range tests {
		resultErr, ok := testEval(tt.input).(*object.Error)
		require.True(t, ok)
		require.Equal(t, tt.expected, resultErr.Message().Value())
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
		s.Declare("name", object.NewString("Joe"), true)
		s.Declare("strings", mod, true)
		e := New(Opts{})
		obj := e.Evaluate(ctx, program, s)
		str, ok := obj.(*object.String)
		require.True(t, ok)
		require.Equal(t, tt.expected, str.Value())
	}
}

func TestMethodCallOnError(t *testing.T) {
	ctx := context.Background()
	input := `flub(1).whatever()`
	program, err := parser.Parse(input)
	require.Nil(t, err)
	obj := New(Opts{}).Evaluate(ctx, program, scope.New(scope.Opts{}))
	errObj, ok := obj.(*object.Error)
	require.True(t, ok)
	require.Equal(t, "name error: \"flub\" is not defined", errObj.Message().Value())
}

func TestPipeExpression(t *testing.T) {
	ctx := context.Background()
	input := `
	func pad(s, chr="_") { chr + s + chr }
	["a","b","c"] | len | string | pad("#")`
	program, err := parser.Parse(input)
	require.Nil(t, err)
	obj := New(Opts{}).Evaluate(ctx, program, scope.New(scope.Opts{}))
	str, ok := obj.(*object.String)
	require.True(t, ok, obj)
	require.Equal(t, "#3#", str.Value())
}

func TestTry(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		input    string
		expected string
	}{
		{`try(1)`, "1"},
		{`try(1 + 2)`, "3"},
		{`try("ok")`, `"ok"`},
		{`try("ok", "fallback")`, `"ok"`},
		{`try(ok("yay"), "fallback")`, `"yay"`},
		{`try(err("kaboom"), "fallback")`, `"fallback"`},
		{`try(error("kaboom"), "fallback")`, `"fallback"`},
		{`try(error("kaboom"), error("ouch"))`, `error("ouch")`},
		{`try(error("kaboom"), func(msg) { msg })`, `"kaboom"`},
		{`try(error("kaboom"), func(msg) { 42 })`, `42`},
	}
	for _, tt := range tests {
		program, err := parser.Parse(tt.input)
		require.Nil(t, err)
		obj := New(Opts{}).Evaluate(ctx, program, scope.New(scope.Opts{}))
		require.Equal(t, tt.expected, obj.Inspect())
	}
}

func TestMisc(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`const x = "ok"`, "ok"},
		{`const x = whut`, errors.New("name error: \"whut\" is not defined")},
		{`const x = 1; const x = 2`, errors.New("assignment error: \"x\" is already set")},
		{`const x = 1; x = 2`, errors.New("assignment error: \"x\" is read-only")},
		{`var x = 1; var x = 2`, errors.New("assignment error: \"x\" is already set")},
		{`var x = huh`, errors.New("name error: \"huh\" is not defined")},
		{`var x = 1; x = huh`, errors.New("name error: \"huh\" is not defined")},
		{`var x = [0]; x[0] = 42; x`, []any{int64(42)}},
		{`var x = 1; x += 42`, int64(43)},
		{`var x = 1; x += y`, errors.New("name error: \"y\" is not defined")},
		{`var x = 1; x += error("kaboom")`, errors.New("kaboom")},
		{`var x = 1; x -= 42`, int64(-41)},
		{`var x = 1; x *= 42`, int64(42)},
		{`var x = 42; x /= 2`, int64(21)},
		{`x := {1}; x[1]`, true},
		{`x := 1; x[1]`, errors.New("type error: int object is not scriptable")},
		{`x := 1; x[1]=2`, errors.New("type error: int is not a container")},
		{`x := [9,8,7]; x[1:]`, []any{int64(8), int64(7)}},
		{`x := [9,8,7]; x[:1]`, []any{int64(9)}},
		{`x := [9,8,7]; x[1:1]`, []any{}},
		{`x := [9,8,7]; x[1:3]`, []any{int64(8), int64(7)}},
		{`x := [9,8,7]; x[-2]`, int64(8)},
		{`x := [9,8,7]; x[-2:]`, []any{int64(8), int64(7)}},
		{`x := [9,8,7]; x[-2:-1]`, []any{int64(8)}},
		{`x := [9,8,7]; x[-7:-1]`, errors.New("slice error: start index is out of range")},
		{`x := [9,8,7]; x[1:-7]`, errors.New("slice error: stop index is out of range")},
		{`1 == 1.0`, true},
		{`1.0 == 1`, true},
		{`1 != 1.0`, false},
		{`1.0 != 1`, false},
		{`1 < 2.0`, true},
		{`1 < 0.1`, false},
		{`[1, 1.0, 2] == [1, 1.0, 2]`, true},
		{`[1, 1.0, 2] == [1, 1, 2.0]`, true},
		{`[1, 1.0, 2] == [1, "1.0", 2]`, false},
		{`[1, 1.0, 2] == [1, 1.1, 2]`, false},
		{`error("foo %s", "bar")`, errors.New("foo bar")},
		{`"hi" in [1, 2, 3]`, false},
		{`"hi" in [1, 2, 3, "hi"]`, true},
		{`"hi" in 42`, errors.New("eval error: right hand side of 'in' operator must be a container")},
		{`range [42, 43]`, []map[string]any{
			{"key": int64(0), "value": int64(42)},
			{"key": int64(1), "value": int64(43)},
		}},
		{`range {1,2} | type`, "set_iter"},
		{`range {1,2}.next | string`, `builtin(next)`},
		{`for _, x := range [98,99] { x }`, int64(99)},
		{`for idx, x := range [98,99] { idx }`, int64(1)},
		{`s := {98, 99}; for item := range s { item }`, int64(99)},
		{`x, y, z := [3, 2, 1]; sprintf("%d-%d-%d", x, y, z)`, "3-2-1"},
		{`x, y := [1]`, errors.New("eval error: invalid multi variable assignment (list size: 1; identifiers: 2)")},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, testEval(tt.input).Interface(), tt.input)
	}
}

func FuzzEval(f *testing.F) {
	testcases := []string{
		"1/2+4+=5-[1,2,{}]",
		" ",
		"!12345",
		"var x = [1,2,3];",
		`; const z = {"foo"}`,
		`"foo_" + 1.34 /= 2.0`,
		`{hey: {there: 1}}`,
		`'foo {x[1:3] + 1}'`,
		`x.func(x=1, y=2).bar`,
		`0A=`,
		`return (('hi'+"there"))`,
		`break`,
		`for i := 0; i < 100; i++ { if i == 42 { continue } }`,
		`return`,
		`func x() { return 42 }; x()`,
		`false ? 3.3 : obj.call([1,var,case])`,
		`select { case <-chan: 1; default: 2; continue == (true) || false }`,
		`var x = 99; x /= 99; import "fmt"; fmt/(x)}**%'3'`,
		string("[\x00]"),
		`func(){}()()`,
	}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	f.Fuzz(func(t *testing.T, input string) {
		prog, err := parser.Parse(input) // Confirms no panics
		if err == nil {
			if prog == nil {
				t.Error("nil program")
			} else {
				New(Opts{}).Evaluate(ctx, prog, scope.New(scope.Opts{})) // Confirms no panics
			}
		}
	})
}
