#!/bin/bash

# Release script for elasticdump
# Usage: ./scripts/release.sh [version]

set -e

VERSION=${1:-}

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.3"
    exit 1
fi

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo "Error: Version must be in format v1.2.3 or v1.2.3-beta"
    exit 1
fi

echo "🚀 Preparing release $VERSION"

# Check if we're on main/master branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
    echo "Error: Must be on main or master branch to create a release"
    exit 1
fi

# Check if working directory is clean
if [[ -n $(git status --porcelain) ]]; then
    echo "Error: Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
make test

# Run linting
echo "🔍 Running linter..."
make lint || echo "Warning: Linting failed, but continuing..."

# Build the project
echo "🔨 Building project..."
make build

# Test the binary
echo "🧪 Testing binary..."
./bin/elasticdump --version

# Create git tag
echo "📝 Creating git tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

# Push tag
echo "⬆️  Pushing tag to remote..."
git push origin "$VERSION"

echo "✅ Release $VERSION created successfully!"
echo ""
echo "🎉 GitHub Actions will now:"
echo "   - Run tests"
echo "   - Build binaries for multiple platforms"
echo "   - Create GitHub release"
echo "   - Build and push Docker images"
echo ""
echo "Check the Actions tab on GitHub for progress."
