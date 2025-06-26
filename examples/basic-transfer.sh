#!/bin/bash

# Basic Transfer Example
# This script demonstrates how to transfer data between two Elasticsearch clusters

echo "=== Basic Transfer Example ==="
echo "Transferring data from source index to destination index..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9200" 
SOURCE_INDEX="products"
DEST_INDEX="products_backup"

# Transfer data between clusters
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${SOURCE_INDEX}" \
    --output="${DEST_HOST}/${DEST_INDEX}" \
    --type=data \
    --concurrency=4 \
    --scrollSize=1000 \
    --verbose

echo "Transfer completed!"
