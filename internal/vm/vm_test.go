package vm

import (
	"context"
	"testing"

	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	constants := []object.Object{
		object.NewInt(3),
		object.NewInt(4),
	}
	code := []op.Code{
		op.LoadConst,
		0,
		0,
		op.LoadConst,
		1,
		0,
		op.BinaryOp,
		op.Code(op.Add),
	}
	vm := New(Options{Main: &object.Code{
		Constants:    constants,
		Instructions: code,
		Symbols:      object.NewSymbolTable(),
	}})
	err := vm.Run(context.Background())
	require.Nil(t, err)

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(7), tos)
}

func TestAddCompilationAndExecution(t *testing.T) {
	program, err := parser.Parse(`
	x := 11
	y := 12
	x + y
	`)
	require.Nil(t, err)

	c := compiler.New(compiler.Options{})
	main, err := c.Compile(program)
	require.Nil(t, err)

	consts := main.Constants
	require.Len(t, consts, 2)

	c1, ok := consts[0].(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(11), c1.Value())

	c2, ok := consts[1].(*object.Int)
	require.True(t, ok)
	require.Equal(t, int64(12), c2.Value())

	vm := New(Options{Main: main})
	require.Nil(t, vm.Run(context.Background()))

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(23), tos)
}

func TestConditional(t *testing.T) {
	program, err := parser.Parse(`
	x := 20
	if x > 10 {
		x = 99
	}
	x
	`)
	require.Nil(t, err)

	c := compiler.New(compiler.Options{})
	main, err := c.Compile(program)
	require.Nil(t, err)

	vm := New(Options{Main: main})
	require.Nil(t, vm.Run(context.Background()))

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(99), tos)
}

func TestConditional3(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 5
	y := 10
	if x > 1 {
		y
	} else {
		99
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(10), result)
}

func TestConditional4(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 5
	y := 22
	z := 33
	if x < 1 {
		x = y
	} else {
		x = z
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(33), result)
}

func TestLoop(t *testing.T) {
	result, err := Run(context.Background(), `
	y := 0
	for {
		y = y + 1
		if y > 10 {
			break
		}
	}
	y
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(11), result)
}

func TestForLoop2(t *testing.T) {
	result, err := Run(context.Background(),
		`x := 0; for y := 0; y < 5; y++ { x = y }; x`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestForRange1(t *testing.T) {
	result, err := Run(context.Background(), `
	x := [1, 2.3, "hello", true]
	output := []
	for i := range x {
		1 + 2
		3 + 4
	}
	99
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestForRange2(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 0
	for _, value := range [5,6,7] {
		x = value
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(7), result)
}

func TestForRange3(t *testing.T) {
	result, err := Run(context.Background(), `
	x, y := [0, 0]
	for i, value := range [5, 6, 7] {
		x = i      // should go up to 2
		y = value  // should go up to 7
	}
	[x, y]
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(2),
		object.NewInt(7),
	}), result)
}

func TestForRange4(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 0
	for i := range ["a", "b", "c"] { x = i }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestForRange5(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 0
	for range ["a", "b", "c"] { x++ }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestForRange6(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 0
	r := range { "a", "b" }
	for r { x++ }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestForRange7(t *testing.T) {
	result, err := Run(context.Background(), `
	x := nil
	y := nil
	count := 0
	f := func() { range [ "a", "b", "c" ] }
	for i, value := f() {
		x = i      // should count 0, 1, 2
		y = value  // should go "a", "b", "c"
		count++
	}
	[x, y, count]
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(2),
		object.NewString("c"),
		object.NewInt(3),
	}), result)
}

func TestIterator(t *testing.T) {
	tests := []testCase{
		{`range { 33, 44, 55 }.next().key`, object.NewInt(33)},
		{`range { 33, 44, 55 }.next().value`, object.True},
		{`range [ 33, 44, 55 ].next().key`, object.NewInt(0)},
		{`range [ 33, 44, 55 ].next().value`, object.NewInt(33)},
		{`range "abcd".next().key`, object.NewInt(0)},
		{`range "abcd".next().value`, object.NewString("a")},
		{`range { a: 33, b: 44 }.next().key`, object.NewString("a")},
		{`range { a: 33, b: 44 }.next().value`, object.NewInt(33)},
	}
	runTests(t, tests)
}

