package object

import (
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/op"
)

// Partial is a partially applied function
type Partial struct {
	fn   Object
	args []Object
}

func (f *Partial) Function() Object {
	return f.fn
}

func (f *Partial) Args() []Object {
	return f.args
}

func (f *Partial) Type() Type {
	return PARTIAL
}

func (f *Partial) Inspect() string {
	var args []string
	for _, arg := range f.args {
		args = append(args, arg.Inspect())
	}
	return fmt.Sprintf("partial(%s, %s)", f.fn.Inspect(), strings.Join(args, ", "))
}

func (f *Partial) Interface() interface{} {
	return f
}

func (f *Partial) Equals(other Object) Object {
	other, ok := other.(*Partial)
	if !ok {
		return False
	}
	if f == other {
		return True
	}
	return False
}

func (f *Partial) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (f *Partial) IsTruthy() bool {
	return true
}

func (f *Partial) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for nil: %v", opType))
}

func NewPartial(fn Object, args []Object) *Partial {
	return &Partial{
		fn:   fn,
		args: args,
	}
}
