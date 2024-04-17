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

func TestSorted(t *testing.T) {
	ctx := context.Background()
	tests := []testCase{
		{
			object.NewList([]object.Object{
				object.NewInt(3),
				object.NewInt(1),
				object.NewInt(2),
			}),
			object.NewList([]object.Object{
				object.NewInt(1),
				object.NewInt(2),
				object.NewInt(3),
			}),
		},
		{
			object.NewList([]object.Object{
				object.NewInt(3),
				object.NewInt(1),
				object.NewString("nope"),
			}),
			object.Errorf("type error: unable to compare string and int"),
		},
		{
			object.NewList([]object.Object{
				object.NewString("b"),
				object.NewString("c"),
				object.NewString("a"),
			}),
			object.NewList([]object.Object{
				object.NewString("a"),
				object.NewString("b"),
				object.NewString("c"),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.input.Inspect(), func(t *testing.T) {
			result := Sorted(ctx, tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSortedWithFunc(t *testing.T) {
	ctx := context.Background()
	// We'll sort this list of integers
	input := object.NewList([]object.Object{
		object.NewInt(3),
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(99),
		object.NewInt(0),
	})
	// This function will be called for each comparison
	callFn := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		require.Len(t, args, 2)
		a := args[0].(*object.Int).Value()
		b := args[1].(*object.Int).Value()
		return object.NewBool(b < a), nil // descending order
	}
	ctx = object.WithCallFunc(ctx, callFn)

	// This sort function isn't actually used here in the test. This value
	// will be passed to callFn but we don't use it.
	var sortFn *object.Function

	// Confirm Sorted returns the expected sorted list
	result := Sorted(ctx, input, sortFn)
	require.Equal(t, object.NewList([]object.Object{
		object.NewInt(99),
		object.NewInt(3),
		object.NewInt(2),
		object.NewInt(1),
		object.NewInt(0),
	}), result)
}
