import { Callout } from 'nextra/components';

# cli

Module `cli` is used to build command line apps written in the Risor. Common
CLI features are supported, including commands, flags, arguments, usage, and
automatic help generation.

<Callout type="info" emoji="ℹ️">
  This module is included by default in the Risor CLI, but must be
  independently installed when using Risor as
  a library using `go get github.com/risor-io/risor/modules/cli`
</Callout>

Behind the scenes, this module uses the [urfave/cli](https://cli.urfave.org/) library.

## Getting Started

Create a file named `myapp.risor` with the following contents. Note that you
must include the shebang line including `--` at the top of the file to ensure
that arguments and options are _passed to the app_, rather than being used as
options by the Risor binary itself.

```risor copy filename="myapp.risor"
#!/usr/bin/env risor --

from cli import app, command as c

app({
    name: "myapp",
    description: "My app description",
    commands: [
        c({
            name: "hello",
            description: "Say hello",
            action: func(ctx) {
                print("Hello, world!")
            },
        }),
    ],
}).run()
```

Now make the file executable:

```
$ chmod +x ./myapp.risor
```

You can now run the app as follows:

```
$ ./myapp.risor hello
Hello, world!
```

## Functions

### app

```go filename="Function signature"
app(options map) app
```

Returns a new app initialized with the given options. A simple app may consist
of just a `name`, `description`, and `action` function. Call `.run()` on the
app to run it.

```risor copy filename="Example"
app := cli.app({
    name: "myapp",
    description: "My app description",
    action: func(ctx) {
        print("Hello, world!")
    },
})

app.run()
```

The `app` function supports the following options:

- `action func(ctx)`: The action to run when the app is run.
- `commands []cli.command`: A list of commands that the app supports.
- `description string`: A short description of the app.
- `flags []cli.flag`: A list of flags that the app supports.
- `help_name string` : Override for the name of the app in help output.
- `name string`: The name of the app.
- `usage_text string`: The usage text for the app.
- `usage string`: The usage string for the app.
- `version string`: The version of the app.

### command

```go filename="Function signature"
command(options map) command
```

Returns a new command initialized with the given options. Commands are provided
to an app via the app's `commands` option.

```go copy filename="Example"
command := cli.command({
    name: "add",
    description: "Add numbers provided as arguments",
    action: func(ctx) {
        sum := 0
        for _, arg := range ctx.args() {
            sum += int(arg)
        }
        print(sum)
    },
})
```

The `command` function supports the following options:

- `action func(ctx)`: The function to call when the command is invoked.
- `aliases []string`: A list of aliases for the command.
- `args_usage string`: A short description of the arguments of this command.
- `args bool`: Whether this command supports arguments.
- `category string`: The category the command is part of.
- `custom_help_template string`: Text template for the command help topic.
- `description string`: A longer explanation of how the command works.
- `flags []cli.flag`: A list of flags that the command supports.
- `help_name string`: Full name of command for help, defaults to full command name, including parent commands.
- `hidden bool`: Hide this command from help or completion.
- `hide_help_command bool`: Whether to hide the command from the help command.
- `hide_help bool` Hide the built-in help command and help flag.
- `name string`: The name of the command.
- `usage_text string`: Custom text to show in USAGE section of help.
- `usage string`: A short description of the usage of this command.
- `use_short_option_handling bool`: Enables short-option handling so the user can combine several single-character bool flags into one.

### flag

```go filename="Function signature"
flag(options map) flag
```

Returns a flag that may be used with an app or command. Supported flag types
include `string`, `int`, `bool`, `float`, `string_slice`, `int_slice`,
and `float_slice`.

A default value for the flag may be provided using the `value` option. The flag
type is inferred from the `value` option if a `type` is not specified.

```go copy filename="Example string flag"
cli.flag({
    name: "food",
    aliases: ["f"],
    usage: "The type of food to eat",
    env_vars: ["FOOD"], // read from this environment variable, if present
    value: "pizza",     // default value
    type: "string",     // flag type: string, int, bool, etc.
})
```

```go copy filename="Example int flag"
cli.flag({
    name: "count",
    aliases: ["c"],
    usage: "The number of items to process",
    value: 1,
})
```

```go copy filename="Example bool flag"
cli.flag({
    name: "verbose",
    aliases: ["v"],
    usage: "Enable verbose output",
    value: false,
})
```

## Types

### app

An app represents the main entry point for a command-line program. It contains
commands, is customized with flags, and is executed via `app.run()`.

### ctx

A ctx object is passed through to each handler action in a cli app. It is
used to retrieve context-specific args and parsed command-line options.

Attributes on the ctx object include:

| Name             | Type                         | Description                                                             |
| ---------------- | ---------------------------- | ----------------------------------------------------------------------- |
| args             | func() []string              | Returns the command-line arguments                                      |
| narg             | func() int                   | Returns the number of arguments                                         |
| value            | func(name string) object     | Returns the value of the flag corresponding to `name`                   |
| count            | func(name string) int        | Returns the count of the flag corresponding to `name`                   |
| flag_names       | func() []string              | Returns the names of all flags used by this context and parent contexts |
| local_flag_names | func() []string              | Returns the names of all flags used by this context                     |
| is_set           | func(name string) bool       | Returns true if the flag corresponding to `name` is set                 |
| set              | func(name string, value obj) | Sets the value of the flag corresponding to `name`                      |
| num_flags        | func() int                   | Returns the number of flags set                                         |
| bool             | func(name string) bool       | Returns the value of the bool flag corresponding to `name`              |
| int              | func(name string) int        | Returns the value of the int flag corresponding to `name`               |
| string           | func(name string) string     | Returns the value of the string flag corresponding to `name`            |
| string_slice     | func(name string) []string   | Returns the value of the string slice flag corresponding to `name`      |

### command

A command represents a sub-command of an app. It contains its own flags and
has an associated action. Commands may have sub-commands.

### flag

A flag is used to parse command-line flags in a cli app. Flags may be specified
on a cli app directly, as well as on commands.
