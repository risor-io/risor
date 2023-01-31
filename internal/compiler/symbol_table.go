package compiler

import "fmt"

type Scope string

const (
	ScopeBuiltin Scope = "builtin"
	ScopeFree    Scope = "free"
	ScopeGlobal  Scope = "global"
	ScopeLocal   Scope = "local"
)

type SymbolAttrs struct {
	IsConstant bool
	IsBuiltin  bool
	// Type       string
}

type Symbol struct {
	Name  string
	Index int
	Scope Scope
	Attrs SymbolAttrs
}

type SymbolTable struct {
	parent  *SymbolTable
	symbols map[string]*Symbol
	free    []*Symbol
}

func (t *SymbolTable) NewChild() *SymbolTable {
	return &SymbolTable{
		parent:  t,
		symbols: map[string]*Symbol{},
	}
}

func (t *SymbolTable) Insert(name string, attrs SymbolAttrs) (*Symbol, error) {
	if _, ok := t.symbols[name]; ok {
		return nil, fmt.Errorf("symbol %q already exists", name)
	}
	s := &Symbol{
		Name:  name,
		Index: len(t.symbols),
		Attrs: attrs,
	}
	if attrs.IsBuiltin {
		s.Scope = ScopeBuiltin
	} else if t.parent == nil {
		s.Scope = ScopeGlobal
	} else {
		s.Scope = ScopeLocal
	}
	t.symbols[name] = s
	return s, nil
}

func (t *SymbolTable) Lookup(name string) (*Symbol, bool) {
	if s, ok := t.symbols[name]; ok {
		return s, true
	}
	if t.parent != nil {
		return t.parent.Lookup(name)
	}
	return nil, false
}

func (t *SymbolTable) ShallowLookup(name string) (*Symbol, bool) {
	s, ok := t.symbols[name]
	return s, ok
}

func (t *SymbolTable) Size() int {
	return len(t.symbols)
}

func (t *SymbolTable) Names() []string {
	names := make([]string, len(t.symbols))
	for name := range t.symbols {
		names = append(names, name)
	}
	return names
}

func (t *SymbolTable) Parent() *SymbolTable {
	return t.parent
}

func (t *SymbolTable) Free() []*Symbol {
	return t.free
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: map[string]*Symbol{},
	}
}
