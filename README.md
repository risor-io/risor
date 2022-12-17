# Tamarin

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/cloudcmds/tamarin/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/cloudcmds/tamarin/tree/main)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/cloudcmds/tamarin)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudcmds/tamarin?style=flat-square)](https://goreportcard.com/report/github.com/cloudcmds/tamarin)
[![Releases](https://img.shields.io/github/release/cloudcmds/tamarin/all.svg?style=flat-square)](https://github.com/cloudcmds/tamarin/releases)

A fun and pragmatic scripting language written in Go. May be used as a CLI or embedded as a library.

## Getting Started

Install the Tamarin CLI as follows:

```
go install github.com/cloudcmds/tamarin@latest
```

You can then use the CLI in three different ways. Running `tamarin`
by itself runs an interaction REPL:

```bash
tamarin
```

Pass a short script to execute using the `-c` flag:

```bash
tamarin -c 'uuid.v4()'
```

Or pass the path to a Tamarin script to execute:

```bash
tamarin ./examples/hello.tm
```

To use Tamarin as a library, see the example
[cmd/simple-example/main.go](./cmd/simple-example/main.go).

## Project Status

This project is early so Tamarin should be considered **alpha** software.
Please visit the [Github discussions](https://github.com/cloudcmds/tamarin/discussions)
page to share thoughts and questions! Feedback now is useful to guide the
project as it matures.

## Goals

Tamarin is meant to be useful for cloud, web, and general scripting.
The syntax should be familiar to Go, Typescript, and Python developers.
Tamarin should be "batteries included" with a handy standard library.

Miscellaneous popular Go open source libraries may be integrated into the
Tamarin standard library. Please share suggestions for what you think would
be good additions on this front!

We are using Tamarin to evaluate user-submitted scripts in sandboxes, but
there are many other situations where Tamarin may be useful. Because of this
initial use case, we chose to not support file I/O operations at this point,
however those will be added as an opt-in feature.

It is **not** currently considered important to be performance competitive
with widely used scripting languages. If you are looking for top performance
you should probably be using sandboxed V8/Javascript or Lua.

## Feature Overview

- Familiar syntax inspired by Go, Typescript, and Python
- Growing standard library which generally wraps the Go stdlib
- Includes higher level libraries that are beyond the Go stdlib
- Currently libraries include: `json`, `math`, `rand`, `strings`, `time`, `uuid`, `strconv`, `sql`
- Built-in types include: `set`, `map`, `list`, `result`, and more
- Functions are values; closures are supported
- Cancel evaluation using Go contexts
- Library may be imported using the `import` keyword
- Easy HTTP requests via the `fetch` built-in function
- Pipe expressions to create processing chains
- Error handling inspired by Rust, using a Result type
- String templates similar to Python's f-strings

## Making Tamarin Scripts Executable

Make a Tamarin script executable by adding a shebang line (`#!`). For example,
for a script named `hello.tm`:

     $ cat hello.tm
     #!/usr/bin/env tamarin
     print("Hello world!")

Execution then works as you would expect:

     $ chmod 755 hello.tm
     $ ./hello.tm
     Hello world!

## Language Features and Syntax

See [docs/Features.md](./docs/Features.md).

## Syntax Highlighting

A [Tamarin VSCode extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.tamarin-language)
is already available which currently only offers syntax highlighting.

You can also make use of the [Tamarin TextMate grammar](./vscode/syntaxes/tamarin.grammar.json).

![](docs/assets/syntax-highlighting.png?raw=true)

## Further Documentation

Work in progress. See [assorted.tm](./examples/assorted.tm).

## Contributing

🎉 This project is just getting started. Tamarin is intended to be a community
project and you can lend a hand in many different ways.

- Please ask questions and share ideas in [Github discussions](https://github.com/cloudcmds/tamarin/discussions)
- Share Tamarin on any social channels that may appreciate it
- Let us know about any usage of Tamarin that we can celebrate
- Leave us a star us on Github

## Known Issues & Limitations

- File I/O was intentionally omitted for now
- No concurrency support yet
- No user-defined types yet

## Pipe Expressions

A single value may be passed between successive calls using pipe expressions.

```
var array = ["gophers", "are", "burrowing", "rodents"]

var sentence = array | strings.join(" ") | strings.to_upper

print(sentence)
```

Output:

```
GOPHERS ARE BURROWING RODENTS
```

The intent is that if any call in the pipeline returns a `Result` object, that
object is unwrapped before providing it to the next call (or the pipeline aborts
if the Result is an error). This is not yet implemented.

## Error Handling in Tamarin

There are two categories of errors in Tamarin. Severe errors which indicate a programming mistake
create `*object.Error` errors that abort evaluation immediately. These are similar to Go panics.
Operations that may fail return an `*object.Result` object that contains either an `Ok` value or
an `Err` value. This is inspired by [Rust's Result type](https://doc.rust-lang.org/std/result/).
`Result` objects proxy to the `Ok` value in certain situations, to keep adhoc scripting concise.

### Result Methods

- `is_err`: returns `true` if the Result contains an error
- `is_ok`: returns `true` if the Result contains the successful operation value
- `unwrap`: returns the ok value if present and errors otherwise
- `unwrap_or`: returns the ok value if present and the default value otherwise
- `expect`: returns the ok value if present or raises an error with the given message

Other method calls on the `Result` object proxy to calls on the ok value (or
raise an error).

Altogether this means programmers have reliable tools for handling errors as
values without requiring `try: catch` style exceptions. In quick scripts, you
can rely on `Result` method proxying if you want, while in situations requiring
robustness you can use the `Result` methods to check for and unwrap errors.

## Credits

- [Thorsten Ball](https://github.com/mrnugget) and his book [Writing an Interpreter in Go](https://interpreterbook.com/).
- [Steve Kemp](https://github.com/skx) and the work in [github.com/skx/monkey](https://github.com/skx/monkey).

See more information in [CREDITS](./CREDITS).

## License

Released under the MIT License. Copyright Curtis Myzie / [github.com/myzie](https://github.com/myzie).
