# Risor

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/risor-io/risor/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/risor-io/risor/tree/main)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
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

## Getting Started

Head over to [Getting Started](https://risor.io/docs) in the
documentation.

That said, if you use [Homebrew](https://brew.sh/), you can install the [Risor](https://formulae.brew.sh/formula/risor) CLI as follows:

```
brew install risor
```

Having done that, just run `risor` to start the CLI or `risor -h` to see
usage information.

Execute a code snippet directly using the `-c` option:

```go
risor -c "time.now()"
"2023-08-19T16:15:46-04:00"
```

## Quick Example

Here's a short example of how Risor feels like a hybrid of Go and Python, with
new features like pipe expressions for transformations, and with access to portions
of the Go standard library (like the `strings` package):

```go
array := ["gophers", "are", "burrowing", "rodents"]

sentence := array | strings.join(" ") | strings.to_upper

print(sentence)
```

Output:

```
GOPHERS ARE BURROWING RODENTS
```

## Built-in Functions and Modules

30+ built-in functions are included and are documented [here](https://risor.io/docs/builtins).

Modules are included that generally wrap the equivalent Go package. For example,
there is direct correspondence between `base64`, `bytes`, `json`, `math`, `os`,
`rand`, `strconv`, `strings`, and `time` Risor modules and the Go standard library.

Risor modules that are beyond the Go standard library include `aws`, `pgx`, and
`uuid`. Additional modules are being added regularly.

## Using Risor

Risor is designed to be versatile and accommodate a variety of usage patterns. You can leverage Risor in the following ways:

- **REPL**: Risor offers a Read-Evaluate-Print-Loop (REPL) that you can use to interactively write and test scripts. This is perfect for experimentation and debugging.

- **Library**: Risor can be imported as a library into existing Go projects. It provides a simple API for running scripts and interacting with the results, in isolated environments for sandboxing.

- **Executable script runner**: Risor scripts can also be marked as executable, providing a simple way to leverage Risor in your build scripts, automation, and other tasks.

- **API**: (Coming soon) A service and API will be provided for remotely executing and managing Risor scripts. This will allow integration into various web applications, potentially with self-hosted and a managed cloud version.

## Go Interface

It is trivial to embed Risor in your Go program in order to evaluate scripts
that have access to arbitrary Go structs and other types.

The simplest way to use Risor is to call the `Eval` function and provide the
Risor script source code. The result of the script is returned as a Risor object:

```go
result, err := risor.Eval(ctx, "math.min([5, 3, 7])")
min := result.(*object.Int).Value()
```

Provide input to the script using Risor options:

```go
result, err := risor.Eval(ctx, "input | strings.to_upper", risor.WithGlobal("input", "hello"))
fmt.Println(result) // HELLO
```

Use the same mechanism to inject a struct. You can then access fields or call
methods on the struct from the Risor script:

```go
type Example struct {
    Message string
}

ex := &Example{"abc"}

result, err := risor.Eval(ctx, "len(ex.Message)", risor.WithGlobal("ex", ex))
fmt.Println(result) // 3
```

## Syntax Highlighting

A [Risor VSCode extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.risor-language)
is already available which currently only offers syntax highlighting.

You can also make use of the [Risor TextMate grammar](./vscode/syntaxes/risor.grammar.json).

## Contributing

Risor is intended to be a community project. You can lend a hand in various ways:

- Please ask questions and share ideas in [GitHub discussions](https://github.com/risor-io/risor/discussions)
- Share Risor on any social channels that may appreciate it
- Open GitHub issue or a pull request for any bugs you find
- Star the project on GitHub

## Discuss the Project

Please visit the [GitHub discussions](https://github.com/risor-io/risor/discussions)
page to share thoughts and questions.

## Credits

Check [CREDITS.md](./CREDITS.md).

## License

Released under the [Apache License, Version 2.0](./LICENSE).

Copyright Curtis Myzie / [github.com/myzie](https://github.com/myzie).
