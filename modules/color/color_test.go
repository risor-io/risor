package color

import (
	"context"
	"testing"

	"github.com/fatih/color"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	fgRed, ok := m.GetAttr("fg_red")
	require.True(t, ok)

	fnObj, ok := m.GetAttr("color")
	require.True(t, ok)
	fn, ok := fnObj.(*object.Builtin)
	require.True(t, ok)

	result := fn.Call(context.Background(), fgRed)
	require.NotNil(t, result)
	colorObj, ok := result.(*Color)
	require.True(t, ok)

	c := colorObj.Value()
	require.False(t, c.Equals(color.New(color.FgBlue)))
	require.True(t, c.Equals(color.New(color.FgRed)))
}
