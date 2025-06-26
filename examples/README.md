# Elasticdump Examples

This directory contains example scripts demonstrating various use cases for the elasticdump CLI tool.

## Prerequisites

Before running these examples, make sure:

1. **Build the elasticdump binary**:
   ```bash
   cd ..
   go build -o bin/elasticdump
   ```

2. **Make scripts executable**:
   ```bash
   chmod +x examples/*.sh
   ```

3. **Update configuration**: Edit the scripts to match your Elasticsearch setup (hosts, indices, credentials, etc.)

## Available Examples

### 1. Basic Transfer (`basic-transfer.sh`)
Demonstrates simple data transfer between two Elasticsearch clusters.

**Usage:**
```bash
./basic-transfer.sh
```

### 2. Backup Data (`backup-data.sh`)
Shows how to backup Elasticsearch data to a file with timestamp.

**Usage:**
```bash
./backup-data.sh
```

### 3. Restore Data (`restore-data.sh`)
Demonstrates how to restore data from a backup file.

**Usage:**
```bash
./restore-data.sh
```

### 4. Transfer Mappings (`transfer-mappings.sh`)
Shows how to transfer only index mappings between clusters.

**Usage:**
```bash
./transfer-mappings.sh
```

### 5. Transfer Settings (`transfer-settings.sh`)
Demonstrates transferring only index settings between clusters.

**Usage:**
```bash
./transfer-settings.sh
```

### 6. Complete Migration (`complete-migration.sh`)
Shows a complete migration workflow (mappings + settings + data) with error handling.

**Usage:**
```bash
./complete-migration.sh
```

### 7. Authenticated Transfer (`authenticated-transfer.sh`)
Demonstrates data transfer with authentication credentials.

**Usage:**
```bash
./authenticated-transfer.sh
```

**Security Note:** For production use, consider using environment variables for credentials.

### 8. Limited Transfer (`limited-transfer.sh`)
Shows how to transfer only a specific number of documents.

**Usage:**
```bash
./limited-transfer.sh
```

### 9. Batch Processing (`batch-processing.sh`)
Demonstrates processing multiple indices in a batch operation.

**Usage:**
```bash
./batch-processing.sh
```

## Configuration

Before running any script, update the configuration variables at the top of each script:

- `SOURCE_HOST`: Source Elasticsearch cluster URL
- `DEST_HOST`: Destination Elasticsearch cluster URL  
- `INDEX_NAME`: Name of the index to process
- Authentication credentials (if needed)

## Common Parameters

Most scripts use these common elasticdump parameters:

- `--input`: Source (Elasticsearch URL or file path)
- `--output`: Destination (Elasticsearch URL or file path)
- `--type`: Type of data to transfer (`data`, `mapping`, `settings`)
- `--concurrency`: Number of concurrent operations (default: 4)
- `--scrollSize`: Scroll size for large datasets (default: 1000)
- `--limit`: Maximum number of documents to transfer
- `--format`: Output format (`json`, `ndjson`)
- `--verbose`: Enable verbose output

## Tips

1. **Test with small datasets first** before running on production data
2. **Monitor resource usage** during large transfers
3. **Use appropriate concurrency settings** based on your cluster capacity
4. **Always backup critical data** before migration
5. **Verify transfers** by comparing document counts and sample data

## Troubleshooting

If you encounter issues:

1. Check Elasticsearch connectivity: `curl -X GET 'http://localhost:9200/_cluster/health?pretty'`
2. Verify index exists: `curl -X GET 'http://localhost:9200/_cat/indices'`
3. Check elasticdump logs with `--verbose` flag
4. Ensure sufficient disk space for backup operations
5. Verify authentication credentials if using secured clusters
