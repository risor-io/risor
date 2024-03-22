package vm

import (
	"context"
	"testing"
	"time"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/stretchr/testify/require"
)

func TestAddCompilationAndExecution(t *testing.T) {
	program, err := parser.Parse(context.Background(), `
	x := 11
	y := 12
	x + y
	`)
	require.Nil(t, err)

	c, err := compiler.New()
	require.Nil(t, err)

	main, err := c.Compile(program)
	require.Nil(t, err)

	constsCount := main.ConstantsCount()
	require.Equal(t, 2, constsCount)

	c1, ok := main.Constant(0).(int64)
	require.True(t, ok)
	require.Equal(t, int64(11), c1)

	c2, ok := main.Constant(1).(int64)
	require.True(t, ok)
	require.Equal(t, int64(12), c2)

	vm := New(main)
	require.Nil(t, vm.Run(context.Background()))

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(23), tos)
}

func TestConditional(t *testing.T) {
	program, err := parser.Parse(context.Background(), `
	x := 20
	if x > 10 {
		x = 99
	}
	x
	`)
	require.Nil(t, err)

	main, err := compiler.Compile(program)
	require.Nil(t, err)

	vm := New(main)
	require.Nil(t, vm.Run(context.Background()))

	tos, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(99), tos)
}

func TestConditional3(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(),
		`x := 0; for y := 0; y < 5; y++ { x = y }; x`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestForLoop3(t *testing.T) {
	result, err := run(context.Background(), `x := 0; for x < 10 { x++ }; x`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(10), result)
}

func TestForRange1(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	x := 0
	for i := range ["a", "b", "c"] { x = i }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestForRange5(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for range ["a", "b", "c"] { x++ }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestForRange6(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	r := range { "a", "b" }
	for r { x++ }
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestForRange7(t *testing.T) {
	result, err := run(context.Background(), `
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
		{`(range { 33, 44, 55 }).next()`, object.NewInt(33)},
		{`i := range { 33, 44, 55 }; i.next(); i.entry().value`, object.True},
		{`i := range { 33, 44, 55 }; i.next(); i.entry().key`, object.NewInt(33)},
		{`(range [ 33, 44, 55 ]).next()`, object.NewInt(33)},
		{`i := range "abcd"; i.next(); i.entry().key`, object.NewInt(0)},
		{`i := range "abcd"; i.next(); i.entry().value`, object.NewString("a")},
		{`(range { a: 33, b: 44 }).next()`, object.NewString("a")},
		{`i := range { a: 33, b: 44 }; i.next(); i.entry().key`, object.NewString("a")},
		{`i := range { a: 33, b: 44 }; i.next(); i.entry().value`, object.NewInt(33)},
	}
	runTests(t, tests)
}

func TestIndexing(t *testing.T) {
	tests := []testCase{
		{`x := [1, 2]; x[0] = 9; x[0]`, object.NewInt(9)},
		{`x := [1, 2]; x[-1] = 9; x[1]`, object.NewInt(9)},
		{`x := {a: 1}; x["a"] = 9; x["a"]`, object.NewInt(9)},
		{`x := {a: 1}; x["b"] = 9; x["b"]`, object.NewInt(9)},
		{`x := { 1 }; x[1]`, object.True},
		{`x := { 1 }; x[2]`, object.False},
	}
	runTests(t, tests)
}

func TestStackPopping1(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	x := 77
	if x > 0 {
		99 
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestStackBehavior4(t *testing.T) {
	result, err := run(context.Background(), `
	x := -1
	if x > 0 {
		99 
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestAssignmentOperators(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	f := func(x) { 42 + x }
	v := f(1)
	v + 10
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(53), result)
}

func TestSwitch1(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	x := 3
	switch x { default: 99 }
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestSwitch5(t *testing.T) {
	result, err := run(context.Background(), `
	x := 3
	switch x { default: 99 case 3: x; x-1 }
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestStr(t *testing.T) {
	result, err := run(context.Background(), `
	s := "hello"
	s
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("hello"), result)
}

func TestStrLen(t *testing.T) {
	result, err := run(context.Background(), `
	s := "hello"
	len(s)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(5), result)
}

func TestList1(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	{"a": 1, "b": 4-2}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewMap(map[string]object.Object{
		"a": object.NewInt(1),
		"b": object.NewInt(2),
	}), result)
}

func TestSet(t *testing.T) {
	result, err := run(context.Background(), `
	{"a", 4-1}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewSet([]object.Object{
		object.NewString("a"),
		object.NewInt(3),
	}), result)
}

func TestNonLocal(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	x := 1
	f := func(x) { x = 99 }
	f(4)
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(1), result)
}

func TestFrameLocals2(t *testing.T) {
	result, err := run(context.Background(), `
	x := 1
	f := func(y) { x = 99 }
	f(4)
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(99), result)
}

func TestMapKeys(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
	f := func(x) { func() { x } }
	closure := f(22)
	closure()
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(22), result)
}

func TestClosureIncrementer(t *testing.T) {
	result, err := run(context.Background(), `
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

func TestClosureOverLocal(t *testing.T) {
	result, err := run(context.Background(), `
	var testValue = 100
	func getint() {
		var foo = testValue + 1
		func inner() {
			foo
		}
		return inner
	}
	getint()()
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(101), result)
}

func TestClosureManyVariables(t *testing.T) {
	result, err := run(context.Background(), `
	func foo(a, b, c) {
		return func(d) {
			return [a, b, c, d]
		}
	}
	foo("hello", "world", "risor")("go")
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewStringList([]string{"hello", "world", "risor", "go"}), result)
}

func TestRecursiveExample1(t *testing.T) {
	result, err := run(context.Background(), `
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
	result, err := run(context.Background(), `
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
	_, err := run(context.Background(), `const x = 1; x = 2`)
	require.NotNil(t, err)
	require.Equal(t, "compile error: cannot assign to constant \"x\"", err.Error())
}

func TestConstantFunction(t *testing.T) {
	_, err := run(context.Background(), `
	func add(x, y) { x + y }
	add = "bloop"
	`)
	require.NotNil(t, err)
	require.Equal(t, "compile error: cannot assign to constant \"add\"", err.Error())
}

func TestStatementsNilValue(t *testing.T) {
	// The result value of a statement is always nil
	tests := []testCase{
		{`x := 0`, object.Nil},
		{`x := 0; x++`, object.Nil},
		{`x := 0; x--`, object.Nil},
		{`x := 0; x += 1`, object.Nil},
		{`x := 0; x -= 1`, object.Nil},
		{`const x = 0`, object.Nil},
		{`var x = 0`, object.Nil},
		{`x, y := [0, 0]`, object.Nil},
		{`x := [1]; x[0] = 2`, object.Nil},
		{`for i := 0; i < 10; i++ { 42 }`, object.Nil},
		{`x := 0; for i := 0; i < 10; i++ { x = i }`, object.Nil},
	}
	runTests(t, tests)
}

func TestArithmetic(t *testing.T) {
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
		{`2 ** 3`, object.NewInt(8)},
		{`2.0 ** 3.0`, object.NewFloat(8.0)},
		{`1 % 3`, object.NewInt(1)},
		{`3 % 3`, object.NewInt(0)},
		{`11 % 3`, object.NewInt(2)},
		{`-11`, object.NewInt(-11)},
		{`x := -11; -x`, object.NewInt(11)},
		{`-1.5`, object.NewFloat(-1.5)},
	}
	runTests(t, tests)
}

func TestNumericComparisons(t *testing.T) {
	tests := []testCase{
		// Integers
		{`3 < 5`, object.True},
		{`3 <= 5`, object.True},
		{`3 > 5`, object.False},
		{`3 >= 5`, object.False},
		{`3 == 5`, object.False},
		{`3 != 5`, object.True},
		{`2 < 2`, object.False},
		{`2 <= 2`, object.True},
		{`2 > 2`, object.False},
		{`2 >= 2`, object.True},
		{`2 == 2`, object.True},
		{`2 != 2`, object.False},
		// Mixed integers and floats
		{`3.0 < 5`, object.True},
		{`3.0 <= 5`, object.True},
		{`3.0 > 5`, object.False},
		{`3.0 >= 5`, object.False},
		{`3.0 == 5`, object.False},
		{`3.0 != 5`, object.True},
		{`2.0 < 2`, object.False},
		{`2.0 <= 2`, object.True},
		{`2.0 > 2`, object.False},
		{`2.0 >= 2`, object.True},
		{`2.0 == 2`, object.True},
		{`2.0 != 2`, object.False},
		// Floats
		{`3.0 < 5.0`, object.True},
		{`3.0 <= 5.0`, object.True},
		{`3.0 > 5.0`, object.False},
		{`3.0 >= 5.0`, object.False},
		{`3.0 == 5.0`, object.False},
		{`3.0 != 5.0`, object.True},
		{`2.0 < 2.0`, object.False},
		{`2.0 <= 2.0`, object.True},
		{`2.0 > 2.0`, object.False},
		{`2.0 >= 2.0`, object.True},
		{`2.0 == 2.0`, object.True},
		{`2.0 != 2.0`, object.False},
	}
	runTests(t, tests)
}

func TestBooleans(t *testing.T) {
	tests := []testCase{
		{`true`, object.True},
		{`false`, object.False},
		{`!true`, object.False},
		{`!false`, object.True},
		{`!!true`, object.True},
		{`!!false`, object.False},
		{`false == false`, object.True},
		{`false == true`, object.False},
		{`false != false`, object.False},
		{`false != true`, object.True},
		{`true == true`, object.True},
		{`true == false`, object.False},
		{`true != true`, object.False},
		{`true != false`, object.True},
		{`type(true)`, object.NewString("bool")},
		{`type(false)`, object.NewString("bool")},
	}
	runTests(t, tests)
}

func TestTruthiness(t *testing.T) {
	tests := []testCase{
		{`!0`, object.True},
		{`!5`, object.False},
		{`![]`, object.True},
		{`![1]`, object.False},
		{`!{}`, object.True},
		{`!{1}`, object.False},
		{`!""`, object.True},
		{`!"a"`, object.False},
		{`bool(0)`, object.False},
		{`bool(5)`, object.True},
		{`bool([])`, object.False},
		{`bool([1])`, object.True},
		{`bool({})`, object.False},
		{`bool({1})`, object.True},
		{`bool({foo: 1})`, object.True},
	}
	runTests(t, tests)
}

func TestControlFlow(t *testing.T) {
	tests := []testCase{
		{`if false { 3 }`, object.Nil},
		{`if true { 3 }`, object.NewInt(3)},
		{`if false { 3 } else { 4 }`, object.NewInt(4)},
		{`if true { 3 } else { 4 }`, object.NewInt(3)},
		{`if false { 3 } else if false { 4 } else { 5 }`, object.NewInt(5)},
		{`if true { 3 } else if false { 4 } else { 5 }`, object.NewInt(3)},
		{`if false { 3 } else if true { 4 } else { 5 }`, object.NewInt(4)},
		{`x := 1; if x > 5 { 99 } else { 100 }`, object.NewInt(100)},
		{`x := 1; if x > 0 { 99 } else { 100 }`, object.NewInt(99)},
		{`x := 1; y := x > 0 ? 77 : 88; y`, object.NewInt(77)},
		{`x := (1 > 2) ? 77 : 88; x`, object.NewInt(88)},
		{`x := (2 > 1) ? 77 : 88; x`, object.NewInt(77)},
		{`x := 1; switch x { case 1: 99; case 2: 100 }`, object.NewInt(99)},
		{`switch 2 { case 1: 99; case 2: 100 }`, object.NewInt(100)},
		{`switch 3 { case 1: 99; default: 42 case 2: 100 }`, object.NewInt(42)},
		{`switch 3 { case 1: 99; case 2: 100 }`, object.Nil},
		{`switch 3 { case 1, 3: 99; case 2: 100 }`, object.NewInt(99)},
		{`switch 3 { case 1: 99; case 2, 4-1: 100 }`, object.NewInt(100)},
		{`x := 3; switch bool(x) { case true: "wow" }`, object.NewString("wow")},
		{`x := 0; switch bool(x) { case true: "wow" }`, object.Nil},
	}
	runTests(t, tests)
}

func TestLength(t *testing.T) {
	tests := []testCase{
		{`len("")`, object.NewInt(0)},
		{`len([])`, object.NewInt(0)},
		{`len({})`, object.NewInt(0)},
		{`len("hello")`, object.NewInt(5)},
		{`len([1, 2, 3])`, object.NewInt(3)},
		{`len({"abc": 1})`, object.NewInt(1)},
		{`len({"abc"})`, object.NewInt(1)},
		{`len("ᛛᛥ")`, object.NewInt(2)},
		{`len(string(byte_slice([0, 1, 2])))`, object.NewInt(3)},
	}
	runTests(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []testCase{
		{`len("hello")`, object.NewInt(5)},
		{`keys({"a": 1})`, object.NewList([]object.Object{
			object.NewString("a"),
		})},
		{`byte(9)`, object.NewByte(9)},
		{`byte_slice([9])`, object.NewByteSlice([]byte{9})},
		{`float_slice([9])`, object.NewFloatSlice([]float64{9})},
		{`type(3.14159)`, object.NewString("float")},
		{`type("hi".contains)`, object.NewString("builtin")},
		{`sprintf("%d-%d", 1, 2)`, object.NewString("1-2")},
		{`int("99")`, object.NewInt(99)},
		{`float("2.5")`, object.NewFloat(2.5)},
		{`string(99)`, object.NewString("99")},
		{`string(2.5)`, object.NewString("2.5")},
		{`ord("a")`, object.NewInt(97)},
		{`chr(97)`, object.NewString("a")},
		{`encode("hi", "hex")`, object.NewString("6869")},
		{`encode("hi", "base64")`, object.NewString("aGk=")},
		{`iter("abc").next()`, object.NewString("a")},
		{`i := iter("abc"); i.next(); i.entry().key`, object.NewInt(0)},
		{`i := iter("abc"); i.next(); i.entry().value`, object.NewString("a")},
		{`reversed("abc")`, object.NewString("cba")},
		{`reversed([1, 2, 3])`, object.NewList([]object.Object{
			object.NewInt(3),
			object.NewInt(2),
			object.NewInt(1),
		})},
		{`sorted([3, -2, 2])`, object.NewList([]object.Object{
			object.NewInt(-2),
			object.NewInt(2),
			object.NewInt(3),
		})},
		{`any([])`, object.False},
		{`any([0, false, {}])`, object.False},
		{`any([0, false, {foo: 42}])`, object.True},
		{`all([])`, object.True},
		{`all([1, false, {foo: 42}])`, object.False},
		{`all([1, true, {foo: 42}])`, object.True},
	}
	runTests(t, tests)
}

func TestTry(t *testing.T) {
	tests := []testCase{
		{`try(1)`, object.NewInt(1)},
		{`try(1, 2)`, object.NewInt(1)},
		{`try(func() { error("oops") }, "nope")`, object.NewString("nope")},
		{`try(func() { error("oops") }, func() { error("oops") })`, object.Nil},
		{`try(func() { error("oops") }, func() { error("oops") }, 1)`, object.NewInt(1)},
		{`x := 0; y := 0; z := try(func() {
			x = 11
			error("oops")
			x = 12
		  }, func() {
			y = 21
			error("oops")
			y = 22
		  }, 33); [x, y, z]`, object.NewList([]object.Object{
			object.NewInt(11),
			object.NewInt(21),
			object.NewInt(33),
		})},
	}
	runTests(t, tests)
}

func TestMultiVarAssignment(t *testing.T) {
	tests := []testCase{
		{`a, b := [3, 4]; a`, object.NewInt(3)},
		{`a, b := [3, 4]; b`, object.NewInt(4)},
		{`a, b, c := [3, 4, 5]; a`, object.NewInt(3)},
		{`a, b, c := [3, 4, 5]; b`, object.NewInt(4)},
		{`a, b, c := [3, 4, 5]; c`, object.NewInt(5)},
		{`a, b := "ᛛᛥ"; a`, object.NewString("ᛛ")},
		{`a, b := "ᛛᛥ"; b`, object.NewString("ᛥ")},
		{`a, b := {42, 43}; a`, object.NewInt(42)},
		{`a, b := {42, 43}; b`, object.NewInt(43)},
		{`a, b := {foo: 1, bar: 2}; a`, object.NewString("bar")},
		{`a, b := {foo: 1, bar: 2}; b`, object.NewString("foo")},
	}
	runTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []testCase{
		{`func add(x, y) { x + y }; add(3, 4)`, object.NewInt(7)},
		{`func add(x, y) { x + y }; add(3, 4) + 5`, object.NewInt(12)},
		{`func inc(x, amount=1) { x + amount }; inc(3)`, object.NewInt(4)},
		{`func factorial(n) { if (n == 1) { return 1 } else { return n * factorial(n - 1) } }; factorial(5)`, object.NewInt(120)},
		{`z := 10; y := func(x, inc=100) { x + z + inc }; y(3)`, object.NewInt(113)},
		{`func(x="a", y="b") { x + y }()`, object.NewString("ab")},
		{`func(x="a", y="b") { x + y + "c" }()`, object.NewString("abc")},
		{`func(x="a", y="b") { x + y + "c" }("W")`, object.NewString("Wbc")},
		{`func(x="a", y="b") { x + y + "c" }("W", "X")`, object.NewString("WXc")},
		{`func(x="a", y="b") { return "X"; x + y + "c" }()`, object.NewString("X")},
		{`x := 1; func() { y := 10; x + y }()`, object.NewInt(11)},
		{`x := 1; func() { func() { y := 10; x + y } }()()`, object.NewInt(11)},
	}
	runTests(t, tests)
}

func TestContainers(t *testing.T) {
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
}

func TestStrings(t *testing.T) {
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
		{`a := 1; b := "ok"; '{a + 1}-{b | strings.to_upper}'`, object.NewString("2-OK")},
		{`func(a, b) { return 'A: {a} B: {b}' }("hi", "bye")`, object.NewString("A: hi B: bye")},
	}
	runTests(t, tests)
}

func TestPipes(t *testing.T) {
	tests := []testCase{
		{`"hello" | strings.to_upper`, object.NewString("HELLO")},
		{`"hello" | len`, object.NewInt(5)},
		{`func() { "hello" }() | len`, object.NewInt(5)},
		{`["a", "b"] | strings.join(",") | strings.to_upper`, object.NewString("A,B")},
		{`func() { "a" } | call`, object.NewString("a")},
		{`"abc" | getattr("to_upper") | call`, object.NewString("ABC")},
		{`"abc" | func(s) { s.to_upper() }`, object.NewString("ABC")},
		{`[11, 12, 3] | math.sum`, object.NewFloat(26)},
		{`"42" | json.unmarshal`, object.NewFloat(42)},
	}
	runTests(t, tests)
}

func TestQuicksort(t *testing.T) {
	result, err := run(context.Background(), `
	func quicksort(arr) {
		if len(arr) < 2 {
			return arr
		} else {
			pivot := arr[0]
			less := arr[1:].filter(func(x) { x <= pivot })
			more := arr[1:].filter(func(x) { x > pivot })
			return quicksort(less) + [pivot] + quicksort(more)
		}
	}
	quicksort([10, 5, 2, 3])
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList(
		[]object.Object{
			object.NewInt(2),
			object.NewInt(3),
			object.NewInt(5),
			object.NewInt(10),
		}), result)
}

func TestMergesort(t *testing.T) {
	result, err := run(context.Background(), `
	func mergesort(arr) {
		length := len(arr)
		if length <= 1 {
			return arr
		}
		mid := length / 2
		left := mergesort(arr[:mid])
		right := mergesort(arr[mid:])
		output := list(length)
		i, j, k := [0, 0, 0]
		for i < len(left) {
			for j < len(right) && right[j] <= left[i] {
				output[k] = right[j]
				k++
				j++
			}
			output[k] = left[i]
			k++
			i++
		}
		for j < len(right) {
			output[k] = right[j]
			k++
			j++
		}
		return output
	}
	", ".join(mergesort([1, 9, -1, 4, 3, 2, 7, 8, 5, 6, 0]).map(string))
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9"), result)
}

func TestRecursiveIsPrime(t *testing.T) {
	result, err := run(context.Background(), `
	func is_prime(n, i=2) {
		// Base cases
		if (n <= 2) { return n == 2 }
		if (n % i == 0) { return false }
		if (i * i > n) { return true }
		// Check for next divisor
    	return is_prime(n, i + 1);
	}
	ints := []
	for i := 1; i < 30; i++ { ints.append(i) }
	primes := ints.filter(is_prime)
	", ".join(primes.map(string))
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("2, 3, 5, 7, 11, 13, 17, 19, 23, 29"), result)
}

func TestAndShortCircuit(t *testing.T) {
	// AND should short-circuit, so data[5] should not be evaluated
	result, err := run(context.Background(), `
	data := []
	if len(data) && data[5] {
		"nope!"
	} else {
		"worked!"
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("worked!"), result)
}

func TestOrShortCircuit(t *testing.T) {
	// OR should short-circuit, so data[5] should not be evaluated
	result, err := run(context.Background(), `
	data := [1]
	if len(data) || data[5] {
		"worked!"
	} else {
		"nope!"
	}
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewString("worked!"), result)
}

func TestLoopBreak(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for i := 0; i < 10; i++ {
		if i == 3 { break }
		x = i
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestLoopContinue(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for i := 0; i < 10; i++ {
		if i > 3 { continue }
		x = i
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestRangeLoopBreak(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for i := range [0, 1, 2, 3, 4] {
		if i == 3 { break }
		x = i
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestRangeLoopContinue(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for i := range [0, 1, 2, 3, 4] {
		if i > 3 { continue }
		x = i
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestSimpleLoopBreak(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	for {
		x++
		if x == 2 { break }
		max := math.max(1, 2) // inject some extra instructions
	}
	x
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(2), result)
}

func TestNestedLoops(t *testing.T) {
	result, err := run(context.Background(), `
	x, y, z := [0, 0, 0]
	for {
		x++ // This should execute 3 times total
		if x == 3 { break }
		// We should reach this point twice, with x as 1 then 2
		for i := range [0, 1, 2, 3] {
			y++ // This should execute 8 times total
			if i > 1 { continue }
			// We should reach this point 4 times total
			for h := 0; h < 10; h++ {
				z++ // This should execute 16 times total (4 times per inner loop)
				if h == 3 { break }
			}
		}
	}
	[x, y, z]
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(3),
		object.NewInt(8),
		object.NewInt(16),
	}), result)
}

func TestSimpleLoopContinue(t *testing.T) {
	result, err := run(context.Background(), `
	x := 0
	y := 0
	for {
		x++
		if x < 2 { continue }
		// We'll reach here on x in [2, 3, 4, 5, 6]
		if x > 5 { break }
		// We'll reach here on x in [2, 3, 4, 5]; so y should increment 4 times
		y++
		max := math.max(1, 2) // inject some extra instructions
	}
	y
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestManyLocals(t *testing.T) {
	result, err := run(context.Background(), `
	func example(x) {
		a := x + 1
		b := a + 1
		c := b + 1
		d := c + 1
		e := d + 1
		f := e + 1
		g := f + 1
		h := g + 1
		i := h + 1
		j := i + 1
		k := j + 1
		l := k + 1
		return l
	}
	example(0)
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(12), result)
}

func TestIncorrectArgCount(t *testing.T) {
	type testCase struct {
		input       string
		expectedErr string
	}
	tests := []testCase{
		{`func ex() { 1 }; ex(1)`, "type error: function takes no arguments (1 given)"},
		{`func ex(x) { x }; ex()`, "type error: function takes 1 argument (0 given)"},
		{`func ex(x) { x }; ex(1, 2)`, "type error: function takes 1 argument (2 given)"},
		{`func ex(x, y) { 1 }; ex()`, "type error: function takes 2 arguments (0 given)"},
		{`func ex(x, y) { 1 }; ex(0)`, "type error: function takes 2 arguments (1 given)"},
		{`func ex(x, y) { 1 }; ex(1, 2, 3)`, "type error: function takes 2 arguments (3 given)"},
		{`func ex() { 1 }; [1, 2].filter(ex)`, "type error: function takes no arguments (1 given)"},
		{`func ex() { 1 }; "foo" | ex`, "type error: function takes no arguments (1 given)"},
		{`"foo" | "bar"`, "type error: object is not callable (got string)"},
	}
	for _, tt := range tests {
		_, err := run(context.Background(), tt.input)
		require.NotNil(t, err)
		require.Equal(t, tt.expectedErr, err.Error())
	}
}

type testData struct {
	Count int
}

func (t *testData) Increment() {
	t.Count++
}

func (t testData) GetCount() int {
	return t.Count
}

type testStruct struct {
	A int
	B string
	C *testData
}

func TestNestedProxies(t *testing.T) {
	s := &testStruct{
		A: 1,
		B: "foo",
		C: &testData{
			Count: 3,
		},
	}
	opts := runOpts{
		Globals: map[string]interface{}{"s": s},
	}
	result, err := run(context.Background(), `
	s.C.Increment()
	s.C.GetCount()
	`, opts)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestProxy(t *testing.T) {
	type test struct {
		Data []byte
	}
	opts := runOpts{
		Globals: map[string]interface{}{
			"s": &test{Data: []byte("foo")},
		},
	}
	result, err := run(context.Background(), `s.Data`, opts)
	require.Nil(t, err)
	require.Equal(t, object.NewByteSlice([]byte("foo")), result)
}

func TestHalt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	_, err := run(ctx, `for {}`)
	require.NotNil(t, err)
	require.Equal(t, context.DeadlineExceeded, err)
}

func TestReturnGlobalVariable(t *testing.T) {
	result, err := run(context.Background(), `
	x := 3
	func test() { x }
	test()
	`)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), result)
}

func TestNakedReturn(t *testing.T) {
	result, err := run(context.Background(), `func test(a) { return }; test(15)`)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestGlobalNames(t *testing.T) {
	ctx := context.Background()
	source := `
	count := 1
	func inc(a, b) { a + b }
	m := {one: 1}
	foo := func() { "bar" }
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))

	globals := vm.GlobalNames()
	globalsMap := map[string]bool{}
	for _, g := range globals {
		globalsMap[g] = true
	}
	require.True(t, globalsMap["count"])
	require.True(t, globalsMap["inc"])
	require.True(t, globalsMap["m"])
	require.True(t, globalsMap["foo"])
}

func TestGetGlobal(t *testing.T) {
	ctx := context.Background()
	source := `func inc(a, b) { a + b }`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))

	obj, err := vm.Get("inc")
	require.Nil(t, err)
	fn, ok := obj.(*object.Function)
	require.True(t, ok)
	require.Equal(t, "inc", fn.Name())
}

func TestCall(t *testing.T) {
	ctx := context.Background()
	source := `func inc(a, b) { a + b }`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))

	obj, err := vm.Get("inc")
	require.Nil(t, err)
	fn, ok := obj.(*object.Function)
	require.True(t, ok)

	result, err := vm.Call(ctx, fn, []object.Object{
		object.NewInt(9),
		object.NewInt(1),
	})
	require.Nil(t, err)
	require.Equal(t, object.NewInt(10), result)
}

func TestCallWithClosure(t *testing.T) {
	ctx := context.Background()
	source := `
	func get_counter() {
		count := 10
		return func() {
			count++
			return count
		}
	}
	counter := get_counter()
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))

	obj, err := vm.Get("counter")
	require.Nil(t, err)
	counter, ok := obj.(*object.Function)
	require.True(t, ok)

	// The counter's first value will be 11. Confirm it counts up from there.
	for i := int64(11); i < 100; i++ {
		obj, err := vm.Call(ctx, counter, []object.Object{})
		require.Nil(t, err)
		require.Equal(t, object.NewInt(i), obj)
	}
}

func TestFreeVariableAssignment(t *testing.T) {
	ctx := context.Background()
	source := `
	func get_counters() {
		a := 0
		b := 0
		c := 0
		func incA() {
			a++
			return a
		}
		func incB() {
			b++
			return b
		}
		func incC() {
			c++
			return c
		}
		return [incA, incB, incC]
	}
	incA, incB, incC := get_counters()
	incA(); incA()                 // 1, 2
	incB(); incB(); incB()         // 1, 2, 3
	incC(); incC(); incC(); incC() // 1, 2, 3, 4
	[incA(), incB(), incC()]       // [3, 4, 5]
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))
	result, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(3),
		object.NewInt(4),
		object.NewInt(5),
	}), result)
}

func TestInterpolatedStringClosures1(t *testing.T) {
	ctx := context.Background()
	source := `
	func foo(a, b, c) {
		return func(d) {
			return '{strings.to_upper(a)}-{b}-{c}-{d}'
		}
	}
	foo("foo", "bar", "baz")("go")
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))
	result, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewString("FOO-bar-baz-go"), result)
}

func TestInterpolatedStringClosures2(t *testing.T) {
	ctx := context.Background()
	source := `
	x := 3
	func foo(a, b="bar") {
		count := 42
		return func(a) {
			return 'a: {a} b: {b} count: {count-2} x: {x+1}'
		}
	}
	foo("IGNORED")("HEY")
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))
	result, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewString("a: HEY b: bar count: 40 x: 4"), result)
}

func TestClone(t *testing.T) {
	ctx := context.Background()
	source := `
	x := 3
	func inc() {
		x++
	}
	inc()
	x
	`
	vm, err := newVM(ctx, source)
	require.Nil(t, err)
	require.Nil(t, vm.Run(ctx))
	result, ok := vm.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(4), result)

	clone, err := vm.Clone()
	require.Nil(t, err)
	value, err := clone.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), value)
}

func TestCloneWithAnonymousFunc(t *testing.T) {
	registered := map[string]*object.Function{}

	// Custom built-in function to be called from the Risor script to register
	// an anonymous function
	registerFunc := func(ctx context.Context, args ...object.Object) object.Object {
		name := args[0].(*object.String).Value()
		fn := args[1].(*object.Function)
		registered[name] = fn
		return object.Nil
	}

	ctx := context.Background()
	source := `
	x := 3
	register("inc", func() {
		x++
		return x
	})
	`
	globals := map[string]any{
		"register": object.NewBuiltin("register", registerFunc),
	}
	machine, err := newVM(ctx, source, runOpts{Globals: globals})
	require.Nil(t, err)
	require.Nil(t, machine.Run(ctx))

	// x should be 3 in the original VM
	value, err := machine.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), value)

	// Confirm the "inc" function was registered
	incFunc, ok := registered["inc"]
	require.True(t, ok)

	// Create a clone of the VM and confirm it also has x = 3
	clone, err := machine.Clone()
	require.Nil(t, err)
	value, err = clone.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), value)

	// Call the "inc" function in the clone and confirm it increments x to 4
	// in both the clone and the original VM
	_, err = clone.Call(ctx, incFunc, nil)
	require.Nil(t, err)

	// Clone's x is now 4
	value, err = clone.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), value)

	// Original's x is now 4
	value, err = machine.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), value)
}

func TestIncrementalEvaluation(t *testing.T) {
	ctx := context.Background()
	ast, err := parser.Parse(ctx, "x := 3")
	require.Nil(t, err)

	comp, err := compiler.New()
	require.Nil(t, err)
	main, err := comp.Compile(ast)
	require.Nil(t, err)

	v := New(main)
	require.Nil(t, v.Run(ctx))
	value, err := v.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), value)

	ast, err = parser.Parse(ctx, "x + 7")
	require.Nil(t, err)
	_, err = comp.Compile(ast)
	require.Nil(t, err)
	require.Nil(t, v.Run(ctx))
	value, err = v.Get("x")
	require.Nil(t, err)
	require.Equal(t, object.NewInt(3), value)

	tos, ok := v.TOS()
	require.True(t, ok)
	require.Equal(t, object.NewInt(10), tos)
}

func TestImports(t *testing.T) {
	tests := []testCase{
		{`import simple_math; simple_math.add(3, 4)`, object.NewInt(7)},
		{`import simple_math; int(simple_math.pi)`, object.NewInt(3)},
		{`import data; data.mydata["count"]`, object.NewInt(1)},
		{`import data; data.mydata["count"] = 3; data.mydata["count"]`, object.NewInt(3)},
		{`import data as d; d.mydata["count"]`, object.NewInt(1)},
		{`import math as m; m.min(3,-7)`, object.NewFloat(-7)},
	}
	runTests(t, tests)
}

func TestFromImport(t *testing.T) {
	tests := []testCase{
		{`from a.data import mapValue; mapValue["3"]`, object.NewInt(3)},
		{`from a.function import plusOne; plusOne(1)`, object.NewInt(2)},
		{`from a import function; function.plusOne(1)`, object.NewInt(2)},
		{`from a.b import data as b_data; from a.function import plusOne; plusOne(b_data.mapValue["1"]) `, object.NewInt(2)},
		{`from math import min; min(3,-7)`, object.NewFloat(-7)},
		{`from math import min as m; m(3,-7)`, object.NewFloat(-7)},
		{
			`from math import (min as a, max as b); [a(1,2), b(1,2)]`,
			object.NewList([]object.Object{
				object.NewFloat(1),
				object.NewFloat(2),
			}),
		},
	}
	runTests(t, tests)
}

func TestBadImports(t *testing.T) {
	ctx := context.Background()
	type testCase struct {
		input     string
		expectErr string
	}
	tests := []testCase{
		{`import foo`, `import error: module "foo" not found`},
		{`import foo as bar`, `import error: module "foo" not found`},
		{`import math as`, `parse error: unexpected end of file while parsing an import statement (expected identifier)`},
		{`from foo import bar`, `import error: module "foo" not found`},
		{`from a.b import c`, `import error: module "a/b" not found`},
		{`from a.b import c as d`, `import error: module "a/b" not found`},
		{`from math import foo`, `import error: cannot import name "foo" from "math"`},
		{`from math`, `parse error: from-import is missing import statement`},
		{`from math import`, `parse error: unexpected end of file while parsing a from-import statement (expected identifier)`},
		{`from math import min as`, `parse error: unexpected end of file while parsing a from-import statement (expected identifier)`},
	}
	for _, tt := range tests {
		_, err := run(ctx, tt.input)
		require.NotNil(t, err)
		require.Equal(t, tt.expectErr, err.Error())
	}
}

func TestModifyModule(t *testing.T) {
	_, err := run(context.Background(), `math.max = 123`)
	require.Error(t, err)
	require.Equal(t, "attribute error: cannot modify module attributes", err.Error())
}

func TestEarlyForRangeReturn(t *testing.T) {
	code := `
func operation(c) {
	for range [1, 2, 3] {
		return "result"
	}
}
func main() {
	items := ['ab', 'cd']
	results := []
	for _, item := range items {
		value := operation(item)
		results.append(value)
	}
	return results
}
main()
`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("result"),
		object.NewString("result"),
	}), result)
}

func TestDeferStatementGlobalClosure(t *testing.T) {
	code := `
	x := 0
	func foo(value) { defer func() { x = value }() }
	foo(4)
	x
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(4), result)
}

func TestDeferStatementOrdering(t *testing.T) {
	code := `
	l := []
	func foo(value) {
		defer l.append(value+1) // 3
		defer l.append(value)   // 2
		l.append(1)             // 1
	}
	foo(2)
	l
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
	}), result)
}

func TestDeferStatementAnon(t *testing.T) {
	code := `
	func() {
		x := 42
		defer func() { x = 1 }()
		defer func() { x = 2 }()
		return x
	}()
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewInt(42), result)
}

