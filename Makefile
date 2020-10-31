run-proto:
	mkdir -p api/pb
	protoc --gogofaster_out=api/pb \
		--proto_path schema \
		--proto_path ~/go/pkg/mod/github.com/gogo/protobuf@v1.3.1/ \
		schema/login.proto

	protoc --plugin=./front/node_modules/.bin/protoc-gen-ts_proto \
		--proto_path=schema \
		--ts_proto_out=front/src/schema \
		schema/login.proto
