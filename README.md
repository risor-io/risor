# Tamarin

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/cloudcmds/tamarin/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/cloudcmds/tamarin/tree/main)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/cloudcmds/tamarin)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudcmds/tamarin?style=flat-square)](https://goreportcard.com/report/github.com/cloudcmds/tamarin)
[![Releases](https://img.shields.io/github/release/cloudcmds/tamarin/all.svg?style=flat-square)](https://github.com/cloudcmds/tamarin/releases)

A fast and flexible embedded scripting language for Go projects. Tamarin code is
auto-compiled and runs within the lightweight Tamarin Virtual Machine, which is
written in pure Go.

Tamarin modules integrate the Go standard library, making it easy to write
scripts using standard library functions that you're already familiar with as
a Go developer.

Thanks to the new VM in v2, Tamarin is now **127x faster** than in v1!

## The Alpha Release of v2.0.0

As of June 2023, the Tamarin v2.0.0 alpha release has entered the testing phase.
I'm **seeking community feedback** before finalizing the exact language semantics
and feature set for v2 of the language. Please share your thoughts and questions
on the [discussions page](https://github.com/cloudcmds/tamarin/discussions).

The main change since v1 is the new compiler and virtual machine, which drastically
improved execution speed. While most language behaviors are unchanged, several
refinements have been made as well, necessitating the major revision to v2.
Specifically error propagation, the `try` built-in, and reduced use of the
`Result` type by other built-in functions all are major changes.

## Using Tamarin

Tamarin is designed to be versatile and accommodate a variety of usage patterns. You can leverage Tamarin in the following ways:

- **REPL**: Tamarin offers a Read-Evaluate-Print-Loop (REPL) that you can use to interactively write and test scripts. This is perfect for experimentation and debugging.

- **Library**: Tamarin can be imported as a library into existing Go projects. It provides a simple API for running scripts and interacting with the results, in isolated environments for sandboxing.

- **Executable script runner**: Tamarin scripts can also be marked as executable, providing a simple way to leverage Tamarin in your build scripts, automation, and other tasks.

- **API**: (Coming soon) A service and API will be provided for remotely executing and managing Tamarin scripts. This will allow integration into various web applications, potentially with self-hosted and a managed cloud version.

## Quick Example

Here's a short example of how Tamarin feels like a hybrid of Go and Python, with
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

## Documentation

Documentation is available at [cloudcmds.github.io/tamarin](https://cloudcmds.github.io/tamarin/).

## Getting Started

The [Quick Start](https://cloudcmds.github.io/tamarin/quick-start/) in the
documentation is where you should head to get started.

If you use Homebrew, you can install the Tamarin CLI as follows:

```
brew tap cloudcmds/tamarin
brew install tamarin
```

Having done that, just run `tamarin` to start the CLI or `tamarin -h` to see
usage information.

## Discuss the Project

Please visit the [GitHub discussions](https://github.com/cloudcmds/tamarin/discussions)
page to share thoughts and questions.

## Benchmark

Execution time in seconds for computing the 35th Fibonacci number:

![](bench/fib35.png?raw=true)

## Feature Overview

- Familiar syntax inspired by Go and Python
- Growing standard library which generally wraps the Go stdlib
- Currently libraries include: `json`, `math`, `rand`, `strings`, `time`, `uuid`, `strconv`, `pgx`
- Built-in types include: `set`, `map`, `list`, `result`, and more
- Functions are values and closures are supported
- Cancel Tamarin executions using Go contexts
- Library may be imported using the `import` keyword
- Easy HTTP requests via the `fetch` built-in function
- Pipeline expressions to create processing chains
- String templates similar to Python's f-strings

## Language Features and Syntax

See [docs/Features.md](./docs/Features.md).

## Syntax Highlighting

A [Tamarin VSCode extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.tamarin-language)
is already available which currently only offers syntax highlighting.

You can also make use of the [Tamarin TextMate grammar](./vscode/syntaxes/tamarin.grammar.json).

![](docs/assets/syntax-highlighting.png?raw=true)

## Contributing

Tamarin is intended to be a community project and you can lend a hand in various ways:

- Please ask questions and share ideas in [GitHub discussions](https://github.com/cloudcmds/tamarin/discussions)
- Share Tamarin on any social channels that may appreciate it
- Let us know about any usage of Tamarin that we can celebrate
- Leave us a star us on Github
- Open a pull request with fixes for bugs you find

## Known Issues & Limitations

- File I/O was intentionally omitted for now
- No concurrency support yet
- No user-defined types yet

## Credits

- [Thorsten Ball](https://github.com/mrnugget) and his book [Writing an Interpreter in Go](https://interpreterbook.com/).
- [Steve Kemp](https://github.com/skx) and the work in [github.com/skx/monkey](https://github.com/skx/monkey).
- [d5](https://github.com/d5) and the benchmarks in [github.com/d5/tengobench](https://github.com/d5/tengobench).

See more information in [CREDITS](./CREDITS).

## License

Released under the MIT License.
Copyright Curtis Myzie / [github.com/myzie](https://github.com/myzie).
