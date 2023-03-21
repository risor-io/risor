package symbol

import (
	"fmt"
)

type Scope string

const (
	ScopeBuiltin Scope = "builtin"
	ScopeFree    Scope = "free"
	ScopeGlobal  Scope = "global"
	ScopeLocal   Scope = "local"
)

type Attrs struct {
	IsConstant bool
	IsBuiltin  bool
	Value      any
}

type Symbol struct {
	Name   string
	Index  int
	Attrs  Attrs
	IsFree bool
}

type ResolvedSymbol struct {
	Symbol *Symbol
	Scope  Scope
	Depth  int
}

type Table struct {
	parent   *Table
	symbols  map[string]*Symbol
	accessed map[string]bool
	free     []*ResolvedSymbol
}

func (t *Table) NewChild() *Table {
	return &Table{
		parent:   t,
		symbols:  map[string]*Symbol{},
		accessed: map[string]bool{},
	}
}

func (t *Table) Insert(name string, attrs Attrs) (*Symbol, error) {
	if _, ok := t.symbols[name]; ok {
		return nil, fmt.Errorf("symbol %q already exists", name)
	}
	s := &Symbol{
		Name:  name,
		Index: len(t.symbols),
		Attrs: attrs,
	}
	t.symbols[name] = s
	// fmt.Println("Insert symbol:", name, s.Index, t.DefaultScope())
	return s, nil
}

func (t *Table) DefaultScope() Scope {
	if t.parent == nil {
		return ScopeGlobal
	}
	return ScopeLocal
}

func (t *Table) Lookup(name string) (*ResolvedSymbol, bool) {
	fmt.Println("Lookup", name, t.symbols)
	if s, ok := t.symbols[name]; ok {
		t.accessed[name] = true
		return &ResolvedSymbol{
			Symbol: s,
			Scope:  t.DefaultScope(),
			Depth:  0,
		}, true
	}
	if t.parent == nil {
		return nil, false
	}
	rs, found := t.parent.Lookup(name)
	if !found {
		return nil, false
	}
	if rs.Scope == ScopeGlobal {
		t.accessed[name] = true
		return rs, true
	}
	resolution := &ResolvedSymbol{
		Symbol: rs.Symbol,
		Scope:  ScopeFree,
		Depth:  rs.Depth + 1,
	}
	t.free = append(t.free, resolution)
	t.accessed[name] = true
	fmt.Printf("FREE SYMBOL: %s %+v\n", name, resolution)
	return resolution, true
}

func (t *Table) ShallowLookup(name string) (*Symbol, bool) {
	s, ok := t.symbols[name]
	return s, ok
}

func (t *Table) Size() int {
	return len(t.symbols)
}

func (t *Table) Names() []string {
	names := make([]string, len(t.symbols))
	for name := range t.symbols {
		names = append(names, name)
	}
	return names
}

func (t *Table) Map() map[string]*Symbol {
	return t.symbols
}

func (t *Table) Parent() *Table {
	return t.parent
}

func (t *Table) Free() []*ResolvedSymbol {
	return t.free
}

func NewTable() *Table {
	return &Table{
		symbols:  map[string]*Symbol{},
		accessed: map[string]bool{},
	}
}
