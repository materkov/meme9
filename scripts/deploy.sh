#!/usr/bin/env bash

cd ~/mypage

export DOCKER_HOST=ssh://mypage
export DOCKER_API_VERSION=1.43

docker compose up -d --build meme9-web
docker compose up -d --build meme9-front
docker compose up -d --build meme9-rss
