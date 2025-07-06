package ssh

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	// Test missing arguments
	result := Connect(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 3 arguments")

	// Test invalid host type
	result = Connect(context.Background(), object.NewInt(123),
		object.NewString("user"), object.NewMap(map[string]object.Object{
			"password": object.NewString("pass"),
		}))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")
}

func TestModule(t *testing.T) {
	mod := Module()
	require.NotNil(t, mod)
	require.Equal(t, "ssh", mod.Name().Value())

	// Test that all expected functions are present
	functions := []string{"connect"}
	for _, fn := range functions {
		attr, found := mod.GetAttr(fn)
		require.True(t, found, "Function %s not found in module", fn)
		require.NotNil(t, attr, "Function %s is nil", fn)
	}
}
