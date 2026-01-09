#!/usr/bin/env bash

# Unified script to generate protobuf code for both Go (api) and TypeScript (front8)
set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

SCHEMA_DIR="$SCRIPT_DIR"
API_DIR="$PROJECT_ROOT/api"
FRONT8_DIR="$PROJECT_ROOT/front8"

echo "Generating protobuf code..."
echo "Schema directory: $SCHEMA_DIR"
echo "Go output: $API_DIR/pb"
echo "TypeScript output: $FRONT8_DIR/src/schema"

# Generate Go protobuf code for api module
echo ""
echo "=== Generating Go code for api module ==="
cd "$API_DIR"

# Generate all service proto files (generate separately to avoid go_package conflicts)
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/posts.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/auth.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/users.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/subscriptions.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/likes.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/photos.proto"

echo "Go code generated successfully"

# Generate TypeScript protobuf code for front8
echo ""
echo "=== Generating TypeScript code for front8 ==="
cd "$FRONT8_DIR"

# Ensure schema directory exists
mkdir -p src/schema

# Generate types using @protobuf-ts/plugin and Twirp clients using twirp-ts
protoc \
  --plugin=./node_modules/.bin/protoc-gen-ts \
  --plugin=./node_modules/.bin/protoc-gen-twirp_ts \
  --proto_path="$SCHEMA_DIR" \
  --ts_out=src/schema \
  --twirp_ts_out=src/schema \
  "$SCHEMA_DIR/posts.proto" \
  "$SCHEMA_DIR/auth.proto" \
  "$SCHEMA_DIR/users.proto" \
  "$SCHEMA_DIR/subscriptions.proto" \
  "$SCHEMA_DIR/photos.proto" \
  "$SCHEMA_DIR/likes.proto"

echo "TypeScript code generated successfully"
echo ""
echo "✅ All protobuf code generation complete!"
