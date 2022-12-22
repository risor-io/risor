package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringHashKey(t *testing.T) {
	a := NewString("hello")
	b := NewString("hello")
	c := NewString("goodbye")
	d := NewString("goodbye")

	require.Equal(t, a.HashKey(), b.HashKey())
	require.Equal(t, c.HashKey(), d.HashKey())
	require.NotEqual(t, a.HashKey(), c.HashKey())

	require.Equal(t, HashKey{Type: STRING, StrValue: "hello"}, a.HashKey())
}
