#!/usr/bin/env bash

# Build
#cd ..
docker build --platform=linux/amd64 -t 7385cbca-brainy-vulpecula.registry.twcstorage.ru/front8/front8:latest .
docker push --platform=linux/amd64 7385cbca-brainy-vulpecula.registry.twcstorage.ru/front8/front8:latest

# Start image
cd ~/mypage
export DOCKER_HOST=ssh://mypage-ru
docker compose up -d --pull always static8
