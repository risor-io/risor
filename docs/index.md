# Tamarin

Tamarin is an embedded scripting language for Go projects.

At a high level, Tamarin is especially useful for adding user defined functionality
to existing Go programs or libraries. Almost any program behavior can be made
extensible by integrating Tamarin as a library and allowing interactions between
Tamarin and your existing Go structs and interfaces.

You may also find that the Tamarin CLI can be handy for miscellaneous scripting
tasks, thanks to the self-contained binary Go build.

## Use Cases

- Allow users of your backend service to customize event processing without
  requiring a compilation step.
- Add customization hooks to any CLI written in Go.
- Enable users of your library to write scripts that can call methods and
  access fields on your Go structs.
- Add dynamic behaviors to a Go web server to customize initialization or
  request handling without recompiling the server.
- Extend game engines written in Go with a scripting interface.

## Why Choose Tamarin?

There are already some really handy embedded scripting languages for Go. Here is
the list on [awesome-go](https://github.com/avelino/awesome-go#embeddable-scripting-languages).
Tamarin is different in a few ways and you can consider whether this makes it a
good match for your project.

- General purpose, but with built-in capabilities for HTTP requests and more.
- Familiar syntax for Go and Python developers.
- Expose Go's standard library functionality to scripts.
- Expressive and powerful list, map, and set data types.
- Pipe expressions to easily express processing pipelines.
- First-class error handling mechanisms including a Result type.

## Getting Started

It's easiest to give Tamarin a try by installing the CLI, which offers a REPL.
There are also a number of [examples](https://github.com/cloudcmds/tamarin/tree/main/cmd)
in the Github repository that show how to use Tamarin as a library.

Head over to [Quick Start](quick-start.md) for some information on both options.
