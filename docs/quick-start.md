# Quick Start

Here is how to get up and running with Tamarin either as a CLI or as a library.

If this is your first time trying Tamarin, we recommend starting with the CLI.

## Install using Homebrew

Install the Tamarin CLI using Homebrew as follows:

```
brew tap cloudcmds/tamarin
brew install tamarin
```

You should then be able to run `tamarin -h` to see usage information.

## Install the CLI from Source

If you have Go installed on your system, you can build and install by running:

```
go install github.com/cloudcmds/tamarin@latest
```

The `tamarin` binary should now be present in `$HOME/go/bin` or in the location
corresponding to your GOPATH directory.

## Add Tamarin as a Library

Use `go get` to add Tamarin as a library dependency of your Go project:

```
go get github.com/cloudcmds/tamarin@v0.0.14
```

## Run the REPL

Running the `tamarin` command without any options will start the REPL:

```
$ tamarin
Tamarin

>>> print("Hello gophers!")
Hello gophers!
>>>
```

## Execute a Tamarin Script String

Run `tamarin -c "code-to-execute"` to directly evaluate a given code string:

```
$ tamarin -c "uuid.v4()"
"0432500a-504a-435e-84de-16abf17b302f"
```

## Run a Script

To run a Tamarin script in a file, pass the path to the command:

```
$ tamarin ./examples/pipe.tm
GOPHERS ARE BURROWING RODENTS
```

## VSCode Extension

VSCode users can quickly enable Tamarin syntax highlighting by installing the
[Tamarin VSCode Extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.tamarin-language).
