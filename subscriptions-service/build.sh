#!/usr/bin/env bash

# Build
cd ..
docker build --platform=linux/amd64 -f subscriptions-service/Dockerfile -t 7385cbca-brainy-vulpecula.registry.twcstorage.ru/subscriptions-service/subscriptions-service:latest .
docker push --platform=linux/amd64 7385cbca-brainy-vulpecula.registry.twcstorage.ru/subscriptions-service/subscriptions-service:latest

# Start image
cd ~/mypage
export DOCKER_HOST=ssh://mypage-ru
docker compose up -d --pull always subscriptions-service
