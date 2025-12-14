#!/usr/bin/env bash

# Generate TypeScript types and Twirp clients from proto files
cd "$(dirname "$0")/.."

mkdir -p src/schema

# Generate types using @protobuf-ts/plugin and Twirp clients using twirp-ts
protoc \
  --plugin=./node_modules/.bin/protoc-gen-ts \
  --plugin=./node_modules/.bin/protoc-gen-twirp_ts \
  --proto_path=../schema \
  --ts_out=src/schema \
  --twirp_ts_out=src/schema \
  ../schema/posts.proto \
  ../schema/auth.proto \
  ../schema/users.proto \
  ../schema/subscriptions.proto

echo "TypeScript types and Twirp clients generated in src/schema/"
