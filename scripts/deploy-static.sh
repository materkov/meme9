#!/usr/bin/env bash

cd ~/mypage

export DOCKER_HOST=ssh://mypage
export DOCKER_API_VERSION=1.43

docker compose up -d --build static7
