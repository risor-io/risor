package object

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/internal/ast"
)

// Function wraps ast.Identifier array and ast.BlockStatement and implements Object interface.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Defaults   map[string]ast.Expression
	Scope      interface{} // avoids circular package dependency; is a scope.Scope
}

// Type returns the type of this object.
func (f *Function) Type() Type {
	return FUNCTION_OBJ
}

// Inspect returns a string-representation of the given object.
func (f *Function) Inspect() string {
	var out bytes.Buffer
	parameters := make([]string, 0)
	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (f *Function) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", f.Type(), method)
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (f *Function) ToInterface() interface{} {
	return "<FUNCTION>"
}
