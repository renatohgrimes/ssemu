#!/bin/bash

echo "Creating base environment..."

mkdir -p bin
mkdir -p bin/client
mkdir -p bin/otel
mkdir -p bin/prometheus/config
mkdir -p bin/prometheus/data
mkdir -p bin/elastic

cp configs/ssemu.env bin
cp configs/prometheus.yml bin/prometheus/config
cp configs/otel-collector-config.yml bin/otel

echo "You must paste your client files here!" > bin/client/README.txt