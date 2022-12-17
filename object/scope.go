package object

// Scope is an interface that can be used to implement variable storage
// for a specific context.
type Scope interface {

	// Name of the scope, to aid debugging
	Name() string

	// IsReadOnly returns true iff the scope contains a variable with the
	// given name and that variable is marked as read-only.
	IsReadOnly(name string) bool

	// Get returns the object associated with the given name, and a boolean
	// indicating whether the object was found.
	Get(name string) (Object, bool)

	// Declare adds a new variable to the scope. If the variable already
	// exists, an error is returned.
	Declare(name string, obj Object, readOnly bool) error

	// Update the variable with the given name. If the variable does not
	// exist, an error is returned.
	Update(name string, obj Object) error

	// Contents returns a map of all variables in the scope.
	Contents() map[string]Object
}
