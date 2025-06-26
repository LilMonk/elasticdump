#!/bin/bash

# Authenticated Transfer Example
# This script demonstrates how to transfer data with authentication

echo "=== Authenticated Transfer Example ==="
echo "Transferring data with authentication..."

# Configuration
SOURCE_HOST="https://source-elastic.example.com:9200"
DEST_HOST="https://dest-elastic.example.com:9200"
INDEX_NAME="secure_logs"

# Authentication credentials
SOURCE_USERNAME="source_user"
SOURCE_PASSWORD="source_password"
DEST_USERNAME="dest_user"
DEST_PASSWORD="dest_password"

echo "Transferring authenticated data for index '${INDEX_NAME}'..."

# Note: For production use, consider using environment variables for credentials
# export ELASTICDUMP_SOURCE_USERNAME="source_user"
# export ELASTICDUMP_SOURCE_PASSWORD="source_password"

# Transfer with authentication
../bin/elasticdump transfer \
    --input="${SOURCE_HOST}/${INDEX_NAME}" \
    --output="${DEST_HOST}/${INDEX_NAME}" \
    --username="${SOURCE_USERNAME}" \
    --password="${SOURCE_PASSWORD}" \
    --type=data \
    --concurrency=2 \
    --scrollSize=500 \
    --verbose

echo "Authenticated transfer completed!"
echo ""
echo "Security Note:"
echo "For production environments, consider:"
echo "1. Using environment variables for credentials"
echo "2. Using API keys instead of username/password"
echo "3. Enabling SSL certificate verification"
