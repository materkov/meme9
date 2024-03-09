FROM golang:1.22.1-alpine3.19
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app /build/cmd

FROM alpine:3.19.1
RUN apk add --no-cache ca-certificates
COPY --from=0 /app /app
CMD ["/app"]
