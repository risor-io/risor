package scope

import (
	"fmt"

	"github.com/myzie/tamarin/internal/object"
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

	// children environments
	children []*Scope
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
		return fmt.Errorf("variable already exists: %s", name)
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
			return fmt.Errorf("cannot update %s since it is read-only", name)
		}
		s.store[name] = obj
		return nil
	}
	if s.parent != nil {
		return s.parent.Update(name, obj)
	}
	return fmt.Errorf("unknown variable: %s", name)
}

func (s *Scope) Contents() map[string]object.Object {
	contents := make(map[string]object.Object, len(s.store))
	for k, v := range s.store {
		contents[k] = v
	}
	return contents
}

func (s *Scope) NewChild(opts Opts) *Scope {
	opts.Parent = s
	child := New(opts)
	s.children = append(s.children, child)
	return child
}

func (s *Scope) Children() []*Scope {
	return s.children
}

func (s *Scope) AddBuiltin(name string, fn object.BuiltinFunction) error {
	if err := s.Declare(name, &object.Builtin{Fn: fn}, true); err != nil {
		return fmt.Errorf("failed to define %s: %w", name, err)
	}
	return nil
}

func (s *Scope) AddBuiltins(funcs []Builtin) error {
	for _, f := range funcs {
		if err := s.AddBuiltin(f.Name, f.Func); err != nil {
			return err
		}
	}
	return nil
}

type Builtin struct {
	Name string
	Func object.BuiltinFunction
}
