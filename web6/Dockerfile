FROM golang:1.19.1-alpine3.16
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/materkov/meme9/web6/pkg.BuildTime=$(date)'" -o /meme9 /build/cmd

FROM alpine:3.16
RUN apk add --no-cache ca-certificates
COPY --from=0 /meme9 /meme9
CMD ["/meme9"]
