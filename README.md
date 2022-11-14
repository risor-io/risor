# Tamarin

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/cloudcmds/tamarin/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/cloudcmds/tamarin/tree/main)

Cloud scripting language.

## Usage

To execute a Tamarin script, pass the path of a script to the tamarin binary:

     $ tamarin ./example/hello.tm

Scripts can be made executable by adding a suitable shebang line:

     $ cat hello.tm
     #!/usr/bin/env tamarin
     print("Hello world!")

Execution then works as you would expect:

     $ chmod 755 hello.tm
     $ ./hello.tm
     Hello world!

## Further Documentation

Work in progress. See [example.tm](./example.tm).

## Credits

- [Thorsten Ball](https://github.com/mrnugget) and his book [Writing an Interpreter in Go](https://interpreterbook.com/).
- [Steve Kemp](https://github.com/skx) and the work in [github.com/skx/monkey](https://github.com/skx/monkey).

See more information in [CREDITS](./CREDITS).
