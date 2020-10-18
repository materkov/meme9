run-proto:
	mkdir -p pb
	protoc --gogofaster_out=pb \
		--proto_path ../schema \
		--proto_path ~/go/pkg/mod/github.com/gogo/protobuf@v1.3.1/ \
		../schema/login.proto
