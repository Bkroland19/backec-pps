#!/bin/bash

echo "Fixing Go dependencies..."

# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Verify dependencies
go mod verify

echo "Dependencies fixed successfully!"
echo "You can now run: go run main.go"
