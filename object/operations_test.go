package object

import (
	"testing"

	"github.com/risor-io/risor/op"
	"github.com/stretchr/testify/require"
)

func TestCompareNonComparable(t *testing.T) {
	s1 := NewSet(nil)
	s2 := NewSet(nil)
	result := Compare(op.LessThan, s1, s2)
	resultErr, ok := result.(*Error)
	require.True(t, ok)
	require.Equal(t, Errorf("type error: expected a comparable object (got set)"), resultErr)
}

func TestCompareUnknownComparison(t *testing.T) {
	obj1 := NewInt(1)
	obj2 := NewInt(2)
	require.Panics(t, func() {
		Compare(op.CompareOpType(op.Halt), obj1, obj2)
	})
}

func TestAndOperator(t *testing.T) {
	type testCase struct {
		left  Object
		right Object
		want  Object
	}
	testCases := []testCase{
		{NewInt(1), NewInt(1), NewInt(1)},
		{NewInt(1), NewInt(0), NewInt(0)},
		{NewInt(0), NewInt(1), NewInt(0)},
		{NewInt(0), NewInt(0), NewInt(0)},
		{NewInt(1), NewBool(true), NewBool(true)},
		{NewInt(1), NewBool(false), NewBool(false)},
		{NewInt(0), NewBool(true), NewInt(0)},
		{NewInt(0), NewBool(false), NewInt(0)},
		{NewBool(true), NewInt(1), NewInt(1)},
		{NewBool(true), NewInt(0), NewInt(0)},
	}
	for _, tc := range testCases {
		result := BinaryOp(op.And, tc.left, tc.right)
		require.Equal(t, tc.want, result)
	}
}

func TestOrOperator(t *testing.T) {
	type testCase struct {
		left  Object
		right Object
		want  Object
	}
	testCases := []testCase{
		{NewInt(1), NewInt(1), NewInt(1)},
		{NewInt(1), NewInt(0), NewInt(1)},
		{NewInt(0), NewInt(1), NewInt(1)},
		{NewInt(0), NewInt(0), NewInt(0)},
		{NewInt(1), NewBool(true), NewInt(1)},
		{NewInt(1), NewBool(false), NewInt(1)},
		{NewInt(0), NewBool(true), NewBool(true)},
		{NewInt(0), NewBool(false), NewBool(false)},
		{NewBool(true), NewInt(1), NewBool(true)},
		{NewBool(true), NewInt(0), NewBool(true)},
	}
	for _, tc := range testCases {
		result := BinaryOp(op.Or, tc.left, tc.right)
		require.Equal(t, tc.want, result)
	}
}
