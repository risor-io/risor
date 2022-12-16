package object

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/ast"
)

// Function wraps ast.Identifier array and ast.BlockStatement and implements Object interface.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Defaults   map[string]ast.Expression
	Scope      interface{} // avoids circular package dependency; is a scope.Scope
}

func (f *Function) Type() Type {
	return FUNCTION
}

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

func (f *Function) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (f *Function) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", f.Type(), method)
}

func (f *Function) ToInterface() interface{} {
	return "Function()"
}

func (f *Function) Equals(other Object) Object {
	if other.Type() == FUNCTION && f == other.(*Function) {
		return True
	}
	return False
}
