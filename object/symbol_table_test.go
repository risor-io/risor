package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	table := NewSymbolTable()

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
	table := NewSymbolTable()
	block := table.NewBlock()

	block.InsertVariable("a", NewInt(42))

	require.Equal(t, uint16(1), table.Size())

	locals := table.Variables()
	require.Len(t, locals, 1)
	value := locals[0]
	require.Equal(t, NewInt(42), value)
}

func TestFreeVar(t *testing.T) {
	main := NewSymbolTable()
	outerFunc := main.NewChild()
	innerFunc := outerFunc.NewChild()

	outerFunc.InsertVariable("a", NewInt(42))

	_, found := innerFunc.Lookup("whut")
	require.False(t, found)

	res, found := innerFunc.Lookup("a")
	require.True(t, found)

	exp := &Resolution{
		Symbol: &Symbol{
			Name:  "a",
			Index: 0,
			Value: NewInt(42),
		},
		Code:  ScopeFree,
		Depth: 1,
	}
	require.Equal(t, exp, res)

	freeVars := innerFunc.Free()
	require.Len(t, freeVars, 1)
	require.Equal(t, exp, freeVars[0])

	require.Len(t, outerFunc.Free(), 0)
}
