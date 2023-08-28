package exec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Command struct {
	value  *exec.Cmd
	stdin  object.Object
	stdout object.Object
	stderr object.Object
}

func (c *Command) Inspect() string {
	var args []string
	for _, arg := range c.value.Args {
		args = append(args, fmt.Sprintf("%q", arg))
	}
	return fmt.Sprintf("exec.command(%s)", strings.Join(args, ", "))
}

func (c *Command) Type() object.Type {
	return "exec.command"
}

func (c *Command) Value() *exec.Cmd {
	return c.value
}

func (c *Command) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "path":
		return object.NewString(c.value.Path), true
	case "dir":
		return object.NewString(c.value.Dir), true
	case "env":
		var env []object.Object
		for _, e := range c.value.Env {
			env = append(env, object.NewString(e))
		}
		return object.NewList(env), true
	case "stdin":
		if c.stdin == nil {
			return object.Nil, true
		}
		return c.stdin, true
	case "stdout":
		if c.stdout == nil {
			return object.Nil, true
		}
		return c.stdout, true
	case "stderr":
		if c.stderr == nil {
			return object.Nil, true
		}
		return c.stderr, true
	case "run":
		return object.NewBuiltin("exec.command.run", func(ctx context.Context, args ...object.Object) object.Object {
			if err := c.value.Run(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "combined_output":
		return object.NewBuiltin("exec.command.combined_output", func(ctx context.Context, args ...object.Object) object.Object {
			output, err := c.value.CombinedOutput()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewByteSlice(output)
		}), true
	case "environ":
		return object.NewBuiltin("exec.command.environ", func(ctx context.Context, args ...object.Object) object.Object {
			env := c.value.Environ()
			var envStr []object.Object
			for _, e := range env {
				envStr = append(envStr, object.NewString(e))
			}
			return object.NewList(envStr)
		}), true
	case "output":
		return object.NewBuiltin("exec.command.output", func(ctx context.Context, args ...object.Object) object.Object {
			output, err := c.value.Output()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewByteSlice(output)
		}), true
	case "start":
		return object.NewBuiltin("exec.command.start", func(ctx context.Context, args ...object.Object) object.Object {
			if err := c.value.Start(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "wait":
		return object.NewBuiltin("exec.command.wait", func(ctx context.Context, args ...object.Object) object.Object {
			if err := c.value.Wait(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	}
	return nil, false
}

func (c *Command) SetAttr(name string, value object.Object) error {
	switch name {
	case "path":
		path, err := object.AsString(value)
		if err != nil {
			return err.Value()
		}
		c.value.Path = path
	case "dir":
		dir, err := object.AsString(value)
		if err != nil {
			return err.Value()
		}
		c.value.Dir = dir
	case "env":
		env, err := object.AsList(value)
		if err != nil {
			return err.Value()
		}
		var envStr []string
		for _, e := range env.Value() {
			item, err := object.AsString(e)
			if err != nil {
				return err.Value()
			}
			envStr = append(envStr, item)
		}
		c.value.Env = envStr
	case "stdin":
		stdin, err := object.AsReader(value)
		if err != nil {
			return err.Value()
		}
		c.value.Stdin = stdin
		c.stdin = value
	case "stdout":
		stdout, err := object.AsWriter(value)
		if err != nil {
			return err.Value()
		}
		c.value.Stdout = stdout
		c.stdout = value
	case "stderr":
		stderr, err := object.AsWriter(value)
		if err != nil {
			return err.Value()
		}
		c.value.Stderr = stderr
		c.stderr = value
	default:
		return fmt.Errorf("attribute error: exec.command object has no attribute %q", name)
	}
	return nil
}

func (c *Command) Interface() interface{} {
	return c.value
}

func (c *Command) String() string {
	return c.Inspect()
}

func (c *Command) Compare(other object.Object) (int, error) {
	return 0, errors.New("type error: unable to compare exec.command")
}

func (c *Command) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Command) IsTruthy() bool {
	return true
}

func (c *Command) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("eval error: unsupported operation for exec.command: %v ", opType))
}

func (c *Command) Cost() int {
	return 8
}

func (c *Command) MarshalJSON() ([]byte, error) {
	return nil, errors.New("type error: unable to marshal exec.command")
}

func NewCommand(cmd *exec.Cmd) *Command {
	return &Command{value: cmd}
}
