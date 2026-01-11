FROM golang:1.24.9-alpine3.22
RUN apk add --no-cache git gcc musl-dev
WORKDIR /build

# Copy api module first (needed for replace directive)
COPY api /api
# Copy service go.mod and go.sum
COPY subscriptions-service/go.mod subscriptions-service/go.sum ./
# Download dependencies
RUN go mod download

# Copy service code
COPY subscriptions-service .
RUN go build -o /subscriptions-service ./cmd

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=0 /subscriptions-service /subscriptions-service
CMD ["/subscriptions-service"]

