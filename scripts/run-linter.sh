#!/usr/bin/env bash

go version
golangci-lint version

cd ../web6/
golangci-lint run

cd ../api/
golangci-lint run

cd ../realtime/
golangci-lint run
