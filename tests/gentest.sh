#!/bin/bash

TIMESTAMP=$(date +"%Y-%m-%d-%H-%M")
TEST_FILE="test-${TIMESTAMP}.tm"

cat <<EOF >>${TEST_FILE}
// UPDATE THIS FILE FOR YOUR TEST
// github issue: xyz
// expected value: 10
// expected type: int

var a = 10
EOF

echo "created test file: ${TEST_FILE}"
