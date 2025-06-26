#!/bin/bash

# Limited Transfer Example
# This script demonstrates how to transfer a limited number of documents

echo "=== Limited Transfer Example ==="
echo "Transferring limited number of documents..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9200"
SOURCE_INDEX="large_dataset"
DEST_INDEX="sample_dataset"
LIMIT=10000  # Transfer only first 10,000 documents

echo "Transferring first ${LIMIT} documents from '${SOURCE_INDEX}' to '${DEST_INDEX}'..."

# Transfer with limit
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${SOURCE_INDEX}" \
    --output="${DEST_HOST}/${DEST_INDEX}" \
    --type=data \
    --limit=${LIMIT} \
    --concurrency=2 \
    --scrollSize=1000 \
    --verbose

echo "Limited transfer completed!"
echo "Transferred up to ${LIMIT} documents."
echo ""
echo "Verify the count:"
echo "curl -X GET '${DEST_HOST}/${DEST_INDEX}/_count?pretty'"
