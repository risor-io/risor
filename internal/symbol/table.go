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
	Name  string
	Index int
	Scope Scope
	Attrs Attrs
}

type Table struct {
	parent  *Table
	symbols map[string]*Symbol
	free    []*Symbol
}

func (t *Table) NewChild() *Table {
	return &Table{
		parent:  t,
		symbols: map[string]*Symbol{},
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
	if t.parent == nil {
		s.Scope = ScopeGlobal
	} else if attrs.IsBuiltin {
		s.Scope = ScopeBuiltin
	} else {
		s.Scope = ScopeLocal
	}
	t.symbols[name] = s
	fmt.Println("Insert symbol:", name, s.Index, s)
	return s, nil
}

func (t *Table) Lookup(name string) (*Symbol, bool) {
	if s, ok := t.symbols[name]; ok {
		return s, true
	}
	if t.parent != nil {
		return t.parent.Lookup(name)
	}
	return nil, false
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

func (t *Table) Free() []*Symbol {
	return t.free
}

func NewTable() *Table {
	return &Table{
		symbols: map[string]*Symbol{},
	}
}
