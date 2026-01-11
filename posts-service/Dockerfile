FROM golang:1.24.9-alpine3.22
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build

# Copy api module first (needed for replace directive)
COPY api /api
# Copy service go.mod and go.sum
COPY posts-service/go.mod posts-service/go.sum ./
# Download dependencies
RUN go mod download

# Copy service code
COPY posts-service .
RUN go build -o /posts-service ./cmd

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=0 /posts-service /posts-service
CMD ["/posts-service"]

