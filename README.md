# Tamarin

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
