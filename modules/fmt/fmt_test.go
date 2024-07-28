package fmt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestPrintableValue(t *testing.T) {
	type testCase struct {
		obj      object.Object
		expected any
	}
	cases := []testCase{
		{object.NewString("hello"), "hello"},
		{object.NewByte(5), byte(5)},
		{object.NewInt(42), int64(42)},
		{object.NewFloat(42.42), 42.42},
		{object.NewBool(true), true},
		{object.NewBool(false), false},
		{object.Errorf("error"), errors.New("error")},
		{obj: object.Nil, expected: nil},
		{ // strings printed inside lists are quoted in Risor
			obj: object.NewList([]object.Object{
				object.NewString("hello"),
				object.NewInt(42),
			}),
			expected: `["hello", 42]`,
		},
		{ // strings printed inside maps are quoted in Risor
			obj: object.NewMap(map[string]object.Object{
				"a": object.NewInt(42),
				"b": object.NewString("hello"),
				"c": object.Nil,
			}),
			expected: `{"a": 42, "b": "hello", "c": nil}`,
		},
		{
			obj: object.NewSet([]object.Object{
				object.NewInt(42),
				object.NewString("hi there"),
			}),
			expected: `{42, "hi there"}`,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.expected), func(t *testing.T) {
			got := printableValue(tc.obj)
			require.Equal(t, tc.expected, got)
		})
	}
}
