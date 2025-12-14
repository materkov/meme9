#!/usr/bin/env bash

# Generate TypeScript types from proto files
cd "$(dirname "$0")/.."

mkdir -p src/schema

# Generate types for each service proto file
# Use jsonName=true to preserve snake_case field names in JSON (matching Twirp output)
protoc \
  --plugin=./node_modules/.bin/protoc-gen-ts_proto \
  --proto_path=../schema \
  --ts_proto_out=src/schema \
  --ts_proto_opt=esModuleInterop=true,outputServices=generic-definitions,outputClientImpl=false,useJsonName=true \
  ../schema/feed.proto \
  ../schema/posts.proto \
  ../schema/auth.proto \
  ../schema/users.proto \
  ../schema/subscriptions.proto

echo "TypeScript types generated in src/schema/"

