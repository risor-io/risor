# Tamarin

Tamarin is an embedded scripting language for Go projects.

By integrating Tamarin into your existing Go program or library, you can add
dynamic scripting and behavior without requiring users to recompile your program.

You may also find that the Tamarin CLI can be handy for miscellaneous command
line scripting tasks, thanks to the simple, single binary Go build and convenient
syntax.

```go
["welcome", "to", "tamarin", "👋"] | strings.join(" ")
```

## Use Cases

- Allow users of your Go program to customize event processing without recompilation.
- Add customization hooks to any CLI written in Go.
- Enable users of your library to write scripts that interact with Go structs.
- Add dynamic behaviors to a Go web server to customize initialization or
  request handling.
- Extend game engines written in Go with a scripting interface.
- Sandbox execution of user scripts in a SaaS application.

## Why Choose Tamarin?

There are already some really handy embedded scripting languages for Go. Here is
the great list on [awesome-go](https://github.com/avelino/awesome-go#embeddable-scripting-languages).
Tamarin is different in a few ways and you can consider whether this makes it a
good match for your project:

- General purpose, but with built-in capabilities for HTTP requests and more.
- Familiar syntax for Go and Python developers.
- Exposes a portion of the Go standard library to scripts.
- Expressive and intuitive list, map, string, set, and time data types.
- Pipe expressions to easily express processing pipelines.
- First-class error handling mechanisms.
- Easily customizable built-in functions.

## Getting Started

Head over to [Quick Start](quick-start.md) for information on how to start using
Tamarin as a CLI or a library. There are also a variety of
[examples](https://github.com/cloudcmds/tamarin/tree/main/cmd) on the Github that
demonstrate using Tamarin as a library.
