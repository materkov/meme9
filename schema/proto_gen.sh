#!/usr/bin/env bash

# Unified script to generate protobuf code for both Go (web7) and TypeScript (front8)
set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

SCHEMA_DIR="$SCRIPT_DIR"
WEB7_DIR="$PROJECT_ROOT/web7"
FRONT8_DIR="$PROJECT_ROOT/front8"

echo "Generating protobuf code..."
echo "Schema directory: $SCHEMA_DIR"
echo "Go output: $WEB7_DIR/pb"
echo "TypeScript output: $FRONT8_DIR/src/schema"

# Generate Go protobuf code for web7
echo ""
echo "=== Generating Go code for web7 ==="
cd "$WEB7_DIR"

# Generate all service proto files (generate separately to avoid go_package conflicts)
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/posts.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/auth.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/users.proto"
protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/subscriptions.proto"

# Generate json_api.proto if it exists
if [ -f "$SCHEMA_DIR/json_api.proto" ]; then
    protoc -I "$SCHEMA_DIR" --go_out=pb --twirp_out=pb "$SCHEMA_DIR/json_api.proto"
fi

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
  "$SCHEMA_DIR/subscriptions.proto"

echo "TypeScript code generated successfully"
echo ""
echo "✅ All protobuf code generation complete!"
