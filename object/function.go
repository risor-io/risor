package object

import (
	"bytes"
	"strings"

	"github.com/cloudcmds/tamarin/ast"
)

// Function contains the AST for user defined function and implements Object interface.
type Function struct {
	name       string
	parameters []*ast.Ident
	body       *ast.Block
	defaults   map[string]ast.Expression
	scope      Scope
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Name() string {
	if f.name == "" {
		return "anonymous"
	}
	return f.name
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	parameters := make([]string, 0)
	for _, p := range f.parameters {
		ident := p.String()
		if def, ok := f.defaults[p.String()]; ok {
			ident += "=" + def.String()
		}
		parameters = append(parameters, ident)
	}
	out.WriteString("func")
	if f.name != "" {
		out.WriteString(" " + f.name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {")
	lines := strings.Split(f.body.String(), "\n")
	if len(lines) == 1 {
		out.WriteString(" " + lines[0] + " }")
	} else if len(lines) == 0 {
		out.WriteString(" }")
	} else {
		for _, line := range lines {
			out.WriteString("\n    " + line)
		}
		out.WriteString("\n}")
	}
	return out.String()
}

func (f *Function) Body() *ast.Block {
	return f.body
}

func (f *Function) Parameters() []*ast.Ident {
	return f.parameters
}

func (f *Function) Defaults() map[string]ast.Expression {
	return f.defaults
}

func (f *Function) Scope() Scope {
	return f.scope
}

func (f *Function) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (f *Function) Interface() interface{} {
	return "function()"
}

func (f *Function) Equals(other Object) Object {
	if other.Type() == FUNCTION && f == other.(*Function) {
		return True
	}
	return False
}

func (f *Function) IsTruthy() bool {
	return true
}

func NewFunction(
	name string,
	parameters []*ast.Ident,
	body *ast.Block,
	defaults map[string]ast.Expression,
	scope Scope,
) *Function {
	return &Function{
		name:       name,
		parameters: parameters,
		body:       body,
		defaults:   defaults,
		scope:      scope,
	}
}
