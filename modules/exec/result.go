package exec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Result struct {
	cmd *exec.Cmd
}

func (r *Result) Type() object.Type {
	return "exec.result"
}

func (r *Result) Inspect() string {
	return fmt.Sprintf("exec.result(status: %d)", 200)
}

func (r *Result) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "pid":
		return object.NewInt(int64(r.cmd.Process.Pid)), true
	case "stdout":
		buffer, ok := r.cmd.Stdout.(*object.Buffer)
		if !ok {
			return object.NewError(fmt.Errorf("eval error: exec.result.stdout does not support stdout type %T", r.cmd.Stdout)), true
		}
		return buffer, true
	case "stderr":
		buffer, ok := r.cmd.Stderr.(*object.Buffer)
		if !ok {
			return object.NewError(fmt.Errorf("eval error: exec.result.stderr does not support stderr type %T", r.cmd.Stderr)), true
		}
		return buffer, true
	case "json":
		return object.NewBuiltin("exec.result.json",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 0 {
					return object.NewArgsError("json", 0, len(args))
				}
				return r.JSON()
			},
		), true
	}
	return nil, false
}

func (r *Result) JSON() object.Object {
	var data []byte
	switch stdout := r.cmd.Stdout.(type) {
	case *object.Buffer:
		data = stdout.Value().Bytes()
	default:
		return object.NewError(fmt.Errorf("eval error: exec.result.json does not support stdout type %T", stdout))
	}
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return object.Errorf("value error: json.unmarshal failed with: %s", err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.Errorf("type error: json.unmarshal failed")
	}
	return scriptObj
}

func (r *Result) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("eval error: exec.result does not support attribute assignment")
}

func (r *Result) Interface() interface{} {
	return ""
}

func (r *Result) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare exec.result")
}

func (r *Result) Equals(other object.Object) object.Object {
	if r == other {
		return object.True
	}
	return object.False
}

func (r *Result) IsTruthy() bool {
	return true
}

func (r *Result) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for exec.result: %v", opType))
}

func (r *Result) Cost() int {
	return 0
}

func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status string `json:"status"`
	}{
		Status: "200",
	})
}

func NewResult(cmd *exec.Cmd) *Result {
	return &Result{cmd: cmd}
}
