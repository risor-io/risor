package fmt

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestPrintableValue(t *testing.T) {
	type testCase struct {
		obj      object.Object
		expected any
	}

	testTime, err := time.Parse("2006-01-02", "2021-01-01")
	require.NoError(t, err)

	builtin := func(ctx context.Context, args ...object.Object) object.Object {
		return nil
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
		{obj: object.NewTime(testTime), expected: "2021-01-01T00:00:00Z"},
		{obj: object.NewBuiltin("foo", builtin), expected: "builtin(foo)"},
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
			got := object.PrintableValue(tc.obj)
			require.Equal(t, tc.expected, got)
		})
	}
}
