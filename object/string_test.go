package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringBasics(t *testing.T) {
	value := NewString("abcd")
	require.Equal(t, STRING, value.Type())
	require.Equal(t, "abcd", value.Value)
	require.Equal(t, "string(abcd)", value.String())
	require.Equal(t, `"abcd"`, value.Inspect())
	require.Equal(t, "abcd", value.Interface())
	require.True(t, value.Equals(NewString("abcd")).(*Bool).Value)
}

func TestStringCompare(t *testing.T) {
	a := NewString("a")
	b := NewString("b")
	A := NewString("A")
	tests := []struct {
		first    Comparable
		second   Object
		expected int
	}{
		{a, b, -1},
		{b, a, 1},
		{a, a, 0},
		{a, A, 1},
		{A, a, -1},
	}
	for _, tc := range tests {
		result, err := tc.first.Compare(tc.second)
		require.Nil(t, err)
		require.Equal(t, tc.expected, result,
			"first: %v, second: %v", tc.first, tc.second)
	}
}

func TestStringReverse(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"abc", "cba"},
	}
	for _, tc := range tests {
		result := NewString(tc.s).Reversed().Value
		require.Equal(t, tc.expected, result, "s: %v", tc.s)
	}
}

func TestStringGetItem(t *testing.T) {
	tests := []struct {
		s           string
		index       int64
		expected    string
		expectedErr string
	}{
		{"", 0, "", "index error: index out of range: 0"},
		{"a", 0, "a", ""},
		{"a", -1, "a", ""},
		{"a", -2, "a", "index error: index out of range: -2"},
		{"012345", 5, "5", ""},
		{"012345", -1, "5", ""},
		{"012345", -2, "4", ""},
	}
	for _, tc := range tests {
		msg := fmt.Sprintf("%v[%d]", tc.s, tc.index)
		result, err := NewString(tc.s).GetItem(NewInt(tc.index))
		if tc.expectedErr != "" {
			require.NotNil(t, err, msg)
			require.Equal(t, tc.expectedErr, err.Message, msg)
		} else {
			resultStr, ok := result.(*String)
			require.True(t, ok, msg)
			require.Equal(t, tc.expected, resultStr.Value, msg)
		}
	}
}
