FROM golang:1.21.0-alpine3.17
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/materkov/meme9/web6/pkg.BuildTime=$(date)'" -o /meme9 /build/cmd

FROM alpine:3.17
RUN apk add --no-cache ca-certificates
COPY --from=0 /meme9 /meme9
CMD ["/meme9"]
