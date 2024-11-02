#!/bin/bash

echo "Making Windows executable..."

DATE="$(date +%Y.%m.%d.%H%M)"

docker run \
    --name builder.ssemu \
    --rm \
    -e CC="x86_64-w64-mingw32-gcc" \
    -e CGO_ENABLED=1 \
    -e GOOS=windows \
    -e GOARCH=amd64 \
    -v $PWD:/src \
    -v /tmp/ssemu/build/go-build:/root/.cache/go-build \
    ssemu/builder:latest \
    go build -o bin/ssemu.exe -ldflags "-X ssemu/internal/app.Version=$DATE-windows -s -w" cmd/ssemu/main.go