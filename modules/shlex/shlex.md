# shlex

The `shlex` module provides a shell-like argument parser.

The core functionality is provided by github.com/u-root/u-root/pkg/shlex.

## Functions

### argv

```go filename="Function signature"
argv(s string) []string
```
Split a command line according to usual simple shell rules. For example `start --append="foobar foobaz" --nogood 'food'`
will parse into the appropriate argvs to start the command.

```go filename="Example"
>>> shlex.argv("start --append=\"foobar foobaz\" --nogood 'food'")
["start", "--append=foobar foobaz", "--nogood", "food"]
```
