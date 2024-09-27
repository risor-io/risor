package exec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Command struct {
	value *exec.Cmd
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
	case "stdout":
		return c.Stdout(), true
	case "stderr":
		return c.Stderr(), true
	case "run":
		return object.NewBuiltin("exec.command.run", func(ctx context.Context, args ...object.Object) object.Object {
			if err := c.Run(ctx); err != nil {
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
	case "stdout":
		stdout, err := object.AsWriter(value)
		if err != nil {
			return err.Value()
		}
		c.value.Stdout = stdout
	case "stderr":
		stderr, err := object.AsWriter(value)
		if err != nil {
			return err.Value()
		}
		c.value.Stderr = stderr
	default:
		return object.TypeErrorf("type error: exec.command object has no attribute %q", name)
	}
	return nil
}

func (c *Command) Stdout() object.Object {
	if c.value.Stdout == nil {
		return object.Nil
	}
	switch value := c.value.Stdout.(type) {
	case *object.Buffer:
		return object.NewString(value.Value().String())
	default:
		return object.Nil
	}
}

func (c *Command) Stderr() object.Object {
	if c.value.Stderr == nil {
		return object.Nil
	}
	switch value := c.value.Stderr.(type) {
	case *object.Buffer:
		return object.NewString(value.Value().String())
	default:
		return object.Nil
	}
}

func (c *Command) Run(ctx context.Context) error {
	if c.value.Stdout == nil {
		c.value.Stdout = object.NewBuffer(nil)
	}
	if c.value.Stderr == nil {
		c.value.Stderr = object.NewBuffer(nil)
	}
	return c.value.Run()
}

func (c *Command) Interface() interface{} {
	return c.value
}

func (c *Command) String() string {
	return c.Inspect()
}

func (c *Command) Compare(other object.Object) (int, error) {
	return 0, errz.TypeErrorf("type error: unable to compare exec.command")
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
	return object.TypeErrorf("type error: unsupported operation for exec.command: %v ", opType)
}

func (c *Command) Cost() int {
	return 8
}

func (c *Command) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal exec.command")
}

func NewCommand(cmd *exec.Cmd) *Command {
	return &Command{value: cmd}
}
