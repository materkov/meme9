FROM golang:1.21.0-alpine3.17
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/materkov/meme9/api/src/pkg.BuildTime=$(date)'" -o /app /build/cmd

FROM alpine:3.17
RUN apk add --no-cache ca-certificates
COPY --from=0 /app /app
CMD ["/app"]
