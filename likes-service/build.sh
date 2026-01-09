#!/usr/bin/env bash

# Build
cd ..
docker build --platform=linux/amd64 -f likes-service/Dockerfile -t 7385cbca-brainy-vulpecula.registry.twcstorage.ru/likes-service/likes-service:latest .
docker push --platform=linux/amd64 7385cbca-brainy-vulpecula.registry.twcstorage.ru/likes-service/likes-service:latest

# Start image
cd ~/mypage
export DOCKER_HOST=ssh://mypage-ru
docker compose up -d --build likes-service