func TestDeferStatementBuiltin(t *testing.T) {
	code := `
	m := {one: 1, two: 2}
	func test() {
		func() {
			m["three"] = 3
			defer delete(m, "one")
		}()
	}
	test()
	m
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewMap(map[string]object.Object{
		"two":   object.NewInt(2),
		"three": object.NewInt(3),
	}), result)
}

func TestDeferFileClose(t *testing.T) {
	code := `
	func get_lines(path) {
		f := os.open(path)
		defer f.close()
		return string(f.read()).split("\n")
	}
	lines := get_lines("fixtures/jabberwocky.txt")
	lines[0]
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewString("'Twas brillig, and the slithy toves"), result)
}

func TestDeferBehavior(t *testing.T) {
	code := `
	func work(count) {
		c := chan(count)
		spawn(func() {
			defer close(c)
			for i := 0; i < count; i++ {
				c <- i
			}
		})
		return c
	}
	results := []
	for _, value := range work(5) {
		results.append(value)
	}
	results
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(0),
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
		object.NewInt(4),
	}), result)
}

func TestFreeVariables(t *testing.T) {
	code := `
	func test(count) {
		l := []
		func() {
			y := count
			if true {
				l.append(y)
			}
		}()
		return l
	}
	test(5)
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.NewList([]object.Object{object.NewInt(5)}), result)
}

