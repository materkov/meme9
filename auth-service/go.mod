module github.com/materkov/meme9/auth-service

go 1.24.2

require (
	github.com/materkov/meme9/api v0.0.0-20260108223855-e3a4be80c9f6
	github.com/twitchtv/twirp v8.1.3+incompatible
	go.mongodb.org/mongo-driver v1.17.6
	golang.org/x/crypto v0.43.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/materkov/meme9/api => ../api
