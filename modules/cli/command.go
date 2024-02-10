package cli

import (
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	ucli "github.com/urfave/cli/v2"
)

const COMMAND object.Type = "cli.command"

type Command struct {
	value *ucli.Command
}

func (c *Command) Type() object.Type {
	return "cli.command"
}

func (c *Command) Inspect() string {
	return "cli.command"
}

func (c *Command) Interface() interface{} {
	return c.value
}

func (c *Command) IsTruthy() bool {
	return true
}

func (c *Command) Cost() int {
	return 0
}

func (c *Command) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("type error: unable to marshal %s", COMMAND)
}

func (c *Command) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", COMMAND, opType)
}

func (c *Command) Equals(other object.Object) object.Object {
	if other.Type() != "cli.command" {
		return object.False
	}
	return object.NewBool(c.value == other.(*Command).value)
}

func (c *Command) SetAttr(name string, value object.Object) error {
	var errObj *object.Error
	switch name {
	case "name":
		c.value.Name, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "aliases":
		c.value.Aliases, errObj = object.AsStringSlice(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "usage":
		c.value.Usage, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "usage_text":
		c.value.UsageText, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "description":
		c.value.Description, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "args":
		c.value.Args, errObj = object.AsBool(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "args_usage":
		c.value.ArgsUsage, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "category":
		c.value.Category, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "hide_help":
		c.value.HideHelp, errObj = object.AsBool(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "hide_help_command":
		c.value.HideHelpCommand, errObj = object.AsBool(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "hidden":
		c.value.Hidden, errObj = object.AsBool(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "use_short_option_handling":
		c.value.UseShortOptionHandling, errObj = object.AsBool(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "help_name":
		c.value.HelpName, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	case "custom_help_template":
		c.value.CustomHelpTemplate, errObj = object.AsString(value)
		if errObj != nil {
			return errObj.Value()
		}
	default:
		return fmt.Errorf("attribute error: %s object has no attribute %q", COMMAND, name)
	}
	return nil
}

func (c *Command) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "name":
		return object.NewString(c.value.Name), true
	case "aliases":
		return object.NewStringList(c.value.Aliases), true
	case "usage":
		return object.NewString(c.value.Usage), true
	case "usage_text":
		return object.NewString(c.value.UsageText), true
	case "description":
		return object.NewString(c.value.Description), true
	case "args":
		return object.NewBool(c.value.Args), true
	case "args_usage":
		return object.NewString(c.value.ArgsUsage), true
	case "category":
		return object.NewString(c.value.Category), true
	case "hide_help":
		return object.NewBool(c.value.HideHelp), true
	case "hide_help_command":
		return object.NewBool(c.value.HideHelpCommand), true
	case "hidden":
		return object.NewBool(c.value.Hidden), true
	case "use_short_option_handling":
		return object.NewBool(c.value.UseShortOptionHandling), true
	case "help_name":
		return object.NewString(c.value.HelpName), true
	case "custom_help_template":
		return object.NewString(c.value.CustomHelpTemplate), true
	case "subcommands":
		l := make([]object.Object, len(c.value.Subcommands))
		for i, cmd := range c.value.Subcommands {
			l[i] = NewCommand(cmd)
		}
		return object.NewList(l), true
	}
	return nil, false
}

func NewCommand(c *ucli.Command) *Command {
	return &Command{value: c}
}