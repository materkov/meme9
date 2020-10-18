FROM golang:1.15.2-alpine3.12
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /meme9 /build/cmd/main.go

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
EXPOSE 8001
COPY --from=0 /meme9 /meme9
CMD ["/meme9"]
