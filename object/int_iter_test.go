package object

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntIterPositive(t *testing.T) {
	iter := NewIntIter(NewInt(3))
	ctx := context.Background()
	var entries []Object
	for {
		entry, ok := iter.Next(ctx)
		if !ok {
			break
		}
		entries = append(entries, entry)
	}
	require.Len(t, entries, 3)
	require.Equal(t, NewInt(0), entries[0])
	require.Equal(t, NewInt(1), entries[1])
	require.Equal(t, NewInt(2), entries[2])
}

func TestIntIterNegative(t *testing.T) {
	iter := NewIntIter(NewInt(-3))
	ctx := context.Background()
	var entries []Object
	for {
		entry, ok := iter.Next(ctx)
		if !ok {
			break
		}
		entries = append(entries, entry)
	}
	require.Len(t, entries, 3)
	require.Equal(t, NewInt(0), entries[0])
	require.Equal(t, NewInt(-1), entries[1])
	require.Equal(t, NewInt(-2), entries[2])
}
