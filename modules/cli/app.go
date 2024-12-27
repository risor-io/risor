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

func (app *App) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: %s object has no attribute %q", APP, name)
}

func (app *App) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "name":
		return object.NewString(app.value.Name), true
	case "usage":
		return object.NewString(app.value.Usage), true
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
	case "commands":
		var commands []object.Object
		for _, cmd := range app.value.Commands {
			commands = append(commands, NewCommand(cmd))
		}
		return object.NewList(commands), true
	}
	return nil, false
}

func NewApp(ctx context.Context, opts *object.Map) (*App, error) {
	app := &ucli.App{}

	var err error
	app.Name, err = getMapStr(opts, "name")
	if err != nil {
		return nil, err
	}
	app.HelpName, err = getMapStr(opts, "help_name")
	if err != nil {
		return nil, err
	}
	app.Usage, err = getMapStr(opts, "usage")
	if err != nil {
		return nil, err
	}
	app.UsageText, err = getMapStr(opts, "usage_text")
	if err != nil {
		return nil, err
	}
	app.Description, err = getMapStr(opts, "description")
	if err != nil {
		return nil, err
	}
	app.Version, err = getMapStr(opts, "version")
	if err != nil {
		return nil, err
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

	if commandsOpt := opts.Get("commands"); commandsOpt != object.Nil {
		if commands, err := object.AsList(commandsOpt); err == nil {
			app.Commands = []*ucli.Command{}
			for _, commandOpt := range commands.Value() {
				switch command := commandOpt.(type) {
				case *Command:
					app.Commands = append(app.Commands, command.value)
				default:
					return nil, errz.TypeErrorf("type error: expected a command (got %s)", command.Type())
				}
			}
		}
	}

	if flagsOpt := opts.Get("flags"); flagsOpt != object.Nil {
		if flags, err := object.AsList(flagsOpt); err == nil {
			app.Flags = []ucli.Flag{}
			for _, flagOpt := range flags.Value() {
				flagObj, ok := flagOpt.(*Flag)
				if !ok {
					return nil, errz.TypeErrorf("type error: expected a flag (got %s)", flagOpt.Type())
				}
				app.Flags = append(app.Flags, flagObj.value)
			}
		}
	}

	return &App{value: app}, nil
}

func getMapStr(opts *object.Map, key string) (string, error) {
	if opt := opts.Get(key); opt != object.Nil {
		if str, err := object.AsString(opt); err == nil {
			return str, nil
		}
		return "", errz.TypeErrorf("type error: %s must be a string", key)
	}
	return "", nil
}
