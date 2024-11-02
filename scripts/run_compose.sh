#!/bin/bash

echo "Composing infrastructure..."

docker compose -p ssemu -f $PWD/deployments/docker-compose.yml up -d --no-build --force-recreate