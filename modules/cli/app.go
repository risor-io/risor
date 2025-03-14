package cli

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/os"
	ucli "github.com/urfave/cli/v2"
)

const APP object.Type = "cli.app"

type App struct {
	value *ucli.App
}

func (app *App) Type() object.Type {
	return APP
}

func (app *App) Inspect() string {
	return fmt.Sprintf("%s()", app.Type())
}

func (app *App) Interface() interface{} {
	return app.value
}

func (app *App) IsTruthy() bool {
	return true
}

func (app *App) Cost() int {
	return 0
}

func (app *App) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal %s", APP)
}

func (app *App) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", APP, opType)
}

func (app *App) Equals(other object.Object) object.Object {
	return object.NewBool(app == other)
}

func (app *App) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "name":
		return object.NewString(app.value.Name), true
	case "help_name":
		return object.NewString(app.value.HelpName), true
	case "usage":
		return object.NewString(app.value.Usage), true
	case "usage_text":
		return object.NewString(app.value.UsageText), true
	case "args":
		return object.NewBool(app.value.Args), true
	case "args_usage":
		return object.NewString(app.value.ArgsUsage), true
	case "version":
		return object.NewString(app.value.Version), true
	case "description":
		return object.NewString(app.value.Description), true
	case "default_command":
		return object.NewString(app.value.DefaultCommand), true
	case "commands":
		var commands []object.Object
		for _, cmd := range app.value.Commands {
			commands = append(commands, NewCommand(cmd))
		}
		return object.NewList(commands), true
	case "flags":
		var flags []object.Object
		for _, flag := range app.value.Flags {
			flags = append(flags, NewFlag(flag))
		}
		return object.NewList(flags), true
	case "enable_bash_completion":
		return object.NewBool(app.value.EnableBashCompletion), true
	case "hide_help":
		return object.NewBool(app.value.HideHelp), true
	case "hide_help_command":
		return object.NewBool(app.value.HideHelpCommand), true
	case "hide_version":
		return object.NewBool(app.value.HideVersion), true
	case "run":
		return object.NewBuiltin("cli.app.run",
			func(ctx context.Context, args ...object.Object) object.Object {
				if err := arg.RequireRange("cli.app.run", 0, 1, args); err != nil {
					return err
				}
				var strArgs []string
				if len(args) == 0 {
					strArgs = os.GetDefaultOS(ctx).Args()
				} else {
					var errObj *object.Error
					strArgs, errObj = object.AsStringSlice(args[0])
					if errObj != nil {
						return errObj
					}
				}
				if err := app.value.RunContext(ctx, strArgs); err != nil {
					return object.NewError(err)
				}
				return object.Nil
			}), true
	}
	return nil, false
}

func (app *App) SetAttr(name string, value object.Object) error {
	switch name {
	case "name":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.Name = str.Value()
	case "help_name":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.HelpName = str.Value()
	case "usage":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.Usage = str.Value()
	case "usage_text":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.UsageText = str.Value()
	case "args":
		b, ok := value.(*object.Bool)
		if !ok {
			return object.TypeErrorf("type error: expected bool, got %s", value.Type())
		}
		app.value.Args = b.Value()
	case "args_usage":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.ArgsUsage = str.Value()
	case "version":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.Version = str.Value()
	case "description":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.Description = str.Value()
	case "default_command":
		str, ok := value.(*object.String)
		if !ok {
			return object.TypeErrorf("type error: expected string, got %s", value.Type())
		}
		app.value.DefaultCommand = str.Value()
	case "enable_bash_completion":
		b, ok := value.(*object.Bool)
		if !ok {
			return object.TypeErrorf("type error: expected bool, got %s", value.Type())
		}
		app.value.EnableBashCompletion = b.Value()
	case "hide_help":
		b, ok := value.(*object.Bool)
		if !ok {
			return object.TypeErrorf("type error: expected bool, got %s", value.Type())
		}
		app.value.HideHelp = b.Value()
	case "hide_help_command":
		b, ok := value.(*object.Bool)
		if !ok {
			return object.TypeErrorf("type error: expected bool, got %s", value.Type())
		}
		app.value.HideHelpCommand = b.Value()
	case "hide_version":
		b, ok := value.(*object.Bool)
		if !ok {
			return object.TypeErrorf("type error: expected bool, got %s", value.Type())
		}
		app.value.HideVersion = b.Value()
	default:
		return object.TypeErrorf("type error: %s object has no attribute %q", APP, name)
	}
	return nil
}

