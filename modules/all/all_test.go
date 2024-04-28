package all

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuiltins(t *testing.T) {
	b := Builtins()
	require.True(t, len(b) > 0, "Builtins() should return some builtins")
	require.Contains(t, b, "base64", "Builtins() should include base64")
}
