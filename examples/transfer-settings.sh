#!/bin/bash

# Transfer Settings Example
# This script demonstrates how to transfer only index settings between clusters

echo "=== Transfer Settings Example ==="
echo "Transferring index settings..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9200"
INDEX_NAME="analytics"

echo "Transferring settings for index '${INDEX_NAME}'..."

# Transfer settings only
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --type=settings \
    --verbose

echo "Settings transfer completed!"
echo ""
echo "You can verify the settings were transferred by running:"
echo "curl -X GET '${DEST_HOST}/${INDEX_NAME}/_settings?pretty'"
