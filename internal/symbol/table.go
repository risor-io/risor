package symbol

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/cloudcmds/tamarin/object"
)

type Scope string

const (
	ScopeBuiltin Scope = "builtin"
	ScopeLocal   Scope = "local"
	ScopeGlobal  Scope = "global"
	ScopeFree    Scope = "free"
)

type Symbol struct {
	Name  string
	Index uint16
	Value object.Object
}

type Resolution struct {
	Symbol *Symbol
	Scope  Scope
	Depth  int
}

type Table struct {
	parent    *Table
	symbols   map[string]*Symbol
	variables map[string]*Symbol
	builtins  map[string]*Symbol
	accessed  map[string]bool
	free      map[string]*Resolution
	values    []object.Object
	isBlock   bool
}

func (t *Table) NewChild() *Table {
	return &Table{
		parent:    t,
		symbols:   map[string]*Symbol{},
		variables: map[string]*Symbol{},
		builtins:  map[string]*Symbol{},
		accessed:  map[string]bool{},
		free:      map[string]*Resolution{},
		isBlock:   false,
	}
}

func (t *Table) NewBlock() *Table {
	child := t.NewChild()
	child.isBlock = true
	return child
}

func (t *Table) claimIndex(value object.Object) (uint16, error) {
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

func (t *Table) InsertVariable(name string, value ...object.Object) (*Symbol, error) {
	if _, ok := t.symbols[name]; ok {
		return nil, fmt.Errorf("symbol %q already exists", name)
	}
	var obj object.Object
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

func (t *Table) InsertBuiltin(name string, value ...object.Object) (*Symbol, error) {
	if t.parent != nil {
		return nil, errors.New("cannot insert builtin in child table")
	}
	if _, ok := t.symbols[name]; ok {
		return nil, fmt.Errorf("symbol %q already exists", name)
	}
	priorCount := len(t.builtins)
	if priorCount >= math.MaxUint16 {
		return nil, errors.New("too many symbols")
	}
	s := &Symbol{Name: name, Index: uint16(priorCount)}
	valueCount := len(value)
	if valueCount > 1 {
		return nil, errors.New("expected at most one value")
	} else if valueCount == 1 {
		s.Value = value[0]
	}
	t.symbols[name] = s
	t.builtins[name] = s
	return s, nil
}

func (t *Table) IsBuiltin(name string) bool {
	_, ok := t.builtins[name]
	return ok
}

func (t *Table) IsVariable(name string) bool {
	_, ok := t.variables[name]
	return ok
}

func (t *Table) Get(name string) (*Symbol, bool) {
	s, ok := t.symbols[name]
	return s, ok
}

func (t *Table) Lookup(name string) (*Resolution, bool) {
	// Check if the symbol is defined directly in this table
	if s, ok := t.symbols[name]; ok {
		t.accessed[name] = true
		var scope Scope
		if t.IsBuiltin(name) {
			scope = ScopeBuiltin
		} else if t.parent == nil {
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
	// Check if this is a global or a builtin. These are simple in that we don't
	// care about their depth and their scope always stays unchanged.
	if rs.Scope == ScopeGlobal || rs.Scope == ScopeBuiltin {
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
	}
	return resolution, true
}

func (t *Table) AccessedNames() []string {
	names := make([]string, 0, len(t.accessed))
	for name := range t.accessed {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *Table) InsertedNames() []string {
	names := make([]string, 0, len(t.symbols))
	for name := range t.symbols {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t *Table) Size() uint16 {
	return uint16(len(t.values))
}

func (t *Table) Parent() *Table {
	return t.parent
}

func (t *Table) LocalTable() *Table {
	current := t
	for current.isBlock {
		current = current.parent
	}
	return current
}

func (t *Table) Variables() []object.Object {
	return t.values
}

func (t *Table) Builtins() []object.Object {
	result := make([]object.Object, len(t.builtins))
	for _, s := range t.builtins {
		result[s.Index] = s.Value
	}
	return result
}

func (t *Table) Free() []*Resolution {
	result := make([]*Resolution, 0, len(t.free))
	for _, rs := range t.free {
		result = append(result, rs)
	}
	return result
}

func NewTable() *Table {
	return &Table{
		symbols:   map[string]*Symbol{},
		variables: map[string]*Symbol{},
		builtins:  map[string]*Symbol{},
		accessed:  map[string]bool{},
		free:      map[string]*Resolution{},
	}
}
