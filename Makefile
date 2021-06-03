run-proto:
	mkdir -p api/pb
	protoc --gogofaster_out=api/pb \
		--proto_path schema \
		--proto_path ~/go/pkg/mod/github.com/gogo/protobuf@v1.3.1/ \
		schema/api.proto

	protoc --plugin=./front/node_modules/.bin/protoc-gen-ts_proto \
		--proto_path=schema \
		--ts_proto_out=front/src/schema \
		schema/api.proto

run-proto2:
	mkdir -p web/pb
	protoc --go_out=web \
		--proto_path schema \
		schema/*.proto

	#find web/pb -name '*.twirp.go' -type f | xargs sed -i '' -e 's/google.golang.org\/protobuf/github.com\/gogo\/protobuf/g'

build-api:
	cd api && go build cmd/main.go

build-front:
	cd front && yarn && yarn build

test-all:
	cd api && go test -v ./...

lint-all:
	cd web && test -z $$(gofmt -l . | tee /dev/stderr)

	cd web && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.32.1
	cd web && bin/golangci-lint --version
	cd web && bin/golangci-lint run

deploy:
	cd ~/mypage && docker-compose --context=mypage up -d --build meme9-front meme9-web