func TestChannels(t *testing.T) {
	tests := []testCase{
		{`c := chan(1); c <- 1; <-c`, object.NewInt(1)},
		{`c := make(chan, 1); c <- 1; <-c`, object.NewInt(1)},
		{`c := chan(2); c <- 1; c <- 2; [<-c, <-c]`, object.NewList([]object.Object{
			object.NewInt(1), object.NewInt(2),
		})},
		{`c := chan(1); c <- 1; close(c); <-c`, object.NewInt(1)},
		{`c := chan(); close(c); <-c`, object.Nil},
		{`c := chan(1); c <- "ok"; close(c); [<-c, <-c]`, object.NewList([]object.Object{
			object.NewString("ok"), object.Nil,
		})},
		{`c := chan(2); c <- "a"; c <- "b"; close(c);
		  results := []
		  for _, value := range c { results.append(value) }
		  results`, object.NewList([]object.Object{
			object.NewString("a"),
			object.NewString("b"),
		})},
	}
	runTests(t, tests)
}

func TestChannelErrors(t *testing.T) {
	ctx := context.Background()
	type testCase struct {
		input     string
		expectErr string
	}
	tests := []testCase{
		{`c := chan(1); close(c); c <- 1`, "exec error: send on closed channel"},
		{`c := chan(1); close(c); close(c)`, "exec error: close of closed channel"},
		{`c := chan(1); c <- 1; close(c); c <- 2`, "exec error: send on closed channel"},
	}
	for _, tt := range tests {
		_, err := run(ctx, tt.input)
		require.NotNil(t, err)
		require.Equal(t, tt.expectErr, err.Error())
	}
}

