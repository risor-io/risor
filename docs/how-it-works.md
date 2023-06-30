# How It Works

Risor follows standard patterns for implementing a scripting language. This
includes parsing source code into an Abstract Syntax Tree (AST), compiling the
AST into bytecode, and then executing the bytecode on a lightweight virtual
machine.

In the first version of Risor, v1, the AST was executed more directly, which
was considerably slower. In v2, the added compiler and virtual machine improved
performance by over 100x.

The excellent book [Writing an Interpreter in Go](https://interpreterbook.com/)
was the original inspiration for the project.

## The Internals

Risor includes the following internal components:

- A [lexer](https://github.com/risor-io/risor/tree/main/lexer) which takes
  source code as input and produces a stream of tokens as output.
- A [parser](https://github.com/risor-io/risor/tree/main/parser) which takes
  tokens as an input and produces an abstract syntax tree (AST).
- A [compiler](https://github.com/risor-io/risor/tree/main/compiler) which
  compiles the AST into Risor bytecode instructions.
- A [virtual machine](https://github.com/risor-io/risor/tree/main/vm) which
  executes the program bytecode.
- [Built-in types](https://github.com/risor-io/risor/tree/main/object)
  available to all programs.
- [Built-in functions](https://github.com/risor-io/risor/blob/main/vm/builtins.go)
  that are accesible by default.

## Controlling Execution

The [exec](https://github.com/risor-io/risor/blob/main/exec/exec.go)
package offers a user-friendly API to use Risor as a library.
The provided `context.Context` is used to cancel execution or limit execution
with a timeout. Internally Risor passes this context to all operations to
guarantee that execution quickly stops when the context is canceled.

## Concurrency

A single Risor execution operates within a single goroutine. Multiple Risor
executions may happen concurrently and these are entirely independent. Risor
strictly avoids use of global state for safety and security reasons.
