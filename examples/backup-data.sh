#!/bin/bash

# Backup Example
# This script demonstrates how to backup Elasticsearch data to a file

echo "=== Backup Example ==="
echo "Creating backup of Elasticsearch index..."

# Configuration
ES_HOST="http://localhost:9200"
INDEX_NAME="logs"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/${INDEX_NAME}_${TIMESTAMP}.ndjson"

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

echo "Backing up index '${INDEX_NAME}' to '${BACKUP_FILE}'..."

# Backup data
../bin/elasticdump backup \
    --input="${ES_HOST}/${INDEX_NAME}" \
    --output="${BACKUP_FILE}" \
    --format=ndjson \
    --scrollSize=5000 \
    --verbose

echo "Backup completed: ${BACKUP_FILE}"
echo "File size: $(du -h ${BACKUP_FILE} | cut -f1)"
