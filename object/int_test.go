package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntCompare(t *testing.T) {

	one := NewInt(1)
	two := NewFloat(2.0)
	thr := NewInt(3)

	tests := []struct {
		first    Comparable
		second   Object
		expected int
	}{
		{one, two, -1},
		{two, one, 1},
		{one, one, 0},
		{two, thr, -1},
		{thr, two, 1},
		{two, two, 0},
	}
	for _, tc := range tests {
		result, err := tc.first.Compare(tc.second)
		require.Nil(t, err)
		require.Equal(t, tc.expected, result,
			"first: %v, second: %v", tc.first, tc.second)
	}
}

func TestIntEquals(t *testing.T) {

	oneInt := NewInt(1)
	twoFlt := NewFloat(2.0)
	twoInt := NewInt(2)

	tests := []struct {
		first    Object
		second   Object
		expected bool
	}{
		{oneInt, twoFlt, false},
		{oneInt, twoInt, false},
		{oneInt, oneInt, true},
		{twoInt, twoFlt, true},
		{twoFlt, twoFlt, true},
	}
	for _, tc := range tests {
		result, ok := tc.first.Equals(tc.second).(*Bool)
		require.True(t, ok)
		require.Equal(t, tc.expected, result.Value(),
			"first: %v, second: %v", tc.first, tc.second)
	}
}

func TestIntBasics(t *testing.T) {
	value := NewInt(-3)
	require.Equal(t, INT, value.Type())
	require.Equal(t, int64(-3), value.Value())
	require.Equal(t, "int(-3)", value.String())
	require.Equal(t, "-3", value.Inspect())
	require.Equal(t, int64(-3), value.Interface())
}
