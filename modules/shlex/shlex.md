import { Callout } from 'nextra/components';

# shlex

<Callout type="info" emoji="ℹ️">
  This module requires that Risor has been compiled with the `shlex` Go build tag.
  When compiling **manually**, [make sure you specify `-tags shlex`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source).
</Callout>

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
