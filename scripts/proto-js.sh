#!/usr/bin/env bash

cd ../front6

cp ../schema/api.proto .
protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto --ts_proto_opt=outputEncodeMethods=false --ts_proto_opt=stringEnums=true --ts_proto_out=src/api ./api.proto
sed -i "" "s|encode|toJSON|g" src/api/api.ts
sed -i "" "s|decode|fromJSON|g" src/api/api.ts

sed -i "" "s|.finish();|;\n    //@ts-ignore-line|g" src/api/api.ts

sed -i "" "s|_m0.Reader.create(data)|data|g" src/api/api.ts

rm api.proto
