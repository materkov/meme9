#!/usr/bin/env bash

/Users/m.materkov/Downloads/protoc-25.2-osx-aarch_64/bin/protoc -I ../schema --go_out=pb --twirp_out=pb ../schema/api.proto
