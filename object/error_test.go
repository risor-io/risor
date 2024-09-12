package object

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorString(t *testing.T) {
	e := NewError(errors.New("a"))
	other1 := NewError(errors.New("a"))
	other2 := NewError(errors.New("b"))

	require.Equal(t, "a", e.Message().Value())
	require.True(t, e.Equals(other1).Interface().(bool))
	require.False(t, e.Equals(other2).Interface().(bool))

	cmp, err := e.Compare(other1)
	require.Nil(t, err)
	require.Equal(t, 0, cmp)

	cmp, err = e.Compare(other2)
	require.Nil(t, err)
	require.Equal(t, -1, cmp)

	cmp, err = other2.Compare(e)
	require.Nil(t, err)
	require.Equal(t, 1, cmp)
}
