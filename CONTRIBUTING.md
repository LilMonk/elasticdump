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
   make install dev-deps
   ```
3. Build the project:
   ```bash
   make build
   ```
4. Run tests:
   ```bash
   make test
   ```

## Code Style

- Follow standard Go conventions
- Use `make format` to format your code
- Run `make lint` to check for common mistakes
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

## CI/CD and Releases

### Automated Testing
All pull requests automatically trigger:
- **CI Workflow**: Runs tests, linting, and builds on multiple Go versions
- **Examples Testing**: Validates example scripts against a real Elasticsearch instance
- **Code Quality**: Runs golangci-lint for code quality checks

### Release Process
Releases are automated through GitHub Actions:

1. **Create a release**: Use the release script:
   ```bash
   ./scripts/release.sh v1.2.3
   ```

2. **Automated release workflow** will:
   - Run comprehensive tests
   - Build binaries for multiple platforms (Linux, macOS, Windows)
   - Create GitHub release with changelog
   - Build and push Docker images
   - Generate release notes

### Manual Testing
Before releasing:
```bash
# Run all tests
make test

# Run linting
make lint

# Build and test binary
make build
./bin/elasticdump --version
```

### Version Management
- Version is managed through Git tags (e.g., `v1.2.3`)
- Semantic versioning is used: `MAJOR.MINOR.PATCH`
- Development builds use version `dev`
