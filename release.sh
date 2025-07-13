#!/bin/bash

# Automated release script for querybuilder
# Usage: ./release.sh [patch|minor|major] [message]

set -e

# Default to patch if no argument provided
BUMP_TYPE=${1:-patch}
MESSAGE=${2:-"Release"}

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Current version: $CURRENT_VERSION"

# Remove 'v' prefix for version manipulation
CURRENT_VERSION=${CURRENT_VERSION#v}

# Split version into parts
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Bump version based on type
case $BUMP_TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Invalid bump type. Use: patch, minor, or major"
        exit 1
        ;;
esac

NEW_VERSION="v$MAJOR.$MINOR.$PATCH"
echo "New version: $NEW_VERSION"

# Confirm release
read -p "Create release $NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Release cancelled"
    exit 1
fi

# Run tests first
echo "Running tests..."
go test -v ./...

# Create and push tag
echo "Creating tag..."
git tag -a "$NEW_VERSION" -m "$MESSAGE $NEW_VERSION"

echo "Pushing tag..."
git push origin "$NEW_VERSION"

echo "✅ Release $NEW_VERSION created successfully!"
echo "Check progress at: https://github.com/hojabri/querybuilder/actions"