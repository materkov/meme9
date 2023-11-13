#!/usr/bin/env bash

cd ../web6/
go version
golangci-lint version
golangci-lint run
