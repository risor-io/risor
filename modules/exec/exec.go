package exec

import (
	"context"
	"os/exec"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func CommandFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("command", 1, 1000, args); err != nil {
		return err
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var strArgs []string
	for _, arg := range args[1:] {
		argStr, err := object.AsString(arg)
		if err != nil {
			return err
		}
		strArgs = append(strArgs, argStr)
	}
	return NewCommand(exec.Command(name, strArgs...))
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
	if err != nil {
		return object.NewError(execErr)
	}
	return object.NewString(result)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("exec", map[string]object.Object{
		"command":   object.NewBuiltin("exec.command", CommandFunc),
		"look_path": object.NewBuiltin("exec.look_path", LookPath),
	})
}
