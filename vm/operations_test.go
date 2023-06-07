package vm

import (
	"testing"

	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/op"
	"github.com/stretchr/testify/require"
)

func TestCompareNonComparable(t *testing.T) {
	s1 := object.NewSet(nil)
	s2 := object.NewSet(nil)
	result := compare(op.LessThan, s1, s2)
	resultErr, ok := result.(*object.Error)
	require.True(t, ok)
	require.Equal(t, object.Errorf("object is not comparable: *object.Set"), resultErr)
}

func TestCompareUnknownComparison(t *testing.T) {
	obj1 := object.NewInt(1)
	obj2 := object.NewInt(2)
	require.Panics(t, func() {
		compare(op.CompareOpType(op.Halt), obj1, obj2)
	})
}

func TestAndOperator(t *testing.T) {
	type testCase struct {
		left  object.Object
		right object.Object
		want  object.Object
	}
	testCases := []testCase{
		{object.NewInt(1), object.NewInt(1), object.NewInt(1)},
		{object.NewInt(1), object.NewInt(0), object.NewInt(0)},
		{object.NewInt(0), object.NewInt(1), object.NewInt(0)},
		{object.NewInt(0), object.NewInt(0), object.NewInt(0)},
		{object.NewInt(1), object.NewBool(true), object.NewBool(true)},
		{object.NewInt(1), object.NewBool(false), object.NewBool(false)},
		{object.NewInt(0), object.NewBool(true), object.NewInt(0)},
		{object.NewInt(0), object.NewBool(false), object.NewInt(0)},
		{object.NewBool(true), object.NewInt(1), object.NewInt(1)},
		{object.NewBool(true), object.NewInt(0), object.NewInt(0)},
	}
	for _, tc := range testCases {
		result := binaryOp(op.And, tc.left, tc.right)
		require.Equal(t, tc.want, result)
	}
}

func TestOrOperator(t *testing.T) {
	type testCase struct {
		left  object.Object
		right object.Object
		want  object.Object
	}
	testCases := []testCase{
		{object.NewInt(1), object.NewInt(1), object.NewInt(1)},
		{object.NewInt(1), object.NewInt(0), object.NewInt(1)},
		{object.NewInt(0), object.NewInt(1), object.NewInt(1)},
		{object.NewInt(0), object.NewInt(0), object.NewInt(0)},
		{object.NewInt(1), object.NewBool(true), object.NewInt(1)},
		{object.NewInt(1), object.NewBool(false), object.NewInt(1)},
		{object.NewInt(0), object.NewBool(true), object.NewBool(true)},
		{object.NewInt(0), object.NewBool(false), object.NewBool(false)},
		{object.NewBool(true), object.NewInt(1), object.NewBool(true)},
		{object.NewBool(true), object.NewInt(0), object.NewBool(true)},
	}
	for _, tc := range testCases {
		result := binaryOp(op.Or, tc.left, tc.right)
		require.Equal(t, tc.want, result)
	}
}
