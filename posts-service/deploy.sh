#!/usr/bin/env bash

cd ~/mypage

export DOCKER_HOST=ssh://mypage-ru

docker compose up -d --build posts-service
