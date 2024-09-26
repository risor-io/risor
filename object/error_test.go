package object

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorEquals(t *testing.T) {
	e := NewError(errors.New("a"))
	other1 := NewError(errors.New("a"))
	other2 := NewError(errors.New("b"))

	require.Equal(t, "a", e.Message().Value())
	require.True(t, e.Equals(other1).Interface().(bool))
	require.False(t, e.Equals(other2).Interface().(bool))
}

func TestErrorCompareStr(t *testing.T) {
	e := NewError(errors.New("a"))
	other1 := NewError(errors.New("a"))
	other2 := NewError(errors.New("b"))

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

func TestErrorCompareRaised(t *testing.T) {
	a := NewError(errors.New("a")).WithRaised(true)
	b := NewError(errors.New("a")) // raised is set by default

	require.True(t, a.IsRaised())
	require.True(t, b.IsRaised())

	result, err := a.Compare(b)
	require.Nil(t, err)
	require.Equal(t, 0, result)

	b.WithRaised(false)
	require.False(t, b.IsRaised())

	result, err = a.Compare(b)
	require.Nil(t, err)
	require.Equal(t, 1, result)

	result, err = b.Compare(a)
	require.Nil(t, err)
	require.Equal(t, -1, result)
}

func TestErrorMessage(t *testing.T) {
	a := NewError(errors.New("a"))

	attr, ok := a.GetAttr("error")
	require.True(t, ok)
	fn := attr.(*Builtin)
	result := fn.Call(context.Background())
	require.Equal(t, "a", result.(*String).Value())

	attr, ok = a.GetAttr("message")
	require.True(t, ok)
	fn = attr.(*Builtin)
	result = fn.Call(context.Background())
	require.Equal(t, "a", result.(*String).Value())
}
