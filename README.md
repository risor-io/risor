<h1><img src="https://github.com/risor-io/risor/raw/main/static/images/risor-logo-nopad.png" alt="Risor logo" height="64" valign="middle"> Risor</h1>

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/risor-io/risor/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/risor-io/risor/tree/main)
[![Apache-2.0 license](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/risor-io/risor)
[![Go Report Card](https://goreportcard.com/badge/github.com/risor-io/risor?style=flat-square)](https://goreportcard.com/report/github.com/risor-io/risor)
[![Releases](https://img.shields.io/github/release/risor-io/risor/all.svg?style=flat-square)](https://github.com/risor-io/risor/releases)

Risor is a fast and flexible scripting language for Go developers and DevOps.

Its modules integrate the Go standard library, making it easy to use functions
that you're already familiar with as a Go developer.

Scripts are compiled to bytecode and then run on a lightweight virtual machine.
Risor is written in pure Go.

## Documentation

Documentation is available at [risor.io](https://risor.io).

You might also want to try evaluating Risor scripts [from your browser](https://risor.io/#editor).

## Syntax Example

Here's a short example of how Risor feels like a hybrid of Go and Python. This
demonstrates using Risor's pipe expressions to apply a series of transformations:

```go
array := ["gophers", "are", "burrowing", "rodents"]

sentence := array | strings.join(" ") | strings.to_upper

print(sentence)
```

Output:

```
GOPHERS ARE BURROWING RODENTS
```

## Getting Started

You might want to head over to [Getting Started](https://risor.io/docs) in the
documentation. That said, here are tips for both the CLI and the Go library.

### Risor CLI and REPL

If you use [Homebrew](https://brew.sh/), you can install the
[Risor](https://formulae.brew.sh/formula/risor) CLI as follows:

```
brew install risor
```

Having done that, just run `risor` to start the CLI or `risor -h` to see
usage information.

Execute a code snippet directly using the `-c` option:

```go
risor -c "time.now()"
```

Start the REPL by running `risor` with no options.

### Build and Install the CLI from Source

Build the CLI from source as follows:

```bash
git clone git@github.com:risor-io/risor.git
cd risor/cmd/risor
go install -tags aws,k8s,vault .
```

### Go Library

Use `go get` to add Risor as a dependency of your Go program:

```bash
go get github.com/risor-io/risor
```

Here's an example of using the `risor.Eval` API to evaluate some code:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/risor-io/risor"
)

func main() {
	ctx := context.Background()
	script := "math.sqrt(input)"
	result, err := risor.Eval(ctx, script, risor.WithGlobal("input", 4))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The square root of 4 is:", result)
}
```

## Built-in Functions and Modules

30+ built-in functions are included and are documented [here](https://risor.io/docs/builtins).

Modules are included that generally wrap the equivalent Go package. For example,
there is direct correspondence between `base64`, `bytes`, `filepath`, `json`, `math`, `os`,
`rand`, `regexp`, `strconv`, `strings`, and `time` Risor modules and
the Go standard library.

Risor modules that are beyond the Go standard library currently include
`aws`, `color`, `cli`, `jmespath`, `pgx`, `playwright`, `qrcode`, `uuid`, `vault`, `k8s`, and more.

## Go Interface

It is trivial to embed Risor in your Go program in order to evaluate scripts
that have access to arbitrary Go structs and other types.

The simplest way to use Risor is to call the `Eval` function and provide the script source code.
The result is returned as a Risor object:

```go
result, err := risor.Eval(ctx, "math.min([5, 2, 7])")
// result is 2, as an *object.Int
```

Provide input to the script using Risor options:

```go
result, err := risor.Eval(ctx, "input | strings.to_upper", risor.WithGlobal("input", "hello"))
// result is "HELLO", as an *object.String
```

Use the same mechanism to inject a struct. You can then access fields or call
methods on the struct from the Risor script:

```go
type Example struct {
    Message string
}
example := &Example{"abc"}
result, err := risor.Eval(ctx, "len(ex.Message)", risor.WithGlobal("ex", example))
// result is 3, as an *object.Int
```

## Optional Modules

Risor is designed to have minimal external dependencies in its core libraries.
You can choose to opt into various add-on modules if they are of value in your
application. The modules are present in this same Git repository, but must be
installed with `go get` as separate dependencies:

| Name           | Path                                               | Go Get Command                                            |
| -------------- | -------------------------------------------------- | --------------------------------------------------------- |
| aws            | [modules/aws](./modules/aws)                       | `go get github.com/risor-io/risor/modules/aws`            |
| bcrypt         | [modules/bcrypt](./modules/bcrypt)                 | `go get github.com/risor-io/risor/modules/bcrypt`         |
| cli            | [modules/cli](./modules/cli)                       | `go get github.com/risor-io/risor/modules/cli`            |
| color          | [modules/color](./modules/color)                   | `go get github.com/risor-io/risor/modules/color`          |
| echarts        | [modules/echarts](./modules/echarts)               | `go get github.com/risor-io/risor/modules/echarts`          |
| gha            | [modules/gha](./modules/gha)                       | `go get github.com/risor-io/risor/modules/gha`            |
| goquery        | [modules/goquery](./modules/goquery)               | `go get github.com/risor-io/risor/modules/goquery`        |
| htmltomarkdown | [modules/htmltomarkdown](./modules/htmltomarkdown) | `go get github.com/risor-io/risor/modules/htmltomarkdown` |
| image          | [modules/image](./modules/image)                   | `go get github.com/risor-io/risor/modules/image`          |
| isatty         | [modules/isatty](./modules/isatty)                 | `go get github.com/risor-io/risor/modules/isatty`         |
| jmespath       | [modules/jmespath](./modules/jmespath)             | `go get github.com/risor-io/risor/modules/jmespath`       |
| k8s            | [modules/kubernetes](./modules/kubernetes)         | `go get github.com/risor-io/risor/modules/kubernetes`     |
| pgx            | [modules/pgx](./modules/pgx)                       | `go get github.com/risor-io/risor/modules/pgx`            |
| playwright     | [modules/playwright](./modules/playwright)         | `go get github.com/risor-io/risor/modules/playwright`     |
| qrcode         | [modules/qrcode](./modules/qrcode)                 | `go get github.com/risor-io/risor/modules/qrcode`         |
| s3fs           | [os/s3fs](./os/s3fs)                               | `go get github.com/risor-io/risor/os/s3fs`                |
| sched          | [modules/sched](./modules/sched)                   | `go get github.com/risor-io/risor/modules/sched`          |
| semver         | [modules/semver](./modules/semver)                 | `go get github.com/risor-io/risor/modules/semver`         |
| shlex          | [modules/shlex](./modules/shlex)                   | `go get github.com/risor-io/risor/modules/shlex`          |
| slack          | [modules/slack](./modules/slack)                   | `go get github.com/risor-io/risor/modules/slack`          |
| sql            | [modules/sql](./modules/sql)                       | `go get github.com/risor-io/risor/modules/sql`            |
| tablewriter    | [modules/tablewriter](./modules/tablewriter)       | `go get github.com/risor-io/risor/modules/tablewriter`    |
| template       | [modules/template](./modules/template)             | `go get github.com/risor-io/risor/modules/template`       |
| uuid           | [modules/uuid](./modules/uuid)                     | `go get github.com/risor-io/risor/modules/uuid`           |
| yaml           | [modules/yaml](./modules/yaml)                     | `go get github.com/risor-io/risor/modules/yaml`           |
| vault          | [modules/vault](./modules/vault)                   | `go get github.com/risor-io/risor/modules/vault`          |

These add-ons are included by default when using the Risor CLI. However, when
building Risor into your own program, you'll need to opt-in using `go get` as
described above and then add the modules as globals in Risor scripts as follows:

```go
import (
    "github.com/risor-io/risor"
    "github.com/risor-io/risor/modules/aws"
    "github.com/risor-io/risor/modules/image"
    "github.com/risor-io/risor/modules/pgx"
    "github.com/risor-io/risor/modules/uuid"
)

func main() {
    source := `"nice modules!"`
    result, err := risor.Eval(ctx, source,
        risor.WithGlobals(map[string]any{
            "aws":   aws.Module(),
            "image": image.Module(),
            "pgx":   pgx.Module(),
            "uuid":  uuid.Module(),
        }))
    // ...
}
```

## Syntax Highlighting

A [Risor VSCode extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.risor-language)
is already available which currently only offers syntax highlighting.

You can also make use of the [Risor TextMate grammar](./vscode/syntaxes/risor.grammar.json).

## Benchmarking

There are two Makefile commands that assist with benchmarking and CPU profiling:

```
make bench
make pprof
```

## Contributing

Risor is intended to be a community project. You can lend a hand in various ways:

- Please ask questions and share ideas in [GitHub discussions](https://github.com/risor-io/risor/discussions)
- Share Risor on any social channels that may appreciate it
- Open GitHub issue or a pull request for any bugs you find
- Star the project on GitHub

### Contributing New Modules

Adding modules to Risor is a great way to get involved with the project.
See [this guide](https://risor.io/docs/contributing_modules) for details.

## Community Projects

- [RSX: Package Risor Scripts into Go Binaries](https://github.com/rubiojr/rsx)
- [Awesome Risor](https://github.com/rubiojr/awesome-risor)
- [tree-sitter-risor](https://github.com/applejag/tree-sitter-risor)
- [bench_go_scripting](https://github.com/mna/bench_go_scripting)

## Discuss the Project

Please visit the [GitHub discussions](https://github.com/risor-io/risor/discussions)
page to share thoughts and questions.

There is also a `#risor` Slack channel on the [Gophers Slack](https://gophers.slack.com).

## Credits

Check [CREDITS.md](./CREDITS.md).

## License

Released under the [Apache License, Version 2.0](./LICENSE).

Copyright Curtis Myzie / [github.com/myzie](https://github.com/myzie).
