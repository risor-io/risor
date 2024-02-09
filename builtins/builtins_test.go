package builtins

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	input    object.Object
	expected object.Object
}

func TestMake(t *testing.T) {
	ctx := context.Background()

	tests := []testCase{
		{object.NewBuiltin("list", nil), object.NewList([]object.Object{})},
		{object.NewBuiltin("map", nil), object.NewMap(map[string]object.Object{})},
		{object.NewBuiltin("set", nil), object.NewSet([]object.Object{})},

		{object.NewList([]object.Object{
			object.NewString("ignored"),
		}), object.NewList([]object.Object{})},

		{object.NewMap(map[string]object.Object{
			"ignored": object.NewString("ignored"),
		}), object.NewMap(map[string]object.Object{})},

		{object.NewSet([]object.Object{
			object.NewString("ignored"),
		}), object.NewSet([]object.Object{})},
	}

	for _, tt := range tests {
		t.Run(tt.input.Inspect(), func(t *testing.T) {
			result := Make(ctx, tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestMakeChan(t *testing.T) {
	ctx := context.Background()
	result := Make(ctx, object.NewBuiltin("chan", nil), object.NewInt(4))
	require.IsType(t, &object.Chan{}, result)
	ch, _ := result.(*object.Chan)
	require.Equal(t, 4, ch.Capacity())
}
