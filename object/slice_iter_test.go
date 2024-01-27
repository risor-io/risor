package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceIter(t *testing.T) {
	iter, err := NewSliceIter([]int{5, 6})
	require.Nil(t, err)

	require.True(t, iter.IsTruthy())
	require.Equal(t, "slice_iter(pos=-1 size=2)", iter.Inspect())

	// Call Next to go to position 0 (value 5)
	obj, ok := iter.Next()
	require.True(t, ok)
	require.Equal(t, int64(5), obj.(*Int).value)

	entryObj, ok := iter.Entry()
	require.True(t, ok)
	entry := entryObj.(*Entry)
	require.Equal(t, NewInt(0), entry.Key())
	require.Equal(t, NewInt(5), entry.Value())

	// Call Next to go to position 1 (value 6)
	obj, ok = iter.Next()
	require.True(t, ok)
	require.Equal(t, int64(6), obj.(*Int).value)

	entryObj, ok = iter.Entry()
	require.True(t, ok)
	entry = entryObj.(*Entry)
	require.Equal(t, NewInt(1), entry.Key())
	require.Equal(t, NewInt(6), entry.Value())

	// We should be at the end now
	_, ok = iter.Next()
	require.False(t, ok)
}

func TestSliceIterStrings(t *testing.T) {
	iter, err := NewSliceIter([]string{"apple", "banana"})
	require.Nil(t, err)

	obj, ok := iter.Next()
	require.True(t, ok)
	require.Equal(t, "apple", obj.(*String).value)

	entryObj, ok := iter.Entry()
	require.True(t, ok)
	entry := entryObj.(*Entry)
	require.Equal(t, NewInt(0), entry.Key())
	require.Equal(t, NewString("apple"), entry.Value())

	obj, ok = iter.Next()
	require.True(t, ok)
	require.Equal(t, "banana", obj.(*String).value)

	entryObj, ok = iter.Entry()
	require.True(t, ok)
	entry = entryObj.(*Entry)
	require.Equal(t, NewInt(1), entry.Key())
	require.Equal(t, NewString("banana"), entry.Value())

	_, ok = iter.Next()
	require.False(t, ok)
}
