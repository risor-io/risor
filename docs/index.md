# Tamarin

Tamarin is an embedded scripting language for Go projects.

At a high level, Tamarin is especially useful for adding user defined functionality
to existing Go programs or libraries. Almost any feature can be made extensible
by integrating Tamarin as a library, since Tamarin scripts can interact with Go
structs that already exist in your program.

You may also find that the Tamarin CLI can be handy for miscellaneous scripting
tasks, thanks to the self-contained binary Go build.

## Use Cases

- Allow users of your Go program to customize event processing without recompilation.
- Add customization hooks to any CLI written in Go.
- Enable users of your library to write scripts that interact with Go structs.
- Add dynamic behaviors to a Go web server to customize initialization or
  request handling.
- Extend game engines written in Go with a scripting interface.

## Why Choose Tamarin?

There are already some really handy embedded scripting languages for Go. Here is
the great list on [awesome-go](https://github.com/avelino/awesome-go#embeddable-scripting-languages).
Tamarin is different in a few ways and you can consider whether this makes it a
good match for your project:

- General purpose, but with built-in capabilities for HTTP requests and more.
- Familiar syntax for Go and Python developers.
- Exposes Go's standard library functionality to scripts.
- Expressive and intuitive list, map, string, set, and time data types.
- Pipe expressions to easily express processing pipelines.
- First-class error handling mechanisms including a Result type.

## Getting Started

Head over to [Quick Start](quick-start.md) for information on how to start using
Tamarin as a CLI or a library.

There are also a variety of [examples](https://github.com/cloudcmds/tamarin/tree/main/cmd)
in the Github repository that demonstrate using Tamarin as a library.
