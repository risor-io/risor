# Quick Start

Here's how to get up and running with Risor as a CLI or as a library.
If this is your first time trying Risor, we recommend trying the CLI first.

## Install using Homebrew

Install the [Risor](https://formulae.brew.sh/formula/risor) v2 CLI using [Homebrew](https://brew.sh/) as follows:

```bash
brew install risor
```

You should then be able to run `risor -h` to see usage information.

## Install the CLI from Source

If you have Go installed on your system, you can build and install by running:

```bash
go install github.com/risor-io/risor/v2@latest
```

The `risor` binary should now be present in `$HOME/go/bin` or in the location
corresponding to your GOPATH directory.

## Add Risor as a Library

Use `go get` to add Risor as a library dependency of your Go project:

```bash
go get github.com/risor-io/risor/v2@v2.0.0-alpha.1
```

## Run the REPL

Running the `risor` command without any options will start the REPL:

```go
$ risor
Risor

>>> print("Hello gophers!")
Hello gophers!
>>>
```

Entering `ctrl+c` or `ctrl+d` will exit the program.

## Execute a Risor String

Run `risor -c "code-to-execute"` to directly evaluate a given code string:

```go
$ risor -c "uuid.v4()"
"0432500a-504a-435e-84de-16abf17b302f"
```

## Run a Script

To run a Risor script in a file, pass the path to the command.

```go title="example.tm" linenums="1"
my_array := ["gophers", "are", "burrowing", "rodents"]
sentence := my_array | strings.join(" ") | strings.to_upper
print(sentence)
```

With the above `example.tm` file on disk, run Risor as follows:

```bash
$ risor ./example.tm
GOPHERS ARE BURROWING RODENTS
```

## VSCode Extension

VSCode users can quickly enable syntax highlighting using the
[Risor VSCode Extension](https://marketplace.visualstudio.com/items?itemName=CurtisMyzie.risor-language).

## TextMate Grammar

A TextMate grammar file is available
[here](https://github.com/risor-io/risor/blob/main/vscode/syntaxes/risor.grammar.json).
This may help with syntax highlighting in other editors.
