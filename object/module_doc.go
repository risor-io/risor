package object

import (
	"fmt"
	"strings"

	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/rdoc"
)

var _ Object = (*ModuleDoc)(nil)

// ModuleDoc contains documentation for a Risor function.
type ModuleDoc struct {
	*base

	docs *rdoc.Module
}

func (d *ModuleDoc) Type() Type {
	return MODULE_DOC
}

func (d *ModuleDoc) Inspect() string {
	if d.docs == nil {
		return "module_doc()"
	}
	return fmt.Sprintf("module_doc(%s)", d.docs.Name)
}

func (d *ModuleDoc) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Module: %s", d.docs.Name))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("Description: %s", d.docs.Description))
	s.WriteString("\n\n")
	return s.String()
}

func (d *ModuleDoc) Interface() interface{} {
	return nil
}

func (d *ModuleDoc) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(d.docs.Name), true
	case "desc":
		return NewString(d.docs.Description), true
	case "functions":
		var fns []Object
		for _, f := range d.docs.Functions {
			fns = append(fns, NewFunctionDoc(f))
		}
		return NewList(fns), true
	}
	return nil, false
}

func (d *ModuleDoc) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for module_doc: %v", opType))
}

func (d *ModuleDoc) Equals(other Object) Object {
	if d == other {
		return True
	}
	return False
}

func (d *ModuleDoc) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal module_doc")
}

func NewModuleDoc(docs *rdoc.Module) *ModuleDoc {
	return &ModuleDoc{docs: docs}
}
