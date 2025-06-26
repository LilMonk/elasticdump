# Contributing to Elasticdump

Thank you for your interest in contributing to Elasticdump! We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/elasticdump.git
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

1. Ensure you have Go 1.21+ installed
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run tests:
   ```bash
   go test ./...
   ```
4. Build the project:
   ```bash
   go build -o elasticdump
   ```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Add tests for new functionality

## Submitting Changes

1. Commit your changes with a descriptive message
2. Push to your fork
3. Create a pull request with:
   - Clear description of changes
   - Reference to any related issues
   - Test results

## Reporting Issues

When reporting issues, please include:
- Go version
- Operating system
- Elasticsearch version
- Full command and error output
- Steps to reproduce

## Feature Requests

Feature requests are welcome! Please:
- Check existing issues first
- Provide clear use case description
- Include example usage if possible
