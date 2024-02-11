# cli

Module `cli` is used to build command line programs written in the Risor
language. This wraps `github.com/urfave/cli` v2.

## Functions

### app

```go filename="Function signature"
app(opts map) app
```

Creates a new app with the given options. The app can then be run with its `run`
method.

```go copy filename="Example"
cli.app({
    name: "myapp",
    description: "My app description",
    commands: [],
}).run()
```

## Types

### app

App is the main entry point for a command line program. It contains commands,
is customized with flags, and is run with the `run` method.

### ctx

### command

### flag

