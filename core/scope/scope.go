package scope

import (
	"fmt"
	"sort"

	"github.com/cloudcmds/tamarin/core/object"
)

// Scope stores our functions, variables, constants, etc.
type Scope struct {
	// name of the scope
	name string

	// holds variables, including functions
	store map[string]object.Object

	// marks named variables as read-only
	readOnly map[string]bool

	// optional parent environment
	parent *Scope
}

type Opts struct {
	Name   string
	Parent *Scope
}

// New creates and returns a new, empty scope
func New(opts Opts) *Scope {
	return &Scope{
		name:     opts.Name,
		parent:   opts.Parent,
		store:    map[string]object.Object{},
		readOnly: map[string]bool{},
	}
}

func (s *Scope) Name() string {
	return s.name
}

func (s *Scope) IsReadOnly(name string) bool {
	return s.readOnly[name]
}

func (s *Scope) Get(name string) (object.Object, bool) {
	if obj, ok := s.store[name]; ok {
		return obj, true
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil, false
}

func (s *Scope) Declare(name string, obj object.Object, readOnly bool) error {
	if _, exists := s.store[name]; exists {
		return fmt.Errorf("assignment error: %q is already set", name)
	}
	s.store[name] = obj
	if readOnly {
		s.readOnly[name] = true
	}
	return nil
}

func (s *Scope) Update(name string, obj object.Object) error {
	if _, ok := s.store[name]; ok {
		if s.IsReadOnly(name) {
			return fmt.Errorf("assignment error: %q is read-only", name)
		}
		s.store[name] = obj
		return nil
	}
	if s.parent != nil {
		return s.parent.Update(name, obj)
	}
	return fmt.Errorf("name error: %q is not defined", name)
}

func (s *Scope) Contents() map[string]object.Object {
	contents := make(map[string]object.Object, len(s.store))
	for k, v := range s.store {
		contents[k] = v
	}
	return contents
}

func (s *Scope) Clear() {
	for k := range s.store {
		delete(s.store, k)
		delete(s.readOnly, k)
	}
}

func (s *Scope) NewChild(opts Opts) *Scope {
	opts.Parent = s
	child := New(opts)
	return child
}

func (s *Scope) Keys() []string {
	var keys []string
	for k := range s.store {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *Scope) AddBuiltin(b *object.Builtin) error {
	if err := s.Declare(b.Name(), b, true); err != nil {
		return fmt.Errorf("failed to define %s: %w", b.Key(), err)
	}
	return nil
}

func (s *Scope) AddBuiltins(funcs []*object.Builtin) error {
	for _, f := range funcs {
		if err := s.AddBuiltin(f); err != nil {
			return err
		}
	}
	return nil
}
