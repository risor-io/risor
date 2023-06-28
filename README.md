# Risor

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/risor-io/risor/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/risor-io/risor/tree/main)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/risor-io/risor)
[![Go Report Card](https://goreportcard.com/badge/github.com/risor-io/risor?style=flat-square)](https://goreportcard.com/report/github.com/risor-io/risor)
[![Releases](https://img.shields.io/github/release/risor-io/risor/all.svg?style=flat-square)](https://github.com/risor-io/risor/releases)

A fast and flexible embedded scripting language for Go projects. Risor compiles
scripts to bytecode internally which it then runs on a lightweight Virtual
Machine (VM). Risor is written in pure Go.

Risor modules integrate the Go standard library, making it easy to write
scripts using functions that you're already familiar with as a Go developer.

## Notice: Project Renamed

Risor is a young project and until June 28, 2023 was known as _Tamarin_. For
various reasons, the project needed a new name that would take the project into
the future. Risor is a fun name, a bit shorter, and I can get a domain name for
the project. Thanks for bearing with me during this update!

## Documentation

Documentation is available at [risor-io.github.io/risor](https://risor-io.github.io/risor/).

## Getting Started

The [Quick Start](https://risor-io.github.io/risor/quick-start/) in the
documentation is where you should head to get started.

If you use Homebrew, you can install the Risor CLI as follows:

```
brew tap risor-io/risor
brew install risor
```

Having done that, just run `risor` to start the CLI or `risor -h` to see
usage information.

## Using Risor

Risor is designed to be versatile and accommodate a variety of usage patterns. You can leverage Risor in the following ways:

- **REPL**: Risor offers a Read-Evaluate-Print-Loop (REPL) that you can use to interactively write and test scripts. This is perfect for experimentation and debugging.

- **Library**: Risor can be imported as a library into existing Go projects. It provides a simple API for running scripts and interacting with the results, in isolated environments for sandboxing.

- **Executable script runner**: Risor scripts can also be marked as executable, providing a simple way to leverage Risor in your build scripts, automation, and other tasks.

- **API**: (Coming soon) A service and API will be provided for remotely executing and managing Risor scripts. This will allow integration into various web applications, potentially with self-hosted and a managed cloud version.

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

## Syntax Highlighting

A [Risor VSCode extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.tamarin-language)
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
