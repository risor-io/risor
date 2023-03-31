package symbol

import (
	"testing"

	"github.com/cloudcmds/tamarin/object"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	table := NewTable()

	require.Nil(t, table.Parent())
	require.Equal(t, uint16(0), table.Size())

	a, err := table.InsertVariable("a")
	require.Nil(t, err)
	require.Equal(t, uint16(0), a.Index)
	require.Equal(t, "a", a.Name)
	require.Nil(t, a.Value)

	b, err := table.InsertVariable("b")
	require.Nil(t, err)
	require.Equal(t, uint16(1), b.Index)
	require.Equal(t, "b", b.Name)
	require.Nil(t, b.Value)

	c, err := table.InsertBuiltin("c")
	require.Nil(t, err)
	require.Equal(t, uint16(0), c.Index)
	require.Equal(t, "c", c.Name)
	require.Nil(t, c.Value)

	// The size is the count of variables, not including builtins
	require.Equal(t, uint16(2), table.Size())

	require.True(t, table.IsVariable("a"))
	require.True(t, table.IsVariable("b"))
	require.False(t, table.IsVariable("c"))

	require.False(t, table.IsBuiltin("a"))
	require.False(t, table.IsBuiltin("b"))
	require.True(t, table.IsBuiltin("c"))

	require.Len(t, table.Variables(), 2)
	require.Len(t, table.Builtins(), 1)
}

func TestBlock(t *testing.T) {
	table := NewTable()
	block := table.NewBlock()

	block.InsertVariable("a", object.NewInt(42))

	require.Equal(t, uint16(1), table.Size())

	locals := table.Variables()
	require.Len(t, locals, 1)
	value := locals[0]
	require.Equal(t, object.NewInt(42), value)
}