func NewApp(ctx context.Context, opts *object.Map) (*App, error) {
	app := &ucli.App{}

	// Handle string fields
	for _, field := range []string{
		"name", "help_name", "usage", "usage_text", "args_usage",
		"version", "description", "default_command",
	} {
		if opt := opts.Get(field); opt != object.Nil {
			str, ok := opt.(*object.String)
			if !ok {
				return nil, object.TypeErrorf("type error: %s must be a string", field)
			}
			switch field {
			case "name":
				app.Name = str.Value()
			case "help_name":
				app.HelpName = str.Value()
			case "usage":
				app.Usage = str.Value()
			case "usage_text":
				app.UsageText = str.Value()
			case "args_usage":
				app.ArgsUsage = str.Value()
			case "version":
				app.Version = str.Value()
			case "description":
				app.Description = str.Value()
			case "default_command":
				app.DefaultCommand = str.Value()
			}
		}
	}

	// Handle boolean fields
	for _, field := range []string{
		"args", "enable_bash_completion", "hide_help",
		"hide_help_command", "hide_version",
	} {
		if opt := opts.Get(field); opt != object.Nil {
			b, ok := opt.(*object.Bool)
			if !ok {
				return nil, object.TypeErrorf("type error: %s must be a boolean", field)
			}
			switch field {
			case "args":
				app.Args = b.Value()
			case "enable_bash_completion":
				app.EnableBashCompletion = b.Value()
			case "hide_help":
				app.HideHelp = b.Value()
			case "hide_help_command":
				app.HideHelpCommand = b.Value()
			case "hide_version":
				app.HideVersion = b.Value()
			}
		}
	}

	// Handle commands
	if commandsOpt := opts.Get("commands"); commandsOpt != object.Nil {
		commands, ok := commandsOpt.(*object.List)
		if !ok {
			return nil, object.TypeErrorf("type error: commands must be a list")
		}
		app.Commands = []*ucli.Command{}
		for _, cmdOpt := range commands.Value() {
			cmd, ok := cmdOpt.(*Command)
			if !ok {
				return nil, object.TypeErrorf("type error: expected a command (got %s)", cmdOpt.Type())
			}
			app.Commands = append(app.Commands, cmd.value)
		}
	}

	// Handle flags
	if flagsOpt := opts.Get("flags"); flagsOpt != object.Nil {
		flags, ok := flagsOpt.(*object.List)
		if !ok {
			return nil, object.TypeErrorf("type error: flags must be a list")
		}
		app.Flags = []ucli.Flag{}
		for _, flagOpt := range flags.Value() {
			flag, ok := flagOpt.(*Flag)
			if !ok {
				return nil, object.TypeErrorf("type error: expected a flag (got %s)", flagOpt.Type())
			}
			app.Flags = append(app.Flags, flag.value)
		}
	}

	stdin := os.GetDefaultOS(ctx).Stdin()
	stdout := os.GetDefaultOS(ctx).Stdout()
	app.Reader = stdin
	app.Writer = stdout
	app.ErrWriter = stdout

	if actionOpt := opts.Get("action"); actionOpt != object.Nil {
		if action, ok := actionOpt.(*object.Function); ok {
			app.Action = func(c *ucli.Context) error {
				callFunc, ok := object.GetCallFunc(c.Context)
				if !ok {
					return fmt.Errorf("no call function found")
				}
				args := []object.Object{NewCtx(c)}
				if _, err := callFunc(c.Context, action, args); err != nil {
					return err
				}
				return nil
			}
		} else {
			return nil, errz.TypeErrorf("type error: action must be a function")
		}
	}

	return &App{value: app}, nil
}
