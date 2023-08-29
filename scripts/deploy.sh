#!/usr/bin/env bash

cd ~/mypage

export DOCKER_HOST=ssh://mypage

docker compose up -d --build meme9-web
#docker-compose --context mypage up -d --build meme9-worker
docker compose  up -d --build meme9-front
