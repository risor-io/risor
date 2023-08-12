package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/op"
)

// Function is a function that has been compiled to bytecode.
type Function struct {
	*base
	name          string
	parameters    []string
	defaults      []Object
	defaultsCount int
	code          *compiler.Code
	fn            *compiler.Function
	instructions  []op.Code
	freeVars      []*Cell
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	parameters := make([]string, 0)
	for i, name := range f.parameters {
		if def := f.defaults[i]; def != nil {
			name += "=" + def.Inspect()
		}
		parameters = append(parameters, name)
	}
	out.WriteString("func")
	if f.name != "" {
		out.WriteString(" " + f.name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {")
	lines := strings.Split(f.code.Source(), "\n")
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

func (f *Function) String() string {
	if f.name != "" {
		return fmt.Sprintf("func %s() { ... }", f.name)
	}
	return "func() { ... }"
}

func (f *Function) Interface() interface{} {
	return nil
}

func (f *Function) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for function: %v", opType))
}

func (f *Function) Equals(other Object) Object {
	if f == other {
		return True
	}
	return False
}

func (f *Function) Instructions() []op.Code {
	if f.instructions == nil {
		count := f.code.InstructionCount()
		f.instructions = make([]op.Code, count)
		for i := 0; i < count; i++ {
			f.instructions[i] = f.code.Instruction(i)
		}
	}
	return f.instructions
}

func (f *Function) FreeVars() []*Cell {
	return f.freeVars
}

func (f *Function) Code() *compiler.Code {
	return f.code
}

func (f *Function) Function() *compiler.Function {
	return f.fn
}

func (f *Function) Parameters() []string {
	return f.parameters
}

func (f *Function) Defaults() []Object {
	return f.defaults
}

func (f *Function) RequiredArgsCount() int {
	return len(f.parameters) - f.defaultsCount
}

func (f *Function) LocalsCount() int {
	return int(f.code.LocalsCount())
}

func (f *Function) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal function")
}

type FunctionOpts struct {
	Name           string
	ParameterNames []string
	Defaults       []Object
	Code           *compiler.Code
}

func NewFunction(fn *compiler.Function) *Function {

	// Parameter defaults
	var defaults []Object
	var defaultsCount int
	for i := 0; i < fn.DefaultsCount(); i++ {
		value := fn.Default(i)
		if value != nil {
			defaultsCount++
			defaults = append(defaults, FromGoType(value))
		} else {
			defaults = append(defaults, nil)
		}
	}

	// Parameter names
	var parameters []string
	for i := 0; i < fn.ParametersCount(); i++ {
		parameters = append(parameters, fn.Parameter(i))
	}

	return &Function{
		name:          fn.Name(),
		code:          fn.Code(),
		parameters:    parameters,
		defaults:      defaults,
		defaultsCount: defaultsCount,
	}
}

func NewClosure(
	fn *Function,
	freeVars []*Cell,
) *Function {
	return &Function{
		name:          fn.name,
		parameters:    fn.parameters,
		defaults:      fn.defaults,
		defaultsCount: fn.defaultsCount,
		code:          fn.Code(),
		freeVars:      freeVars,
	}
}
