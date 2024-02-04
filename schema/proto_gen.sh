#!/usr/bin/env bash

cd ../api
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../web6
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../rss2
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../imgproxy
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto

cd ../realtime
/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto
