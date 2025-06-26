#!/bin/bash

# Batch Processing Example
# This script demonstrates how to process multiple indices in batch

echo "=== Batch Processing Example ==="
echo "Processing multiple indices..."

# Configuration
SOURCE_HOST="http://localhost:9200"
DEST_HOST="http://localhost:9200"

# List of indices to migrate
INDICES=(
    "logs-2024-01"
    "logs-2024-02"
    "logs-2024-03"
    "metrics-2024-01"
    "metrics-2024-02"
)

echo "Migrating ${#INDICES[@]} indices from ${SOURCE_HOST} to ${DEST_HOST}..."

# Process each index
for index in "${INDICES[@]}"; do
    echo ""
    echo "Processing index: ${index}"
    echo "================================"
    
    # Transfer mappings
    echo "Transferring mappings for ${index}..."
    ../bin/elasticdump transfer \
        --input="${SOURCE_HOST}/${index}" \
        --output="${DEST_HOST}/${index}" \
        --type=mapping \
        --verbose
    
    if [ $? -ne 0 ]; then
        echo "Warning: Failed to transfer mappings for ${index}"
        continue
    fi
    
    # Transfer settings
    echo "Transferring settings for ${index}..."
    ../bin/elasticdump transfer \
        --input="${SOURCE_HOST}/${index}" \
        --output="${DEST_HOST}/${index}" \
        --type=settings \
        --verbose
    
    if [ $? -ne 0 ]; then
        echo "Warning: Failed to transfer settings for ${index}"
        continue
    fi
    
    # Transfer data
    echo "Transferring data for ${index}..."
    ../bin/elasticdump transfer \
        --input="${SOURCE_HOST}/${index}" \
        --output="${DEST_HOST}/${index}" \
        --type=data \
        --concurrency=3 \
        --scrollSize=2000 \
        --verbose
    
    if [ $? -eq 0 ]; then
        echo "✓ Successfully migrated ${index}"
    else
        echo "✗ Failed to migrate ${index}"
    fi
done

echo ""
echo "Batch processing completed!"
echo "Check the results by running document counts for each index."
