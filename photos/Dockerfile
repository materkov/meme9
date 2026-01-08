FROM golang:1.24.9-alpine3.22
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app /build/cmd

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=0 /app /app
CMD ["/app"]
