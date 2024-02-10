package cli

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/risor-io/risor/os"
	ucli "github.com/urfave/cli/v2"
)

const APP object.Type = "cli.app"

type App struct {
	app      *ucli.App
	callFunc object.CallFunc
}

func (app *App) Type() object.Type {
	return "app"
}

func (app *App) Inspect() string {
	return "cli.app"
}

func (app *App) Interface() interface{} {
	return app.app
}

func (app *App) IsTruthy() bool {
	return app.app != nil
}

func (app *App) Cost() int {
	return 8
}

func (app *App) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", APP)
}

func (app *App) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", APP, opType)
}

func (app *App) Equals(other object.Object) object.Object {
	if other.Type() != "cli.app" {
		return object.False
	}
	return object.NewBool(app.app == other.(*App).app)
}

func (app *App) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", APP, name)
}

func (app *App) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "name":
		return object.NewString(app.app.Name), true
	case "usage":
		return object.NewString(app.app.Usage), true
	case "run":
		return object.NewBuiltin("cli.app.run", func(ctx context.Context, args ...object.Object) object.Object {
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
			if err := app.app.Run(strArgs); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "commands":
		var commands []object.Object
		for _, cmd := range app.app.Commands {
			commands = append(commands, NewCommand(cmd))
		}
		return object.NewList(commands), true
	}
	return nil, false
}

// func (app *App) Run()

func New(opts *object.Map, callFunc object.CallFunc) (*App, error) {

	app := &ucli.App{}

	if nameOpt := opts.Get("name"); nameOpt != object.Nil {
		if name, err := object.AsString(nameOpt); err == nil {
			app.Name = name
		}
	}

	if usageOpt := opts.Get("usage"); usageOpt != object.Nil {
		if usage, err := object.AsString(usageOpt); err == nil {
			app.Usage = usage
		}
	}

	if descriptionOpt := opts.Get("description"); descriptionOpt != object.Nil {
		if description, err := object.AsString(descriptionOpt); err == nil {
			app.Description = description
		}
	}

	if versionOpt := opts.Get("version"); versionOpt != object.Nil {
		if version, err := object.AsString(versionOpt); err == nil {
			app.Version = version
		}
	}

	if authorsOpt := opts.Get("authors"); authorsOpt != object.Nil {
		if authors, err := object.AsStringSlice(authorsOpt); err == nil {
			app.Authors = []*ucli.Author{}
			for _, author := range authors {
				app.Authors = append(app.Authors, &ucli.Author{Name: author})
			}
		}
	}

	if actionOpt := opts.Get("action"); actionOpt != object.Nil {
		if action, ok := actionOpt.(*object.Function); ok {
			app.Action = func(c *ucli.Context) error {
				args := []object.Object{}
				for _, arg := range c.Args().Slice() {
					args = append(args, object.NewString(arg))
				}
				if _, err := callFunc(context.Background(), action, args); err != nil {
					return err
				}
				return nil
			}
		} else {
			return nil, fmt.Errorf("type error: action must be a function")
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
					return nil, fmt.Errorf("type error: expected a command, got %T", command)
				}
			}
		}
	}

	obj := &App{
		app:      app,
		callFunc: callFunc,
	}

	// obj.app.Action = func(c *ucli.Context) error {

	return obj, nil
}
