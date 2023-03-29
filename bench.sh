#!/bin/bash

pushd cmd/cli >/dev/null
go build
popd >/dev/null

for i in {1..10}; do
    ./cmd/cli/cli -timing ./examples/fib.tm
done
