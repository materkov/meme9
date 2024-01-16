#!/usr/bin/env bash

cd ~/mypage

export DOCKER_HOST=ssh://mypage

docker compose up -d --build meme9-web
