#!/usr/bin/env bash

# Run grpcui from Docker
# The web UI will be available at http://localhost:8082
# It will connect to the gRPC server at host.docker.internal:8081
docker run --rm -it \
    -p 8082:8080 \
    fullstorydev/grpcui:latest \
    -plaintext \
    -vvv \
    host.docker.internal:8081
