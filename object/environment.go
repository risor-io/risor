package object

import (
	"fmt"
	"os"
	"strings"
)

// Environment stores our functions, variables, constants, etc.
type Environment struct {
	// store holds variables, including functions.
	store map[string]Object

	// readonly marks names as read-only.
	readonly map[string]bool

	// outer holds any parent environment.  Our env. allows
	// nesting to implement scope.
	outer *Environment

	// permit stores the names of variables we can set in this
	// environment, if any
	permit []string
}

// NewEnvironment creates new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	r := make(map[string]bool)
	return &Environment{store: s, readonly: r, outer: nil}
}

// NewEnclosedEnvironment create new environment by outer parameter
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewTemporaryScope creates a temporary scope where some values
// are ignored.
//
// This is used as a sneaky hack to allow `foreach` to access all
// global values as if they were local, but prevent the index/value
// keys from persisting.
func NewTemporaryScope(outer *Environment, keys []string) *Environment {
	env := NewEnvironment()
	env.outer = outer
	env.permit = keys
	return env
}

// Names returns the names of every known-value with the
// given prefix.
//
// This function is used by `invokeMethod` to get the methods
// associated with a particular class-type.
func (e *Environment) Names(prefix string) []string {
	var ret []string

	for key := range e.store {
		if strings.HasPrefix(key, prefix) {
			ret = append(ret, key)
		}

		// Functions with an "object." prefix are available
		// to all object-methods.
		if strings.HasPrefix(key, "object.") {
			ret = append(ret, key)
		}
	}
	return ret
}

// Get returns the value of a given variable, by name.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores the value of a variable, by name.
func (e *Environment) Set(name string, val Object) Object {

	//
	// If a variable is constant then we don't allow it to be changed.
	//
	// But constants are scoped, they are not global, so we only need
	// to look in the current scope - not any parent.
	//
	// i.e. The parent-scope might have a constant-value, but
	// we just don't care.  Consider the following code:
	//
	//    const a = 3.13;
	//    function foo() {
	//       let a = 1976;
	//    };
	//
	// The variable inside the function _should_ not be constant.
	//
	cur := e.store[name]
	if cur != nil && e.readonly[name] {
		fmt.Printf("Attempting to modify '%s' denied; it was defined as a constant.\n", name)
		os.Exit(3)
	}

	//
	// Store the (updated) value.
	//
	if len(e.permit) > 0 {
		for _, v := range e.permit {
			// we're permitted to store this variable
			if v == name {
				e.store[name] = val
				return val
			}
		}
		// ok we're not permitted, we must store in the parent
		if e.outer != nil {
			return e.outer.Set(name, val)
		}
		fmt.Printf("scoping weirdness; please report a bug\n")
		os.Exit(5)
	}
	e.store[name] = val
	return val
}

// SetConst sets the value of a constant by name.
func (e *Environment) SetConst(name string, val Object) Object {

	// store the value
	e.store[name] = val

	// flag as read-only.
	e.readonly[name] = true

	return val
}