func TestGoStatement(t *testing.T) {
	tests := []testCase{
		{`go func() { 1 }()`, object.Nil},
		{`x := 0; go func() { x = 1 }(); time.sleep(0.1); x`, object.NewInt(1)},
		{`x := 0; go func() { x = 1 }(); x`, object.NewInt(0)},
		{`c := chan(1); go func() { c <- 1 }(); <-c`, object.NewInt(1)},
		{`func dowork() {
			c := make(chan)
			go func() { defer close(c); c <- 98765 }();
			return c
		  }
		  rxchan := dowork()
		  <-rxchan`, object.NewInt(98765)},
	}
	runTests(t, tests)
}

func TestSpawn(t *testing.T) {
	tests := []testCase{
		{`func test(x) { return x + 1 }; spawn(test, 33).wait()`, object.NewInt(34)},
		{`spawn(func(x=10) { x }).wait()`, object.NewInt(10)},
		{`x := 0; spawn(func() { x = 34 }).wait(); x`, object.NewInt(34)},
		{`l := []; spawn(func() { l.append(1) }).wait(); l`, object.NewList([]object.Object{object.NewInt(1)})},
		{`l := []; spawn(func(x) { x.append(1) }, l).wait(); l`, object.NewList([]object.Object{object.NewInt(1)})},
		{`
		func work(x) { return x ** 2 }
		threads := []
		for i := 0; i < 5; i++ { threads.append(spawn(work, i))	}
		threads.map(func(t) { t.wait() })
		`, object.NewList([]object.Object{
			object.NewInt(0),
			object.NewInt(1),
			object.NewInt(4),
			object.NewInt(9),
			object.NewInt(16),
		})},
	}
	runTests(t, tests)
}

