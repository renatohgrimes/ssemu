#!/bin/bash

echo "Initializing tests..."

set -e

cleanup() {
    echo "Cleaning up tests..."   
    docker rm test.server.ssemu -f
    docker network rm test.network.ssemu -f
}

trap cleanup EXIT

# test env

cp configs/ssemu.env bin
cp test/ssemu.test.env bin

# test data

rm -f bin/database.sqlite3
rm -f bin/database.sqlite3-shm
rm -f bin/database.sqlite3-wal

# logs

rm -rf bin/logs

echo "Running tests..."

# test network
docker network create test.network.ssemu 

# test server
docker run --net test.network.ssemu --name test.server.ssemu -d -v $PWD/bin:/ssemu/bin ssemu:latest 

# test runner
docker run \
    --net test.network.ssemu \
    --name test.runner.ssemu \
    --rm \
    -v $PWD:/src \
    -v /tmp/ssemu/test/go-build:/root/.cache/go-build \
    -e EMU_TESTING_HOST=test.server.ssemu \
    -e EMU_TESTING_CLIENT=/src/bin/client \
    -e EMU_TESTING_LOGS=/src/bin/logs_tests.txt \
    ssemu/builder:latest \
    go test -timeout 30s ./...
