#!/bin/bash

# Restore Example
# This script demonstrates how to restore data from a backup file

echo "=== Restore Example ==="
echo "Restoring data from backup file..."

# Configuration
ES_HOST="http://localhost:9200"
INDEX_NAME="logs_restored"
BACKUP_FILE="./backups/logs_20241226_143022.ndjson"

# Check if backup file exists
if [ ! -f "${BACKUP_FILE}" ]; then
    echo "Error: Backup file '${BACKUP_FILE}' not found!"
    echo "Please update the BACKUP_FILE variable with the correct path."
    exit 1
fi

echo "Restoring from '${BACKUP_FILE}' to index '${INDEX_NAME}'..."

# Restore data
../bin/elasticdump restore \
    --input="${BACKUP_FILE}" \
    --output="${ES_HOST}/${INDEX_NAME}" \
    --concurrency=2 \
    --verbose

echo "Restore completed!"
echo "You can verify the restore by checking the document count:"
echo "curl -X GET '${ES_HOST}/${INDEX_NAME}/_count?pretty'"
