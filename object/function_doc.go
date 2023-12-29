package object

import (
	"fmt"
	"strings"

	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/rdoc"
)

var _ Object = (*FunctionDoc)(nil)

// FunctionDoc contains documentation for a Risor function.
type FunctionDoc struct {
	*base

	docs *rdoc.Function
}

func (d *FunctionDoc) Type() Type {
	return FUNCTION_DOC
}

func (d *FunctionDoc) Inspect() string {
	if d.docs == nil {
		return "function_doc()"
	}
	return fmt.Sprintf("function_doc(%s)", d.docs.Name)
}

func (d *FunctionDoc) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Documentation for function %s.%s", d.docs.Module, d.docs.Name))
	s.WriteString("\n\n")
	s.WriteString(d.docs.Name)
	s.WriteString(")")
	return s.String()
}

func (d *FunctionDoc) Interface() interface{} {
	return nil
}

func (d *FunctionDoc) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(d.docs.Name), true
	case "sig":
		return NewString(d.docs.Signature), true
	case "desc":
		return NewString(d.docs.Description), true
	case "examples":
		var examples []Object
		for _, ex := range d.docs.Examples {
			for _, stmt := range ex.Statements {
				examples = append(examples, &Example{
					name:   ex.Name,
					code:   stmt.Code,
					result: stmt.Result,
				})
			}
		}
		return NewList(examples), true
	}
	return nil, false
}

func (d *FunctionDoc) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for function_doc: %v", opType))
}

func (d *FunctionDoc) Equals(other Object) Object {
	if d == other {
		return True
	}
	return False
}

func (d *FunctionDoc) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal function_doc")
}

func NewFunctionDoc(docs *rdoc.Function) *FunctionDoc {
	return &FunctionDoc{docs: docs}
}