func TestMaps(t *testing.T) {
	tests := []testCase{
		{`{"a": 1}`, object.NewMap(map[string]object.Object{
			"a": object.NewInt(1),
		})},
		{`{"a": 1,}`, object.NewMap(map[string]object.Object{
			"a": object.NewInt(1),
		})},
		{`{"a": 1,
		  }`, object.NewMap(map[string]object.Object{
			"a": object.NewInt(1),
		})},
		{`{"a": 1,
		   "b": 2}`, object.NewMap(map[string]object.Object{
			"a": object.NewInt(1),
			"b": object.NewInt(2),
		})},
		{`{"a": 1,
			"b": 2
		}`, object.NewMap(map[string]object.Object{
			"a": object.NewInt(1),
			"b": object.NewInt(2),
		})},
	}
	runTests(t, tests)
}

func TestLists(t *testing.T) {
	tests := []testCase{
		{`[1,2,3]`, object.NewList([]object.Object{
			object.NewInt(1),
			object.NewInt(2),
			object.NewInt(3),
		})},
		{`[1,
		   2,
		   3]`, object.NewList([]object.Object{
			object.NewInt(1),
			object.NewInt(2),
			object.NewInt(3),
		})},
		{`[1,
		   2,]`, object.NewList([]object.Object{
			object.NewInt(1),
			object.NewInt(2),
		})},
		{`[1,
		2
		]`, object.NewList([]object.Object{
			object.NewInt(1),
			object.NewInt(2),
		})},
	}
	runTests(t, tests)
}

func TestFunctionStack(t *testing.T) {
	code := `
	for i := range 1 {
		try(func() {
		  42
		  error("kaboom")
		})
	  }
	`
	result, err := run(context.Background(), code)
	require.Nil(t, err)
	require.Equal(t, object.Nil, result)
}

func TestContextDone(t *testing.T) {
	// Context with no deadline does not return a Done channel
	ctx := context.Background()
	d := ctx.Done()
	require.Nil(t, d)

	// Context with deadline returns a Done channel
	tctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	d = tctx.Done()
	require.NotNil(t, d)

	// Context with cancel returns a Done channel
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d = cctx.Done()
	require.NotNil(t, d)
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
			result, err := run(ctx, tt.input)
			require.Nil(t, err)
			require.NotNil(t, result)
			require.Equal(t, tt.expected, result)
		})
	}
}
