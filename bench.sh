#!/bin/bash

go build

for i in {1..10}; do
    ./tamarin -timing ./examples/fib.tm
done
