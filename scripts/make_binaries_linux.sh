#!/bin/bash

echo "Making Linux executable..."

DATE="$(date +%Y.%m.%d.%H%M)"

docker run \
    --name builder.ssemu \
    --rm \
    -e CGO_ENABLED=1 \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    -v $PWD:/src \
    -v /tmp/ssemu/build/go-build:/root/.cache/go-build \
    ssemu/builder:latest \
    go build -o bin/ssemu -ldflags "-X ssemu/internal/app.Version=$DATE-linux -s -w" cmd/ssemu/main.go
