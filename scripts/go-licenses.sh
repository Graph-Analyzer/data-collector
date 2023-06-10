#!/bin/sh

GO_BINARY="$PWD/data-collector"
CONFIG_FILE="$PWD/lichen.yaml"
THIRD_PARTY_LICENSES_FILE="$PWD/THIRD-PARTY-LICENSES.txt"

if [ ! -f "$GO_BINARY" ]; then
    echo "$GO_BINARY does not exist, please first build the project."
    exit 1
fi

if [ ! -f "$CONFIG_FILE" ]; then
    echo "$CONFIG_FILE does not exist, something went wrong."
    exit 1
fi

echo "Start writing third party licenses."
echo "============================================"

lichen --config="$CONFIG_FILE" "$GO_BINARY" | tee "$THIRD_PARTY_LICENSES_FILE"

echo "============================================"
echo "Finished writing third party licenses."
