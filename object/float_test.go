package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatBasics(t *testing.T) {
	value := NewFloat(-2)
	require.Equal(t, FLOAT, value.Type())
	require.Equal(t, float64(-2), value.Value())
	require.Equal(t, "float(-2)", value.String())
	require.Equal(t, "-2", value.Inspect())
	require.Equal(t, float64(-2), value.Interface())
}
