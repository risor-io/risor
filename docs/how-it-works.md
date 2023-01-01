# How It Works

Tamarin includes an interpreter written in Go and uses an approach called
Pratt Parsing to parse expressions. The excellent book
[Writing an Interpreter in Go](https://interpreterbook.com/)
was the inspiration for the project.

## The Internals

As with many other interpreted languages, Tamarin includes:

- A [lexer](https://github.com/cloudcmds/tamarin/tree/main/lexer) which takes
  source code as input and produces a stream of tokens as output.
- A [parser](https://github.com/cloudcmds/tamarin/tree/main/parser) which takes
  tokens as an input and produces an abstract syntax tree (AST).
- An [evaluator](https://github.com/cloudcmds/tamarin/tree/main/evaluator) which
  executes an AST as a program.
- [Built-in types](https://github.com/cloudcmds/tamarin/tree/main/object)
  available to all programs.
- [Built-in functions](https://github.com/cloudcmds/tamarin/blob/main/evaluator/builtins.go#L601)
  that are accesible by default.

## Controlling Execution

The [exec](https://github.com/cloudcmds/tamarin/blob/main/exec/exec.go)
package offers a user-friendly API to use Tamarin as a library.
The provided `context.Context` is used to cancel execution or limit execution
with a timeout. Internally Tamarin passes this context to all operations to
guarantee that execution quickly stops when the context is canceled.

## Concurrency

A single Tamarin execution operates within a single goroutine. Multiple Tamarin
executions may happen concurrently and these are entirely independent. Tamarin
avoids all use of global state.

## Providing Input

When running Tamarin as a library, you can provide input data to the scripts by
passing in a scope that has been pre-populated with some variables. The scope
can be passed via the `exec.Opts` struct that is passed to an execution.
