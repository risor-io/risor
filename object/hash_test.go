package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringHashKey(t *testing.T) {
	a := &String{Value: "hello"}
	b := &String{Value: "hello"}
	c := &String{Value: "goodbye"}
	d := &String{Value: "goodbye"}

	require.Equal(t, a.HashKey(), b.HashKey())
	require.Equal(t, c.HashKey(), d.HashKey())
	require.NotEqual(t, a.HashKey(), c.HashKey())

	require.Equal(t, HashKey{Type: STRING, StrValue: "hello"}, a.HashKey())
}
