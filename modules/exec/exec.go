package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func CommandFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("command", 1, 1000, args); err != nil {
		return err
	}
	var strArgs []string
	// Two forms of arguments are supported:
	// 1. command(["ls", "-l"]) - this is the newly added form
	// 2. command("ls", "-l") - this is the original form
	if len(args) == 1 {
		if list, err := object.AsList(args[0]); err == nil {
			// This is form 1
			for _, arg := range list.Value() {
				argStr, err := object.AsString(arg)
				if err != nil {
					return err
				}
				strArgs = append(strArgs, argStr)
			}
			if len(strArgs) == 0 {
				return object.Errorf("exec.command expected at least one argument in list")
			}
			return NewCommand(exec.CommandContext(ctx, strArgs[0], strArgs[1:]...))
		}
	}
	// This is form 2
	name, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	for _, arg := range args[1:] {
		argStr, err := object.AsString(arg)
		if err != nil {
			return err
		}
		strArgs = append(strArgs, argStr)
	}
	return NewCommand(exec.CommandContext(ctx, name, strArgs...))
}

func LookPath(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("look_path", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result, execErr := exec.LookPath(path)
	if execErr != nil {
		return object.NewError(execErr)
	}
	return object.NewString(result)
}

func Exec(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("exec", 1, 3, args); err != nil {
		return err
	}
	var wasList bool
	var program string
	var optArgs []string
	if list, err := object.AsList(args[0]); err == nil {
		wasList = true
		var args []string
		for _, arg := range list.Value() {
			argStr, err := object.AsString(arg)
			if err != nil {
				return err
			}
			args = append(args, argStr)
		}
		if len(args) == 0 {
			return object.Errorf("exec expected at least one argument in list")
		}
		program = args[0]
		optArgs = args[1:]
	} else {
		program, err = object.AsString(args[0])
		if err != nil {
			return err
		}
		if len(args) > 1 {
			optArgs, err = object.AsStringSlice(args[1])
			if err != nil {
				return err
			}
		}
	}
	cmd := exec.CommandContext(ctx, program, optArgs...)

	mapOffset := 2
	if wasList {
		mapOffset = 1
	}

	if len(args) > mapOffset {
		var params *object.Map
		var errObj *object.Error
		params, errObj = object.AsMap(args[mapOffset])
		if errObj != nil {
			return errObj
		}
		if err := configureCommand(cmd, params); err != nil {
			return object.NewError(err)
		}
	}

	if cmd.Stdout == nil {
		cmd.Stdout = object.NewBuffer(nil)
	}
	if cmd.Stderr == nil {
		cmd.Stderr = object.NewBuffer(nil)
	}
	cmdObj := NewCommand(cmd)
	if err := cmdObj.Run(ctx); err != nil {
		return object.NewError(err)
	}
	return NewResult(cmd)
}

var allowedKeys = map[string]bool{
	"dir":    true,
	"stdin":  true,
	"stdout": true,
	"stderr": true,
	"env":    true,
}

func configureCommand(cmd *exec.Cmd, params *object.Map) error {
	for key := range params.Value() {
		if !allowedKeys[key] {
			return fmt.Errorf("exec found unexpected key %q", key)
		}
	}
	if stdoutObj := params.GetWithDefault("stdout", nil); stdoutObj != nil {
		stdoutBuf, ok := stdoutObj.(io.Writer)
		if !ok {
			return fmt.Errorf("exec expected io.Writer for stdout (got %s)", stdoutObj.Type())
		}
		cmd.Stdout = stdoutBuf
	}
	if stderrObj := params.GetWithDefault("stderr", nil); stderrObj != nil {
		stderrBuf, ok := stderrObj.(io.Writer)
		if !ok {
			return fmt.Errorf("exec expected io.Writer for stderr (got %s)", stderrObj.Type())
		}
		cmd.Stderr = stderrBuf
	}
	if stdinObj := params.GetWithDefault("stdin", nil); stdinObj != nil {
		switch stdinObj := stdinObj.(type) {
		case *object.ByteSlice:
			cmd.Stdin = bytes.NewBuffer(stdinObj.Value())
		case *object.String:
			cmd.Stdin = bytes.NewBufferString(stdinObj.Value())
		case io.Reader:
			cmd.Stdin = stdinObj
		default:
			return fmt.Errorf("exec expected io.Reader for stdin (got %s)", stdinObj.Type())
		}
	}
	if dirObj := params.GetWithDefault("dir", nil); dirObj != nil {
		dirStr, err := object.AsString(dirObj)
		if err != nil {
			return fmt.Errorf("exec expected string for dir (got %s)", dirObj.Type())
		}
		cmd.Dir = dirStr
	}
	if envObj := params.GetWithDefault("env", nil); envObj != nil {
		envMap, err := object.AsMap(envObj)
		if err != nil {
			return fmt.Errorf("exec expected map for env (got %s)", envObj.Type())
		}
		var env []string
		for key, value := range envMap.Value() {
			valueStr, err := object.AsString(value)
			if err != nil {
				return fmt.Errorf("exec expected string for env value (got %s)", value.Type())
			}
			env = append(env, fmt.Sprintf("%s=%s", key, valueStr))
		}
		cmd.Env = env
	}
	return nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("exec", map[string]object.Object{
		"command":   object.NewBuiltin("exec.command", CommandFunc),
		"look_path": object.NewBuiltin("exec.look_path", LookPath),
	}, Exec)
}
