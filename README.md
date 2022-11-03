# Tamarin

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/myzie/tamarin/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/myzie/tamarin/tree/master)

Cloud scripting language.

## Usage

To execute a Tamarin script, pass the path of a script to the tamarin binary:

     $ tamarin ./example/hello.mon

Scripts can be made executable by adding a suitable shebang line:

     $ cat hello.mon
     #!/usr/bin/env tamarin
     print("Hello world!")

Execution then works as you would expect:

     $ chmod 755 hello.mon
     $ ./hello.mon
     Hello, world!

## Further Documentation

Work in progress. See [example.mon](./example.mon).

## Credits

- [Thorsten Ball](https://github.com/mrnugget) and his book [Writing an Interpreter in Go](https://interpreterbook.com/).
- [Steve Kemp](https://github.com/skx) and the work in [github.com/skx/monkey](https://github.com/skx/monkey).

See more information in [CREDITS](./CREDITS).
