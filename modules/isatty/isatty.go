package isatty

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/risor-io/risor/object"
	ros "github.com/risor-io/risor/os"
)

func IsTerminal(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("isatty.is_terminal", 1, numArgs)
	}
	switch arg := args[0].(type) {
	case *object.File:
		osFile := arg.Value()
		if osFileObj, ok := osFile.(*os.File); ok {
			return object.NewBool(isatty.IsTerminal(osFileObj.Fd()))
		}
		return object.Errorf("argument error: unsupported file type provided")
	case *object.Int:
		fd, err := object.AsInt(args[0])
		if err != nil {
			return err
		}
		return object.NewBool(isatty.IsTerminal(uintptr(fd)))
	default:
		return object.Errorf("argument error: expected file or int, got %s", args[0].Type())
	}
}

func IsCygwinTerminal(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 1 {
		return object.NewArgsError("isatty.is_cygwin_terminal", 1, numArgs)
	}
	switch arg := args[0].(type) {
	case *object.File:
		osFile := arg.Value()
		if osFileObj, ok := osFile.(*os.File); ok {
			return object.NewBool(isatty.IsCygwinTerminal(osFileObj.Fd()))
		}
		return object.Errorf("argument error: unsupported file type provided")
	case *object.Int:
		fd, err := object.AsInt(args[0])
		if err != nil {
			return err
		}
		return object.NewBool(isatty.IsCygwinTerminal(uintptr(fd)))
	default:
		return object.Errorf("argument error: expected file or int, got %s", args[0].Type())
	}
}

func IsTTY(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewArgsError("isatty.is_terminal", 0, len(args))
	}
	stdout := ros.GetDefaultOS(ctx).Stdout()
	if stdoutFile, ok := stdout.(*os.File); ok {
		fd := stdoutFile.Fd()
		if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
			return object.True
		}
		return object.False
	}
	return object.Errorf("argument error: unsupported file type provided")
}

func Module() *object.Module {
	return object.NewBuiltinsModule("isatty", map[string]object.Object{
		"is_terminal":        object.NewBuiltin("is_terminal", IsTerminal),
		"is_cygwin_terminal": object.NewBuiltin("is_cygwin_terminal", IsCygwinTerminal),
	}, IsTTY)
}
