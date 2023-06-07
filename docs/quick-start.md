# Quick Start

Here's how to get up and running with Tamarin as a CLI or as a library.
If this is your first time trying Tamarin, we recommend trying the CLI first.

## Install using Homebrew

Install the Tamarin v2 CLI using Homebrew as follows:

```bash
brew tap cloudcmds/tamarin
brew install tamarin
```

You should then be able to run `tamarin -h` to see usage information.

## Install the CLI from Source

If you have Go installed on your system, you can build and install by running:

```bash
go install github.com/cloudcmds/tamarin/v2@latest
```

The `tamarin` binary should now be present in `$HOME/go/bin` or in the location
corresponding to your GOPATH directory.

## Add Tamarin as a Library

Use `go get` to add Tamarin as a library dependency of your Go project:

```bash
go get github.com/cloudcmds/tamarin/v2@v2.0.0-alpha.1
```

## Run the REPL

Running the `tamarin` command without any options will start the REPL:

```go
$ tamarin
Tamarin

>>> print("Hello gophers!")
Hello gophers!
>>>
```

Entering `ctrl+c` or `ctrl+d` will exit the program.

## Execute a Tamarin String

Run `tamarin -c "code-to-execute"` to directly evaluate a given code string:

```go
$ tamarin -c "uuid.v4()"
"0432500a-504a-435e-84de-16abf17b302f"
```

## Run a Script

To run a Tamarin script in a file, pass the path to the command.

```go title="example.tm" linenums="1"
my_array := ["gophers", "are", "burrowing", "rodents"]
sentence := my_array | strings.join(" ") | strings.to_upper
print(sentence)
```

With the above `example.tm` file on disk, run Tamarin as follows:

```bash
$ tamarin ./example.tm
GOPHERS ARE BURROWING RODENTS
```

## VSCode Extension

VSCode users can quickly enable syntax highlighting using the
[Tamarin VSCode Extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.tamarin-language).

## TextMate Grammar

A TextMate grammar file is available
[here](https://github.com/cloudcmds/tamarin/blob/main/vscode/syntaxes/tamarin.grammar.json).
This may help with syntax highlighting in other editors.
