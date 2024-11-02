#!/bin/bash

echo "Running linter..."

docker run --name linter.ssemu --rm -v $PWD:/src -w /src golangci/golangci-lint:v1.60.3 golangci-lint run -v