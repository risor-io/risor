package object

import (
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/op"
)

type Example struct {
	*base
	name   string
	desc   string
	code   string
	result string
	err    string
}

func (ex *Example) Type() Type {
	return EXAMPLE
}

func (ex *Example) Inspect() string {
	return fmt.Sprintf("example(name: %s)", ex.name)
}

func (ex *Example) String() string {
	return ex.Inspect()
}

func (ex *Example) Interface() interface{} {
	return ex
}

func (ex *Example) Equals(other Object) Object {
	if ex == other {
		return True
	}
	return False
}

func (ex *Example) GetAttr(name string) (Object, bool) {
	switch name {
	case "name":
		return NewString(ex.name), true
	case "desc":
		return NewString(ex.desc), true
	case "code":
		return NewString(ex.code), true
	case "result":
		return NewString(ex.result), true
	case "err":
		return NewString(ex.err), true
	}
	return nil, false
}

func (ex *Example) IsTruthy() bool {
	return true
}

func (ex *Example) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for example: %v", opType))
}

func (ex *Example) MarshalJSON() ([]byte, error) {
	type example struct {
		Code   string `json:"code"`
		Name   string `json:"name,omitempty"`
		Desc   string `json:"desc,omitempty"`
		Result string `json:"result,omitempty"`
		Err    string `json:"err,omitempty"`
	}
	return json.Marshal(example{
		Name:   ex.name,
		Code:   ex.code,
		Desc:   ex.desc,
		Result: ex.result,
		Err:    ex.err,
	})
}

type ExampleSpec struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Result string `json:"result,omitempty"`
	Err    string `json:"err,omitempty"`
}

func Examples(specs []ExampleSpec) []*Example {
	var examples []*Example
	for _, spec := range specs {
		examples = append(examples, &Example{
			name:   spec.Name,
			code:   spec.Code,
			result: spec.Result,
			err:    spec.Err,
		})
	}
	return examples
}