func TestStackPopping1(t *testing.T) {
	result, err := Run(context.Background(), `
	x := []
	for i := 0; i < 4; i++ {
		1
		2
		3
		4
		x.append(i)
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(0),
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
	}), result)
}

func TestStackPopping2(t *testing.T) {
	result, err := Run(context.Background(), `
	for i := range [1, 2, 3] {
		1
		2
		3
		4
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestStackBehavior1(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 99
	for i := 0; i < 4; x {
		i++
		1
		2
		3
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestStackBehavior2(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 77
	for i := 0; i < 4; x {
		i++
		1
		2
		3
		4
		if i > 0 {
			break // loop once
		}
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestStackBehavior3(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 77
	if x > 0 {
		99 
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestStackBehavior4(t *testing.T) {
	result, err := Run(context.Background(), `
	x := -1
	if x > 0 {
		99 
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestAssignmentOperators(t *testing.T) {
	result, err := Run(context.Background(), `
	y := 99
	y  = 3
	y += 6
	y /= 9
	y *= 2
	y
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestFunctionCall(t *testing.T) {
	result, err := Run(context.Background(), `
	f := func(x) { 42 + x }
	v := f(1)
	v + 10
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(53), result)
}

func TestSwitch1(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 3
	switch x {
		case 1:
		case 2:
			21
		case 3:
			42
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(42), result)
}

func TestSwitch2(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 1
	switch x {
		case 1:
			99
		case 2:
			42
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestSwitch3(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 3
	switch x {
		case 1:
			99
		case 2:
			42
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestSwitch4(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 3
	switch x { default: 99 }
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestSwitch5(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 3
	switch x { default: 99 case 3: x; x-1 }
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestStr(t *testing.T) {
	result, err := Run(context.Background(), `
	s := "hello"
	s
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("hello"), result)
}

func TestStrLen(t *testing.T) {
	result, err := Run(context.Background(), `
	s := "hello"
	len(s)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(5), result)
}

func TestList1(t *testing.T) {
	result, err := Run(context.Background(), `
	l := [1, 2, 3]
	l
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
	}), result)
}

func TestList2(t *testing.T) {
	result, err := Run(context.Background(), `
	plusOne := func(x) { x + 1 }
	[plusOne(0), 4-2, plusOne(2)]
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
	}), result)
}

func TestMap(t *testing.T) {
	result, err := Run(context.Background(), `
	{"a": 1, "b": 4-2}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewMap(map[string]object.Object{
		"a": object.NewInt(1),
		"b": object.NewInt(2),
	}), result)
}

func TestSet(t *testing.T) {
	result, err := Run(context.Background(), `
	{"a", 4-1}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewSet([]object.Object{
		object.NewString("a"),
		object.NewInt(3),
	}), result)
}

func TestNonLocal(t *testing.T) {
	result, err := Run(context.Background(), `
	y := 3
	z := 99
	f := func() { y = 4 }
	f()
	y
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestFrameLocals1(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 1
	f := func(x) {
		x = 99
	}
	f(4)
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(1), result)
}

func TestFrameLocals2(t *testing.T) {
	result, err := Run(context.Background(), `
	x := 1
	f := func(y) {
		x = 99
	}
	f(4)
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestMapKeys(t *testing.T) {
	result, err := Run(context.Background(), `
	m := {"a": 1, "b": 2}
	keys(m)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("a"),
		object.NewString("b"),
	}), result)
}

func TestClosure(t *testing.T) {
	result, err := Run(context.Background(), `
	f := func(x) { func() { x } }
	closure := f(22)
	closure()
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(22), result)
}

func TestClosureIncrementer(t *testing.T) {
	result, err := Run(context.Background(), `
	f := func(x) {
		func() { x++; x }
	}
	incrementer := f(0)
	incrementer() // 1
	incrementer() // 2
	incrementer() // 3
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestRecursiveExample1(t *testing.T) {
	result, err := Run(context.Background(), `
	func twoexp(n) {
		if n == 0 {
			return 1
		} else {
			return 2 * twoexp(n-1)
		}
	}
	twoexp(4)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(16), result)
}

func TestRecursiveExample2(t *testing.T) {
	result, err := Run(context.Background(), `
	func twoexp(n) {
		a := 1
		b := 2
		c := a * b
		if n == 0 {
			return 1
		} else {
			return c * twoexp(n-1)
		}
	}
	twoexp(4)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(16), result)
}

func TestConstant(t *testing.T) {
	_, err := Run(context.Background(), `const x = 1; x = 2`)
	require.NotNil(t, err)
	require.Equal(t, "cannot assign to constant: x", err.Error())
}

func TestConstantFunction(t *testing.T) {
	_, err := Run(context.Background(), `
	func add(x, y) { x + y }
	add = "bloop"
	`)
	require.NotNil(t, err)
	require.Equal(t, "cannot assign to constant: add", err.Error())
}

func TestMultipleCases(t *testing.T) {

	t.Run("Arithmetic", func(t *testing.T) {
		tests := []testCase{
			{`1 + 2`, object.NewInt(3)},
			{`1 + 2 + 3`, object.NewInt(6)},
			{`1 + 2 * 3`, object.NewInt(7)},
			{`(1 + 2) * 3`, object.NewInt(9)},
			{`5 - 3`, object.NewInt(2)},
			{`12 / 4`, object.NewInt(3)},
			{`3 * (4 + 2)`, object.NewInt(18)},
			{`1.5 + 1.5`, object.NewFloat(3.0)},
			{`1.5 + 2`, object.NewFloat(3.5)},
			{`2 + 1.5`, object.NewFloat(3.5)},
		}
		runTests(t, tests)
	})

	t.Run("Control", func(t *testing.T) {
		tests := []testCase{
			{`x := 1; if x > 5 { 99 } else { 100 }`, object.NewInt(100)},
			{`x := 1; if x > 0 { 99 } else { 100 }`, object.NewInt(99)},
			{`x := 1; y := x > 0 ? 77 : 88; y`, object.NewInt(77)},
			{`x := (1 > 2) ? 77 : 88; x`, object.NewInt(88)},
			{`x := (2 > 1) ? 77 : 88; x`, object.NewInt(77)},
			{`x := 1; switch x { case 1: 99; case 2: 100 }`, object.NewInt(99)},
			{`x := 2; switch x { case 1: 99; case 2: 100 }`, object.NewInt(100)},
			{`x := 3; switch x { case 1: 99; default: 42 case 2: 100 }`, object.NewInt(42)},
			{`x := 3; switch x { case 1: 99; case 2: 100 }`, object.Nil},
			{`x := 3; switch x { case 1, 3: 99; case 2: 100 }`, object.NewInt(99)},
			{`x := 3; switch x { case 1: 99; case 2, 4-1: 100 }`, object.NewInt(100)},
			{`x := 3; switch bool(x) { case true: "wow" }`, object.NewString("wow")},
		}
		runTests(t, tests)
	})

	t.Run("Builtins", func(t *testing.T) {
		tests := []testCase{
			{`len("hello")`, object.NewInt(5)},
			{`len([1, 2, 3])`, object.NewInt(3)},
			{`len({"a": 1})`, object.NewInt(1)},
			{`keys({"a": 1})`, object.NewList([]object.Object{object.NewString("a")})},
			{`type(3.14159)`, object.NewString("float")},
			{`type("hi".contains)`, object.NewString("builtin")},
			{`"hi".contains("h")`, object.True},
			{`"hi".contains("x")`, object.False},
			{`sprintf("%d-%d", 1, 2)`, object.NewString("1-2")},
		}
		runTests(t, tests)
	})

	t.Run("Assignment", func(t *testing.T) {
		tests := []testCase{
			{`a, b := [3, 4]; a`, object.NewInt(3)},
			{`a, b := [3, 4]; b`, object.NewInt(4)},
			{`a, b := "ᛛᛥ"; a`, object.NewString("ᛛ")},
			{`a, b := {42, 43}; a in [42, 43]`, object.True},
			{`a, b := {foo: 1, bar: 2}; a in ["foo", "bar"]`, object.True},
		}
		runTests(t, tests)
	})

	t.Run("Functions", func(t *testing.T) {
		closure := `
z := 10
y := func(x, inc=100) { x + z + inc }
y(3)
`
		tests := []testCase{
			{`func add(x, y) { x + y }; add(3, 4)`, object.NewInt(7)},
			{`func add(x, y) { x + y }; add(3, 4) + 5`, object.NewInt(12)},
			{`func inc(x, amount=1) { x + amount }; inc(3)`, object.NewInt(4)},
			{`func factorial(n) { if (n == 1) { return 1 } else { return n * factorial(n - 1) } }; factorial(5)`, object.NewInt(120)},
			{closure, object.NewInt(113)},
		}
		runTests(t, tests)
	})

	t.Run("DataStructures", func(t *testing.T) {
		tests := []testCase{
			{`true`, object.True},
			{`[1,2,3][2]`, object.NewInt(3)},
			{`"hello"[1]`, object.NewString("e")},
			{`{"x": 10, "y": 20}["x"]`, object.NewInt(10)},
			{`3 in [1, 2, 3]`, object.True},
			{`4 in [1, 2, 3]`, object.False},
			{`{"foo": "bar"}["foo"]`, object.NewString("bar")},
			{`{foo: "bar"}["foo"]`, object.NewString("bar")},
			{`[1, 2, 3, 4, 5].filter(func(x) { x > 3 })`, object.NewList(
				[]object.Object{object.NewInt(4), object.NewInt(5)})},
			{`range [1]`, object.NewListIter(object.NewList([]object.Object{object.NewInt(1)}))},
		}
		runTests(t, tests)
	})

	t.Run("ComparisonAndBoolean", func(t *testing.T) {
		tests := []testCase{
			{`3 < 5`, object.True},
			{`10 > 5`, object.True},
			{`3 <= 5`, object.True},
			{`10 == 5`, object.False},
			{`10 != 5`, object.True},
			{`!true`, object.False},
			{`!false`, object.True},
			{`!!true`, object.True},
			{`!!false`, object.False},
			{`!0`, object.True},
			{`!5`, object.False},
			{`![]`, object.True},
			{`![1]`, object.False},
			{`!{}`, object.True},
			{`!{1}`, object.False},
			{`!""`, object.True},
			{`!"a"`, object.False},
		}
		runTests(t, tests)
	})

	t.Run("Strings", func(t *testing.T) {
		tests := []testCase{
			{`"hello" + " " + "world"`, object.NewString("hello world")},
			{`"hello".contains("e")`, object.True},
			{`"hello".contains("x")`, object.False},
			{`"hello".contains("ello")`, object.True},
			{`"hello".contains("ellx")`, object.False},
			{`"hello".contains("")`, object.True},
			{`"hello"[0]`, object.NewString("h")},
			{`"hello"[1]`, object.NewString("e")},
			{`"hello"[-1]`, object.NewString("o")},
			{`"hello"[-2]`, object.NewString("l")},
		}
		runTests(t, tests)
	})

	t.Run("Pipes", func(t *testing.T) {
		tests := []testCase{
			{`"hello" | strings.to_upper`, object.NewString("HELLO")},
			{`"hello" | len`, object.NewInt(5)},
			{`func() { "hello" }() | len`, object.NewInt(5)},
			{`["a", "b"] | strings.join(",") | strings.to_upper`, object.NewString("A,B")},
		}
		runTests(t, tests)
	})

	t.Run("Imports", func(t *testing.T) {
		tests := []testCase{
			// {`import strings; strings`, object.NewModule("strings", nil)},
		}
		runTests(t, tests)
	})
}

type testCase struct {
	input    string
	expected object.Object
}

func runTests(t *testing.T, tests []testCase) {
	t.Helper()
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Helper()
			result, err := Run(ctx, tt.input)
			require.Nil(t, err)
			require.NotNil(t, result)
			require.Equal(t, tt.expected, result)
		})
	}
}
