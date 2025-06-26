# Elasticdump

A powerful Go CLI tool to migrate data between Elasticsearch clusters. Elasticdump can also be used to backup and restore data with support for various Elasticsearch versions and efficient handling of large datasets.

## Features

- 🚀 **Fast Data Migration**: Migrate data between Elasticsearch clusters efficiently
- 💾 **Backup & Restore**: Comprehensive backup and restore functionality
- 🔄 **Multi-Version Support**: Compatible with various Elasticsearch versions
- ⚡ **High Performance**: Multi-threaded operations for faster processing
- 📊 **Progress Tracking**: Real-time progress bars for long-running operations
- 📝 **Multiple Formats**: Support for JSON and NDJSON output formats
- 🎯 **Flexible Operations**: Transfer data, mappings, or settings independently

## Installation

### Using Go Install

```bash
go install github.com/lilmonk/elasticdump@latest
```

### Building from Source

```bash
git clone https://github.com/lilmonk/elasticdump.git
cd elasticdump
go build -o elasticdump
```

### Using Releases

Download the latest binary from the [releases page](https://github.com/lilmonk/elasticdump/releases).

## Usage

### Basic Transfer

Transfer data between two Elasticsearch clusters:

```bash
elasticdump transfer --input=http://localhost:9200/source_index --output=http://localhost:9200/dest_index
```

### Backup Data

Backup data to a file:

```bash
elasticdump backup --input=http://localhost:9200/myindex --output=backup.ndjson
```

### Restore Data

Restore data from a backup file:

```bash
elasticdump restore --input=backup.ndjson --output=http://localhost:9200/myindex
```

### Transfer Mappings

Transfer only index mappings:

```bash
elasticdump transfer --input=http://localhost:9200/myindex --output=http://localhost:9200/myindex --type=mapping
```

### Transfer Settings

Transfer only index settings:

```bash
elasticdump transfer --input=http://localhost:9200/myindex --output=http://localhost:9200/myindex --type=settings
```

## Commands

### `transfer`

Transfer data, mappings, or settings between Elasticsearch clusters or to/from files.

```bash
elasticdump transfer [flags]
```

**Flags:**
- `--input, -i`: Source Elasticsearch cluster or index (required)
- `--output, -o`: Destination Elasticsearch cluster or index (required)
- `--type, -t`: Type of data to transfer (`data`, `mapping`, `settings`) (default: "data")
- `--limit, -l`: Limit the number of records to transfer (0 = no limit) (default: 0)
- `--concurrency, -c`: Number of concurrent operations (default: 4)
- `--format, -f`: Output format (`json`, `ndjson`) (default: "json")
- `--scrollSize, -s`: Size of the scroll for large datasets (default: 1000)
- `--username, -u`: Username for Elasticsearch authentication
- `--password, -p`: Password for Elasticsearch authentication

### `backup`

Backup Elasticsearch data to a file. This is a convenience wrapper around the transfer command.

```bash
elasticdump backup [flags]
```

**Flags:**
- `--input, -i`: Source Elasticsearch cluster or index (required)
- `--output, -o`: Output file path (required)
- `--type, -t`: Type of data to backup (`data`, `mapping`, `settings`) (default: "data")
- `--limit, -l`: Limit the number of records to backup (0 = no limit) (default: 0)
- `--concurrency, -c`: Number of concurrent operations (default: 4)
- `--format, -f`: Output format (`json`, `ndjson`) (default: "ndjson")
- `--scrollSize, -s`: Size of the scroll for large datasets (default: 1000)
- `--username, -u`: Username for Elasticsearch authentication
- `--password, -p`: Password for Elasticsearch authentication

### `restore`

Restore Elasticsearch data from a backup file.

```bash
elasticdump restore [flags]
```

**Flags:**
- `--input, -i`: Input file path (required)
- `--output, -o`: Destination Elasticsearch cluster or index (required)
- `--type, -t`: Type of data to restore (`data`, `mapping`, `settings`) (default: "data")
- `--concurrency, -c`: Number of concurrent operations (default: 4)
- `--username, -u`: Username for Elasticsearch authentication
- `--password, -p`: Password for Elasticsearch authentication

## Global Flags

- `--verbose, -v`: Verbose output
- `--help, -h`: Help for any command
- `--version`: Show version information

## Examples

### Complete Index Migration

Migrate an entire index including data, mappings, and settings:

```bash
# First, transfer settings and mappings
elasticdump transfer --input=http://source:9200/myindex --output=http://dest:9200/myindex --type=settings
elasticdump transfer --input=http://source:9200/myindex --output=http://dest:9200/myindex --type=mapping

# Then transfer the data
elasticdump transfer --input=http://source:9200/myindex --output=http://dest:9200/myindex --type=data --concurrency=8
```

### Large Dataset with Progress

For large datasets, increase concurrency and scroll size:

```bash
elasticdump transfer \
  --input=http://localhost:9200/large_index \
  --output=http://newcluster:9200/large_index \
  --concurrency=10 \
  --scrollSize=5000 \
  --verbose
```

### Partial Backup

Backup only a subset of documents:

```bash
elasticdump backup \
  --input=http://localhost:9200/myindex \
  --output=partial_backup.ndjson \
  --limit=10000 \
  --format=ndjson
```

### Authentication

Elasticdump supports basic authentication methods for clusters requiring authentication:

#### Using Username and Password Flags

```bash
elasticdump transfer \
  --input=http://source.elasticsearch.com:9200/index \
  --output=http://dest.elasticsearch.com:9200/index \
  --username=elastic \
  --password=your_password
```

## Performance Tips

1. **Increase Concurrency**: Use `--concurrency` flag to increase parallel operations
2. **Optimize Scroll Size**: Adjust `--scrollSize` based on document size and available memory
3. **Use NDJSON Format**: For large datasets, NDJSON format is more memory efficient
4. **Network Proximity**: Run elasticdump close to your Elasticsearch clusters to reduce network latency

## Error Handling

Elasticdump includes robust error handling:

- Automatic retries for transient network errors
- Detailed error messages for debugging
- Graceful handling of malformed documents

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/lilmonk/elasticdump/wiki)
- 🐛 [Issues](https://github.com/lilmonk/elasticdump/issues)
- 💬 [Discussions](https://github.com/lilmonk/elasticdump/discussions)

## Acknowledgments

- [Elasticsearch Go Client](https://github.com/elastic/go-elasticsearch)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
