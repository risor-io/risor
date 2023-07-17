package object

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type ScopeName string

const (
	ScopeLocal  ScopeName = "local"
	ScopeGlobal ScopeName = "global"
	ScopeFree   ScopeName = "free"
)

type Symbol struct {
	Name       string
	Index      uint16
	Value      Object
	IsConstant bool
}

type Resolution struct {
	Symbol    *Symbol
	Scope     ScopeName
	Depth     int
	FreeIndex int
}

func (r *Resolution) String() string {
	return fmt.Sprintf("resolution(symbol: %s scope: %s depth: %d)",
		r.Symbol.Name, r.Scope, r.Depth)
}

type SymbolTable struct {
	parent    *SymbolTable
	symbols   map[string]*Symbol
	variables map[string]*Symbol
	accessed  map[string]bool
	free      map[string]*Resolution
	values    []Object
	isBlock   bool
	freeCount int
}

func (t *SymbolTable) NewChild() *SymbolTable {
	return &SymbolTable{
		parent:    t,
		symbols:   map[string]*Symbol{},
		variables: map[string]*Symbol{},
		accessed:  map[string]bool{},
		free:      map[string]*Resolution{},
		isBlock:   false,
	}
}

func (t *SymbolTable) NewBlock() *SymbolTable {
	child := t.NewChild()
	child.isBlock = true
	return child
}

func (t *SymbolTable) claimIndex(value Object) (uint16, error) {
	if t.isBlock {
		return t.parent.claimIndex(value)
	}
	priorCount := len(t.values)
	if priorCount >= math.MaxUint16 {
		return 0, errors.New("too many symbols")
	}
	t.values = append(t.values, value)
	return uint16(priorCount), nil
}

func (t *SymbolTable) InsertConstant(name string, value ...Object) (*Symbol, error) {
	sym, err := t.InsertVariable(name, value...)
	if err != nil {
		return nil, err
	}
	sym.IsConstant = true
	return sym, nil
}

func (t *SymbolTable) InsertVariable(name string, value ...Object) (*Symbol, error) {
	if _, ok := t.symbols[name]; ok {
		return nil, fmt.Errorf("symbol %q already exists", name)
	}
	var obj Object
	valueCount := len(value)
	if valueCount > 1 {
		return nil, errors.New("expected at most one value")
	} else if valueCount == 1 {
		obj = value[0]
	}
	index, err := t.claimIndex(obj)
	if err != nil {
		return nil, err
	}
	s := &Symbol{Name: name, Index: index, Value: obj}
	t.symbols[name] = s
	t.variables[name] = s
	return s, nil
}

func (t *SymbolTable) InsertBuiltin(name string, value ...Object) (*Symbol, error) {
	if t.parent != nil {
		return nil, errors.New("cannot insert builtin in child table")
	}
	return t.InsertVariable(name, value...)
}

func (t *SymbolTable) SetValue(name string, value Object) error {
	s, ok := t.symbols[name]
	if !ok {
		return fmt.Errorf("symbol %q not found", name)
	}
	s.Value = value
	return nil
}

func (t *SymbolTable) IsVariable(name string) bool {
	_, ok := t.variables[name]
	return ok
}

func (t *SymbolTable) Get(name string) (*Symbol, bool) {
	s, ok := t.symbols[name]
	return s, ok
}

func (t *SymbolTable) IsGlobal() bool {
	if t.parent == nil {
		return true
	}
	if t.isBlock {
		return t.parent.IsGlobal()
	}
	return false
}

func (t *SymbolTable) Lookup(name string) (*Resolution, bool) {
	// Check if the symbol is defined directly in this table
	if s, ok := t.symbols[name]; ok {
		t.accessed[name] = true
		var scope ScopeName
		if t.IsGlobal() {
			scope = ScopeGlobal
		} else {
			scope = ScopeLocal
		}
		return &Resolution{Symbol: s, Scope: scope, Depth: 0}, true
	}
	// Check if the symbol was previously found to be a "free" variable
	if rs, ok := t.free[name]; ok {
		return rs, true
	}
	// At this point, if there is no parent then the symbol is undefined
	if t.parent == nil {
		return nil, false
	}
	// Does a parent table define the symbol?
	rs, found := t.parent.Lookup(name)
	if !found {
		return nil, false
	}
	t.accessed[name] = true
	// Check if this is a global. These are simple in that we don't
	// care about their depth and their scope always stays unchanged.
	if rs.Scope == ScopeGlobal {
		return rs, true
	}
	// Determine if this is a free variable which is defined in an outer scope.
	// Locals may stil be defined in a parent table if this is a block.
	scope := rs.Scope
	depth := rs.Depth
	if !t.isBlock {
		depth++
		scope = ScopeFree
	}
	resolution := &Resolution{Symbol: rs.Symbol, Scope: scope, Depth: depth}
	if scope == ScopeFree {
		t.free[name] = resolution
		resolution.FreeIndex = t.freeCount
		t.freeCount++
	}
	return resolution, true
}

func (t *SymbolTable) AccessedNames() []string {
	names := make([]string, 0, len(t.accessed))
	for name := range t.accessed {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *SymbolTable) InsertedNames() []string {
	names := make([]string, 0, len(t.symbols))
	for name := range t.symbols {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *SymbolTable) Size() uint16 {
	return uint16(len(t.values))
}

func (t *SymbolTable) Parent() *SymbolTable {
	return t.parent
}

func (t *SymbolTable) Root() *SymbolTable {
	current := t
	for current.parent != nil {
		current = current.parent
	}
	return current
}

func (t *SymbolTable) LocalTable() *SymbolTable {
	current := t
	for current.isBlock {
		current = current.parent
	}
	return current
}

func (t *SymbolTable) Variables() []Object {
	return t.values
}

func (t *SymbolTable) Free() []*Resolution {
	result := make([]*Resolution, len(t.free))
	for _, rs := range t.free {
		result[rs.FreeIndex] = rs
	}
	return result
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols:   map[string]*Symbol{},
		variables: map[string]*Symbol{},
		accessed:  map[string]bool{},
		free:      map[string]*Resolution{},
	}
}
