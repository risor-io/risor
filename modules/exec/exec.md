# exec

The `exec` module is used to run external commands.

Like the underlying Go `os/exec` package, this module does not invoke the system
shell, expand glob patterns, or handle other shell features.

## Callable

The `exec` module itself is callable, using one of two signatures. The preferred
form is now:

```go filename="Function signature"
exec(args []string, opts map) result
```

The old form that is still supported for backwards compatibility is:

```go filename="Function signature"
exec(name string, args []string, opts map) result
```

This provides a shorthand way to build and run a command. The function returns a
`result` object containing the stdout and stderr produced by running the command.
The `opts` argument is optional.

```go copy filename="Example"
>>> exec(["echo", "TEST"]).stdout
byte_slice("TEST\n")
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| dir    | string                        | The working directory of the command.    |
| env    | map                           | The environment given to the command.    |
| stdin  | string, byte_slice, or reader | The standard input given to the command. |
| stdout | writer                        | The standard output destination.         |
| stderr | writer                        | The standard error destination.          |

## Functions

### command

Two signatures are supported for the `command` function. The preferred form is now:

```go filename="Function signature"
command(args []string) command
```

The old form that is still supported for backwards compatibility is:

```go filename="Function signature"
command(name string, args ...string) command
```

Creates a new command with the given name and arguments. The command can then
be executed with its `run`, `start`, `output`, or `combined_output` methods.
Before the command is run, its `path`, `dir`, and `env` attributes may be set.
Read more about the [command](#command-1) type below.

```go copy filename="Example"
>>> exec.command(["echo", "TEST"]).output()
byte_slice("TEST\n")
```

### look_path

```go filename="Function signature"
look_path(name string) string
```

Searches for the named executable in the directories contained in the PATH
environment variable. If the name contains a slash, it is tried directly,
without consulting the PATH. Otherwise, the result is the absolute path to
the named executable.

```go copy filename="Example"
>>> exec.look_path("echo")
"/bin/echo"
```

## Types

### command

Represents an external command that is being built and run.

#### Attributes

| Name            | Type              | Description                                                                   |
| --------------- | ----------------- | ----------------------------------------------------------------------------- |
| path            | string            | The path to the executable.                                                   |
| dir             | string            | The working directory of the command.                                         |
| env             | []string          | The environment given to the command.                                         |
| stdout          | byte_slice        | The standard output produced by the command.                                  |
| stderr          | byte_slice        | The standard error produced by the command.                                   |
| run             | func()            | Runs the command and waits for it to complete.                                |
| output          | func() byte_slice | Runs the command and returns its standard output.                             |
| combined_output | func() byte_slice | Runs the command and returns its combined standard output and standard error. |
| start           | func()            | Starts the command but does not wait for it to complete.                      |
| wait            | func()            | Waits for the command to exit.                                                |

#### Examples

```go copy filename="Example"
>>> c := exec.command("pwd")
>>> c.dir = "/dev"
>>> c.run()
>>> c.stdout
"/dev\n"
```

### result

Represents the result of running an external command.

#### Attributes

| Name   | Type       | Description                                  |
| ------ | ---------- | -------------------------------------------- |
| stdout | byte_slice | The standard output produced by the command. |
| stderr | byte_slice | The standard error produced by the command.  |
| pid    | int        | The process ID of the command.               |

#### Examples

```go copy filename="Example"
>>> result := exec("ls")
>>> result.stdout
"file1\nfile2\n"
```
