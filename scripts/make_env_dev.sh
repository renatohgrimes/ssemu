#!/bin/bash

echo "Creating dev environment..."

mkdir -p bin
mkdir -p bin/client

cp configs/ssemu.env bin

sed -i '/OTEL_EXPORTER_OTLP_ENDPOINT/d' bin/ssemu.env

echo "You must paste your client files here!" > bin/client/README.txt

echo "\nEMU_TESTING_DATA=1" >> bin/ssemu.env