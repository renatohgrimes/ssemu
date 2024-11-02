#!/bin/bash

echo "Creating release..."

rm -f release.zip
rm -r bin/prometheus
rm -r bin/otel
rm -r bin/elastic

sed -i '/OTEL_EXPORTER_OTLP_ENDPOINT/d' bin/ssemu.env

zip -r release-linux.zip bin -x "bin/ssemu.exe"
zip -r "release-windows.zip" bin -x "bin/ssemu"