FROM golang:1.22.1-alpine3.19
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/materkov/meme9/web6/src/pkg.BuildTime=$(date)'" -o /meme9 /build/cmd

FROM alpine:3.19.1
RUN apk add --no-cache ca-certificates
COPY --from=0 /meme9 /meme9
CMD ["/meme9"]
