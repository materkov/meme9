#!/usr/bin/env bash

cd ../api
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../web6
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../realtime
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../web7
# Generate all new service proto files (generate separately to avoid go_package conflicts)
protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/feed.proto
protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/posts.proto
protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/auth.proto
protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/users.proto
protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/subscriptions.proto
