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
