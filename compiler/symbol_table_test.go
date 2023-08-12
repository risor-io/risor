package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	table := NewSymbolTable()

	require.Nil(t, table.Parent())
	require.Equal(t, uint16(0), table.Count())

	a, err := table.InsertVariable("a")
	require.Nil(t, err)
	require.Equal(t, uint16(0), a.Index())
	require.Equal(t, "a", a.Name())
	require.Nil(t, a.Value())

	b, err := table.InsertVariable("b")
	require.Nil(t, err)
	require.Equal(t, uint16(1), b.Index())
	require.Equal(t, "b", b.Name())
	require.Nil(t, b.Value())

	c, err := table.InsertVariable("c")
	require.Nil(t, err)
	require.Equal(t, uint16(2), c.Index())
	require.Equal(t, "c", c.Name())
	require.Nil(t, c.Value())

	// The size is the count of variables
	require.Equal(t, uint16(3), table.Count())

	require.True(t, table.IsDefined("a"))
	require.True(t, table.IsDefined("b"))
	require.True(t, table.IsDefined("c"))
}

func TestBlock(t *testing.T) {
	table := NewSymbolTable()
	block := table.NewBlock()

	block.InsertVariable("a", 42)

	require.Equal(t, uint16(1), table.Count())
	require.Equal(t, 42, table.Symbol(0).Value())
}

func TestFreeVar(t *testing.T) {
	main := NewSymbolTable()
	outerFunc := main.NewChild()
	innerFunc := outerFunc.NewChild()

	outerFunc.InsertVariable("a", 42)

	_, found := innerFunc.Resolve("whut")
	require.False(t, found)

	res, found := innerFunc.Resolve("a")
	require.True(t, found)

	exp := &Resolution{
		symbol: &Symbol{
			name:  "a",
			index: 0,
			value: 42,
		},
		scope: Free,
		depth: 1,
	}
	require.Equal(t, exp, res)

	require.Equal(t, uint16(1), innerFunc.FreeCount())
	require.Equal(t, exp, innerFunc.Free(0))
	require.Equal(t, uint16(0), outerFunc.FreeCount())
}

func TestConstant(t *testing.T) {
	main := NewSymbolTable()
	outerFunc := main.NewChild()
	innerFunc := outerFunc.NewChild()

	outerFunc.InsertConstant("a", 42)
	outerFunc.InsertVariable("b", 42)

	resolution, found := innerFunc.Resolve("a")
	require.True(t, found)
	require.True(t, resolution.symbol.isConstant)

	resolution, found = innerFunc.Resolve("b")
	require.True(t, found)
	require.False(t, resolution.symbol.isConstant)
}
