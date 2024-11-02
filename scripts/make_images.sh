#!/bin/bash

echo "Creating builder docker image..." 
docker build --tag ssemu/builder:latest --file $PWD/build/builder.Dockerfile .

echo "Creating server docker image..." 
docker build --build-arg HOST_UID=$(id -u) --tag ssemu:latest --file $PWD/build/server.Dockerfile .