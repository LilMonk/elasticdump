#!/bin/bash

# Transfer Mappings Example
# This script demonstrates how to transfer only index mappings between clusters

echo "=== Transfer Mappings Example ==="
echo "Transferring index mappings..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9201"
INDEX_NAME="users"

echo "Transferring mappings for index '${INDEX_NAME}'..."

# Transfer mappings only
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --type=mapping \
    --verbose

echo "Mapping transfer completed!"
echo ""
echo "You can verify the mappings were transferred by running:"
echo "curl -X GET '${DEST_HOST}/${INDEX_NAME}/_mapping?pretty'"
