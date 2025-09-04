#!/bin/bash

# Build script for generating Go protobuf files
# This script generates the protobuf files needed for the Go API service

set -e

echo "🔧 Building protobuf files for Go API service..."

# Create proto directory if it doesn't exist
mkdir -p proto

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install it first:"
    echo "   macOS: brew install protobuf"
    echo "   Ubuntu: sudo apt-get install protobuf-compiler"
    echo "   Windows: Download from https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate Go files from protobuf
echo "🚀 Generating Go protobuf files..."
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    proto_src/steelbeam.proto

# Move generated files to proto directory
mv proto_src/*.pb.go proto/ 2>/dev/null || true

echo "✅ Protobuf files generated successfully!"
echo "📁 Files created:"
ls -la proto/*.pb.go

echo ""
echo "🎯 Ready for deployment! The Go API service now has all required protobuf files."
