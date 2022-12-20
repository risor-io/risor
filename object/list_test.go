package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListInsert(t *testing.T) {
	one := NewInt(1)
	two := NewInt(2)
	thr := NewInt(3)

	list := &List{Items: []Object{one}}

	list.Insert(5, two)
	require.Equal(t, []Object{one, two}, list.Items)

	list.Insert(-10, thr)
	require.Equal(t, []Object{thr, one, two}, list.Items)

	list.Insert(1, two)
	require.Equal(t, []Object{thr, two, one, two}, list.Items)

	list.Insert(0, two)
	require.Equal(t, []Object{two, thr, two, one, two}, list.Items)
}

func TestListPop(t *testing.T) {
	zero := NewString("0")
	one := NewString("1")
	two := NewString("2")

	list := &List{Items: []Object{zero, one, two}}

	val, ok := list.Pop(1).(*String)
	require.True(t, ok)
	require.Equal(t, "1", val.Value)

	val, ok = list.Pop(1).(*String)
	require.True(t, ok)
	require.Equal(t, "2", val.Value)

	err, ok := list.Pop(1).(*Error)
	require.True(t, ok)
	require.Equal(t, "index error: index out of range: 1", err.Message)
}
