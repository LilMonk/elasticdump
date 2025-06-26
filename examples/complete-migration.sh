#!/bin/bash

# Complete Migration Example
# This script demonstrates a complete migration workflow including mappings, settings, and data

echo "=== Complete Migration Example ==="
echo "Performing complete migration (mappings + settings + data)..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9200"
INDEX_NAME="products"

echo "Migrating index '${INDEX_NAME}' from ${SOURCE_HOST} to ${DEST_HOST}..."

# Step 1: Transfer mappings
echo "Step 1/3: Transferring mappings..."
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --type=mapping \
    --verbose

if [ $? -ne 0 ]; then
    echo "Error: Failed to transfer mappings"
    exit 1
fi

# Step 2: Transfer settings
echo "Step 2/3: Transferring settings..."
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --type=settings \
    --verbose

if [ $? -ne 0 ]; then
    echo "Error: Failed to transfer settings"
    exit 1
fi

# Step 3: Transfer data
echo "Step 3/3: Transferring data..."
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --type=data \
    --concurrency=4 \
    --scrollSize=1000 \
    --verbose

if [ $? -ne 0 ]; then
    echo "Error: Failed to transfer data"
    exit 1
fi

echo "Complete migration finished successfully!"
echo ""
echo "Verification commands:"
echo "Source count: curl -X GET '${SOURCE_HOST}/${INDEX_NAME}/_count?pretty'"
echo "Dest count:   curl -X GET '${DEST_HOST}/${INDEX_NAME}/_count?pretty'"
