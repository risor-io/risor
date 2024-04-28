package isatty

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	fnObj, ok := m.GetAttr("is_terminal")
	require.True(t, ok)

	fn, ok := fnObj.(*object.Builtin)
	require.True(t, ok)
	result := fn.Call(context.Background(), object.NewInt(1)) // fd 1 is stdout
	require.NotNil(t, result)
	_, ok = result.(*object.Bool)
	require.True(t, ok)

	fnObj, ok = m.GetAttr("is_cygwin_terminal")
	require.True(t, ok)

	fn, ok = fnObj.(*object.Builtin)
	require.True(t, ok)
	result = fn.Call(context.Background(), object.NewInt(1)) // fd 1 is stdout
	require.NotNil(t, result)
	_, ok = result.(*object.Bool)
	require.True(t, ok)
}
